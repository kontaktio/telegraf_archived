from kapacitor.udf.agent import Agent, Handler
import kapacitor.udf.udf_pb2 as udf_pb2
import sys
import requests
import hashlib
import struct
from datetime import datetime, timedelta
from uuid import UUID

API_KEY_PARAM = 'api_key'
LOCATION_ID_PARAM = 'location_id'
INFSOFT_URL = 'https://api.infsoft.com'
INFSOFT_BLE_SCAN_PATH = '/v1/locatornode-scan/ble'

class InfsoftWriter(Handler):
    def __init__(self):
        self._api_key = ''
        self._location_id = 0
        self._interval = timedelta(seconds=0)
        self._last_flush = datetime.now()
        self._points = {}

    def info(self):
        response = udf_pb2.Response()
        response.info.wants = udf_pb2.STREAM
        response.info.provides = udf_pb2.STREAM

        response.info.options['apikey'].valueTypes.append(udf_pb2.STRING)
        response.info.options['locationid'].valueTypes.append(udf_pb2.INT)
        response.info.options['interval'].valueTypes.append(udf_pb2.INT)

        return response

    def init(self, init_req):
        for opt in init_req.options:
            if opt.name == 'apikey':
                self._api_key = opt.values[0].stringValue
            elif opt.name == 'locationid':
                self._location_id = opt.values[0].intValue
            elif opt.name == 'interval':
                self._interval = opt.values[0].intValue

        print >> sys.stderr, 'ApiKey: %s, LocationId: %d, Interval: %s' % (self._api_key, self._location_id, str(self._interval))

        success = True
        msg = ''

        if self._api_key == '':
            success = False
            msg = 'apikey is required '
        if self._location_id == 0:
            success = False
            msg = msg + 'locationid is required '
        if self._interval == 0:
            success = False
            msg = msg + 'interval is required '

        response = udf_pb2.Response()
        response.init.success = success
        response.init.error = msg

        return response


    def point(self, point):
        print >> sys.stderr, str(point.fieldsString)
        print >> sys.stderr, str(point.fieldsInt)
        trackingId = point.fieldsString['trackingId']
        sourceId = point.fieldsString['sourceId']
        rssi = point.fieldsInt['rssi']
        tx_power = point.fieldsInt['txPower']

        self._points[(trackingId, sourceId)] = (rssi, tx_power)
        if datetime.now() - self._last_flush > timedelta(seconds=self._interval):
            self._emit_points()
            self._last_flush = datetime.now()
            self._points.clear()

    def _emit_points(self):
        pass

    def _push_to_infsoft(self, scans):
        """
        scans = {
            "source_id": sourceid,
            "items": [
                (tracking_id, rssi, tx_power),
                (tracking_id, rssi, tx_power)
            ]
        }
        """
        items = []
        for tracking_id, rssi, tx_power in scans['items']:
            proximity, major, minor, mac = self._extract_proximity_major_minor_mac(tracking_id)
            items.append({
                "ble_mac": mac,
                "ble_proximityuuid": proximity,
                "ble_major": major,
                "ble_minor": minor,
                "ble_txpower": tx_power,
                "ble_databyte": 0,
                "rssi": rssi
            })
        packet = {
            "ln_eth_mac": self._extract_mac(scans['source_id']),
            "items": items
        }

        requests.post(INFSOFT_URL + INFSOFT_BLE_SCAN_PATH, json=packet, params={
            LOCATION_ID_PARAM: self._location_id,
            API_KEY_PARAM: self._api_key
        })


    def _extract_proximity_major_minor_mac(self, tracking_id):
        hashed = hashlib.md5(tracking_id).digest()
        proximity = UUID(bytes=tracking_id.ljust(16, '\x00'))
        major, minor = struct.unpack('HH', hashed[:4])
        mac = self._extract_mac(tracking_id)
        return (proximity, major, minor, mac)

    def _extract_mac(self, value):
        hashed = hashlib.md5(source_id)
        mac_without_separators = hashed.hexdigest()[:12]
        return ':'.join([mac_without_separators[i:i+2] for i in range(0, 12, 2)])

if __name__ == '__main__':
    agent = Agent()
    handler = InfsoftWriter()
    agent.handler = handler

    print >> sys.stderr, "Starting Infsoft Writer"

    agent.start()
    agent.wait()

    print >> sys.stderr, "Infsoft Writer finished"
