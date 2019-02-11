import argparse
import ConfigParser
import sys
import itertools

from utils import chunks
from api_client import ApiClient
from influx_client import InfluxClient
from telegraf_config import TelegrafConfigFormatter


class Options(object): 
    CONFIG_FILE_SECTION = 'default'

    def __init__(self, args):
        parser = argparse.ArgumentParser(
            description="Generate Streams Telegraf config for K.io API Account",
            formatter_class=argparse.ArgumentDefaultsHelpFormatter)
        parser.add_argument("--api-key", dest="api_key", required=True)
        parser.add_argument("--influxdb-url", dest="influxdb_url", required=True)
        parser.add_argument("--influxdb-port", dest="influxdb_port", default=8086, type=int)
        parser.add_argument("--influxdb-username", dest="influxdb_username", required=True)
        parser.add_argument("--influxdb-password", dest="influxdb_password", required=True)
        parser.add_argument("--influxdb-new-user-password", dest="influxdb_new_user_password", required=True)
        parser.add_argument("--config-file", dest="config_file")
        parser.add_argument("--api-url", dest="api_url", default="https://testapi.kontakt.io/")
        parser.add_argument("--api-venue-id", dest="api_venue_id", default=None)
        parser.add_argument("--mqtt-url", dest="mqtt_url", default="ssl://testrtls.kontakt.io:8083")
        parser.add_argument("--telemetry-types", dest="telemetry_types", nargs='+', default = ['all'])
        parser.add_argument("--streams-per-telegraf", dest="streams_per_telegraf", type=int, default=250)

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

    def get_influx_url(self):
        return self.args['influxdb_url']

    def get_influx_port(self):
        return self.args['influxdb_port']

    def get_influx_username(self):
        return self.args['influxdb_username']

    def get_influx_password(self):
        return self.args['influxdb_password']

    def get_mqtt_url(self):
        return self.args['mqtt_url']

    def get_telemetry_types(self):
        return self.args['telemetry_types']

    def get_streams_per_telegraf(self):
        return self.args['streams_per_telegraf']

    def get_influxdb_new_user_password(self):
        return self.args['influxdb_new_user_password']


options = Options(sys.argv[1:])
api_client = ApiClient(options.get_api_url(), options.get_api_key())
influx_client = InfluxClient(options.get_influx_url(),
    options.get_influx_port(),
    options.get_influx_username(), 
    options.get_influx_password())

company_id = api_client.get_company_id()
influxdb_new_user_password = options.get_influxdb_new_user_password()
influx_client.create_database(company_id)
influx_client.create_user(company_id, influxdb_new_user_password, company_id)

influx_client.create_retention_policy(company_id, 'stream_rp', '3h')

influx_client.create_retention_policy(company_id, 'current_rp', '7d')
influx_client.recreate_continuous_query(company_id, '10s', 'current_rp', 'stream_rp', '10s', '40s')

influx_client.create_retention_policy(company_id, 'recent_rp', '365d')
influx_client.recreate_continuous_query(company_id, '5m', 'recent_rp', 'stream_rp', '5m', '5m')

influx_client.create_retention_policy(company_id, 'history_rp', '365d')
influx_client.recreate_continuous_query(company_id, '1h', 'history_rp', 'stream_rp', '1h', '1h')

influx_client.create_retention_policy(company_id, 'temporary_history_rp', '365d')

unique_ids = api_client.get_telemetry_unique_ids(api_venue_id=options.get_api_venue_id())

topics = ['/stream/%s/%s' % elem for elem in 
    itertools.product(unique_ids, options.get_telemetry_types())]

topic_chunks = chunks(topics, options.get_streams_per_telegraf())

cfg = TelegrafConfigFormatter()

cfg.append_section_name('global_tags')
cfg.append_section_name('agent')
cfg.append_key_value('interval', '5s')
cfg.append_key_value('round_interval', False)
cfg.append_key_value('metric_buffer_limit', 1000)
cfg.append_key_value('flush_buffer_when_full', True)
cfg.append_key_value('collection_jitter', '0s')
cfg.append_key_value('flush_interval', '5s')
cfg.append_key_value('flush_jitter', 0)
cfg.append_key_value('omit_hostname', True)
cfg.append_key_value('debug', True)
cfg.append_key_value('quiet', False)
cfg.append_key_value('logfile', '/var/log/telegraf-config-gen.log')

cfg.append_section_name('outputs.influxdb', True)
cfg.append_key_value('urls', ["%s:%d" % (options.get_influx_url(), options.get_influx_port())])
cfg.append_key_value('database', company_id)
cfg.append_key_value('username', options.get_influx_username())
cfg.append_key_value('password', options.get_influx_password())
cfg.append_key_value('precision', 's')
cfg.append_key_value('timeout', '5s')
cfg.append_key_value('retention_policy', 'stream_rp')

cfg.append_section_name('processors.clickdetect', True)
cfg.append_key_value('field_name', 'clickId')
cfg.append_key_value('out_field_name', 'singleClick')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('order', 0)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastSingleClick')
cfg.append_key_value('out_field_name', 'singleClick')
cfg.append_key_value('threshold', 60)
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('order', 1)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastDoubleTap')
cfg.append_key_value('out_field_name', 'doubleTap')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('threshold', 60)
cfg.append_key_value('order', 2)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastThreshold')
cfg.append_key_value('out_field_name', 'threshold')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('threshold', 60)
cfg.append_key_value('order', 3)

cfg.append_section_name('processors.override', True)
cfg.append_key_value('name_override', 'telemetry')

cfg.append_section_name('inputs.mqtt_consumer', True)
cfg.append_key_value('servers', [options.get_mqtt_url()])
cfg.append_key_value('qos', 0)
cfg.append_key_value('connection_timeout', '30s')
cfg.append_key_value('persistent_session', False)
cfg.append_key_value('username', 'telegraf')
cfg.append_key_value('password', options.get_api_key())
cfg.append_key_value('data_format', 'json')
cfg.append_key_value('tag_keys', ['trackingId'])

idx = 0

for topic_chunk in topic_chunks:
    current = TelegrafConfigFormatter(cfg)
    current.append_key_value('topics', topic_chunk)

    f = open("telegraf.stream.conf.%d" % idx, "w")
    f.write(current.to_string())
    f.close()
    idx = idx + 1
