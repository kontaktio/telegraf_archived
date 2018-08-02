from kapacitor.udf.agent import Agent, Handler
import kapacitor.udf.udf_pb2 as udf_pb2
import sys
import requests
import hashlib
import struct
from datetime import datetime, timedelta
from uuid import UUID
import json

API_KEY_PARAM = 'api_key'
LOCATION_ID_PARAM = 'location_id'
INFSOFT_URL = 'https://api.infsoft.com'
INFSOFT_BLE_SCAN_PATH = '/v1/locatornode-scan/ble'

class InfsoftWriter(Handler):
    def __init__(self, agent):
        self._api_key = ''
        self._location_id = 0
        self._interval = timedelta(seconds=0)
        self._last_flush = datetime.now()
        self._points = {}
        self._tx_power = 0
        self._agent = agent

    def info(self):
        response = udf_pb2.Response()
        response.info.wants = udf_pb2.STREAM
        response.info.provides = udf_pb2.STREAM

        response.info.options['apikey'].valueTypes.append(udf_pb2.STRING)
        response.info.options['locationid'].valueTypes.append(udf_pb2.INT)
        response.info.options['interval'].valueTypes.append(udf_pb2.INT)
        response.info.options['txpower'].valueTypes.append(udf_pb2.INT)

        return response

    def init(self, init_req):
        for opt in init_req.options:
            if opt.name == 'apikey':
                self._api_key = opt.values[0].stringValue
            elif opt.name == 'locationid':
                self._location_id = opt.values[0].intValue
            elif opt.name == 'interval':
                self._interval = opt.values[0].intValue
            elif opt.name == 'txpower':
                self._tx_power = opt.values[0].intValue

        print >> sys.stderr, 'ApiKey: %s, LocationId: %d, Interval: %s, TxPower: %d' % (
            self._api_key, 
            self._location_id, 
            str(self._interval),
            self._tx_power)

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
        if self._tx_power == 0:
            success = False
            msg = msg + 'txpower is required '

        response = udf_pb2.Response()
        response.init.success = success
        response.init.error = msg

        return response


    def point(self, point):        
        trackingId = point.tags['trackingId']
        sourceId = point.tags['sourceId']
        rssi = int(point.fieldsDouble['rssi'])

        self._points[(trackingId, sourceId)] = rssi
        if datetime.now() - self._last_flush > timedelta(seconds=self._interval):
            self._emit_points()
            self._last_flush = datetime.now()
            self._points.clear()

    def snapshot(self):
        response = udf_pb2.Response()
        response.snapshot.snapshot = ''
        return response

    def _emit_points(self):
        print >> sys.stderr, ",".join([
            "%s, %s -> %s" % (k, v, self._points[(k, v)]) for k, v in self._points
        ])
        scans = {}
        for tracking_id, source_id in self._points:
            if source_id not in scans:
                scans[source_id] = []

            rssi = self._points[(tracking_id, source_id)]
            scans[source_id].append((tracking_id, rssi))

        self._push_to_infsoft(scans)


    def _push_to_infsoft(self, scans):
        """
        scans = 
            {
                sourceid: [
                    (tracking_id, rssi),
                    (tracking_id, rssi)
                ],
                sourceid2: [...]
            }
        """
        packets = []
        for source_id in scans:
            items = []
            for tracking_id, rssi in scans[source_id]:
                proximity, major, minor, mac = self._extract_proximity_major_minor_mac(tracking_id)
                items.append({
                    "ble_mac": mac,
                    "ble_proximityuuid": str(proximity),
                    "ble_major": major,
                    "ble_minor": minor,
                    "ble_txpower": self._tx_power,
                    "ble_databyte": 0,
                    "rssi": rssi
                })
            packet = {
                "ln_eth_mac": self._extract_mac(source_id),
                "items": items
            }
            packets.append(packet)

        response = requests.post(INFSOFT_URL + INFSOFT_BLE_SCAN_PATH, data=json.dumps(packets), params={
            LOCATION_ID_PARAM: self._location_id,
            API_KEY_PARAM: self._api_key
        }, headers = {
            "Content-Type": "application/json"
        })
        print >> sys.stderr, "Infsoft responded with code %d" % response.status_code


    def _extract_proximity_major_minor_mac(self, tracking_id):
        hashed = hashlib.md5(tracking_id).digest()
        proximity = UUID(bytes=tracking_id.ljust(16, '\x00'))
        major, minor = struct.unpack('HH', hashed[:4])
        mac = self._extract_mac(tracking_id)
        return (proximity, major, minor, mac)

    def _extract_mac(self, value):
        hashed = hashlib.md5(value)
        mac_without_separators = hashed.hexdigest()[:12]
        return ':'.join([mac_without_separators[i:i+2] for i in range(0, 12, 2)])

if __name__ == '__main__':
    agent = Agent()
    handler = InfsoftWriter(agent)
    agent.handler = handler

    print >> sys.stderr, "Starting Infsoft Writer"

    agent.start()
    agent.wait()

    print >> sys.stderr, "Infsoft Writer finished"
