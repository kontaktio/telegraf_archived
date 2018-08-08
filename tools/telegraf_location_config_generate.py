import argparse
import ConfigParser
import sys
import itertools

from utils import chunks
from api_client import ApiClient
from influx_client import InfluxClient
from telegraf_config import TelegrafConfigFormatter
from kapacitor_client import KapacitorClient

INFSOFT_REALTIME_ENDPOINT = 'https://api.infsoft.com/v1/devices-realtime/ble'
KAPACITOR_POSITION_TASK_NAME = "position_%s"
KAPACITOR_LOCATION_TASK_NAME = "location_%s"

class Options(object): 
    CONFIG_FILE_SECTION = 'default'

    def __init__(self, args):
        parser = argparse.ArgumentParser(
            description='Generate Telegraf config for K.io API Account',
            formatter_class=argparse.ArgumentDefaultsHelpFormatter)
        parser.add_argument('--api-key', dest='api_key', required=True)
        parser.add_argument('--kapacitor-url', dest='kapacitor_url', required=True)
        parser.add_argument('--kapacitor-user', dest='kapacitor_user', required=True)
        parser.add_argument('--kapacitor-pass', dest='kapacitor_pass', required=True)
        parser.add_argument('--influxdb-url', dest='influxdb_url', required=True)
        parser.add_argument('--influxdb-port', dest='influxdb_port', default=8086, type=int)
        parser.add_argument('--influxdb-username', dest='influxdb_username', required=True)
        parser.add_argument('--influxdb-password', dest='influxdb_password', required=True)
        parser.add_argument('--config-file', dest='config_file')
        parser.add_argument('--api-url', dest='api_url', default='https://testapi.kontakt.io/')
        parser.add_argument('--api-venue-id', dest='api_venue_id', default=None)
        parser.add_argument('--data-collection-interval', dest='data_collection_interval', default='30s')
        parser.add_argument('--debug', dest='debug', default=False, type=bool)
        parser.add_argument('--log-file', dest='log_file', default='~/telegraf.log')
        parser.add_argument('--tx-power', dest='tx_power', type=int, default=-77)

        self.args = vars(parser.parse_args(args=args))
        if self.args['config_file'] is not None:
            self._load_from_file()

    def _load_from_file(self):
        parser = ConfigParser.ConfigParser()
        parser.read(self.args['config_file'])
        values = parser.options(self.CONFIG_FILE_SECTION)
        for k in values:
            self.args[k] = parser.get(self.CONFIG_FILE_SECTION, k)

    def get_api_key(self):
        return self.args['api_key']

    def get_api_url(self):
        return self.args['api_url']

    def get_api_venue_id(self):
        return self.args['api_venue_id']

    def get_data_collection_interval(self):
        return self.args['data_collection_interval']

    def get_debug(self):
        return self.args['debug']

    def get_log_file(self):
        return self.args['log_file']

    def get_influx_url(self):
        return self.args['influxdb_url']

    def get_influx_port(self):
        return self.args['influxdb_port']

    def get_influx_username(self):
        return self.args['influxdb_username']

    def get_influx_password(self):
        return self.args['influxdb_password']

    def get_kapacitor_url(self):
        return self.args['kapacitor_url']

    def get_kapacitor_user(self):
        return self.args['kapacitor_user']

    def get_kapacitor_pass(self):
        return self.args['kapacitor_pass']

    def get_tx_power(self):
        return self.args['tx_power']

DEFAULT_INFLUX_PASSWORD = 'aeg0UKOmUIzJap1QV6m8'

options = Options(sys.argv[1:])
api_client = ApiClient(options.get_api_url(), options.get_api_key())

company_id = api_client.get_company_id()
kapacitor_client = KapacitorClient(options, company_id, 'stream_rp')

location_task_name = KAPACITOR_LOCATION_TASK_NAME % company_id
kapacitor_location_task = kapacitor_client.get_task_info(location_task_name)
if 'error' in kapacitor_location_task:
    result = kapacitor_client.create_task(location_task_name, 'location-tpl', {
        'database': {
            'value': company_id,
            'type': 'string'
        }
    })
    if len(result['error']) > 0: 
        raise(Exception(result['error']))

location_engine_configs = api_client.get_location_engine_venues(options.get_api_venue_id())
if len(location_engine_configs) == 0:
    print 'None of selected account venues has Location Engine configured!'
    sys.exit(0)

cfg = TelegrafConfigFormatter()

cfg.append_section_name('global_tags')
cfg.append_section_name('agent')
cfg.append_key_value('interval', options.get_data_collection_interval())
cfg.append_key_value('round_interval', True) # default
cfg.append_key_value('metric_buffer_limit', 1000) # default
cfg.append_key_value('flush_buffer_when_full', True) # default
cfg.append_key_value('collection_jitter', '0s') # default
cfg.append_key_value('flush_interval', '10s') # default
cfg.append_key_value('flush_jitter', '0s') # default
cfg.append_key_value('debug', options.get_debug())
cfg.append_key_value('quiet', False)
cfg.append_key_value('logfile', options.get_log_file())

cfg.append_section_name('outputs.influxdb', True)
cfg.append_key_value('urls', ['%s:%d' % (options.get_influx_url(), options.get_influx_port())])
cfg.append_key_value('database', company_id)
cfg.append_key_value('username', options.get_influx_username())
cfg.append_key_value('password', options.get_influx_password())
cfg.append_key_value('precision', 's') # default
cfg.append_key_value('timeout', '5s') # default
cfg.append_key_value('retention_policy', 'current_rp')

cfg.append_section_name('processors.override', True)
cfg.append_key_value('name_override', 'position')

cfg.append_section_name('processors.trackingidresolver', True)
cfg.append_key_value('tag_name', 'ble_proximityuuid')
cfg.append_key_value('new_tag_name', 'trackingId')

cfg.append_section_name('inputs.http', True)
cfg.append_key_value('timeout', '5s')
cfg.append_key_value('data_format', 'json')
cfg.append_key_value('tag_keys', ['ble_proximityuuid'])

idx = 0
for venue_id, location_engine_config in location_engine_configs.iteritems():
    current = TelegrafConfigFormatter(cfg)
    infsoft_apikey = location_engine_config['infsoft']['apiKey']
    infsoft_locationid = location_engine_config['infsoft']['locationId']
    current.append_key_value('urls', [
        '%s?api_key=%s&location_id=%d' % (
            INFSOFT_REALTIME_ENDPOINT,
            infsoft_apikey,
            infsoft_locationid
        )
    ])

    with open('telegraf.location.conf.%d' % idx, 'w') as f:
        f.write(current.to_string())
        
    idx = idx + 1

    task_name = KAPACITOR_POSITION_TASK_NAME % venue_id
    kapacitor_task = kapacitor_client.get_task_info(task_name)
    if 'error' in kapacitor_task:
        result = kapacitor_client.create_task(task_name, 'infsoft-tpl',
        {
            'apikey': { 'value': infsoft_apikey, 'type': 'string' },
            'locationid': { 'value':  infsoft_locationid, 'type': 'int' },
            'txpower': { 'value': options.get_tx_power(), 'type': 'int' }
        })
        if len(result['error']) > 0: 
            raise(Exception(result['error']))