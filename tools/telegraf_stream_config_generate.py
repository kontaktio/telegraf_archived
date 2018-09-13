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
        parser.add_argument("--config-file", dest="config_file")
        parser.add_argument("--api-url", dest="api_url", default="https://testapi.kontakt.io/")
        parser.add_argument("--api-venue-id", dest="api_venue_id", default=None)
        parser.add_argument("--data-collection-interval", dest="data_collection_interval", default="30s")
        parser.add_argument("--debug", dest="debug", default=False, type=bool)
        parser.add_argument("--log-file", dest="log_file", default="~/telegraf.log")
        parser.add_argument("--mqtt-url", dest="mqtt_url", default="ssl://testrtls.kontakt.io:8083")
        parser.add_argument("--telemetry-types", dest="telemetry_types", nargs='+', default = ['sensor'])
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

    def get_mqtt_url(self):
        return self.args['mqtt_url']

    def get_telemetry_types(self):
        return self.args['telemetry_types']

    def get_streams_per_telegraf(self):
        return self.args['streams_per_telegraf']

DEFAULT_INFLUX_PASSWORD = 'aeg0UKOmUIzJap1QV6m8'

options = Options(sys.argv[1:])
api_client = ApiClient(options.get_api_url(), options.get_api_key())
influx_client = InfluxClient(options.get_influx_url(),
    options.get_influx_port(),
    options.get_influx_username(), 
    options.get_influx_password())

company_id = api_client.get_company_id()
influx_client.create_database(company_id)
influx_client.create_user(company_id, DEFAULT_INFLUX_PASSWORD, company_id)

influx_client.create_retention_policy(company_id, 'stream_rp', '3h')

influx_client.create_retention_policy(company_id, 'current_rp', '7d')
influx_client.recreate_continuous_query(company_id, '10s', 'current_rp', 'stream_rp', '10s', '20s')

influx_client.create_retention_policy(company_id, 'recent_rp', '30d')
influx_client.recreate_continuous_query(company_id, '5m', 'recent_rp', 'stream_rp', '5m', '5m')

influx_client.create_retention_policy(company_id, 'history_rp', '365d')
influx_client.recreate_continuous_query(company_id, '1h', 'history_rp', 'stream_rp', '1h', '5m')

unique_ids = api_client.get_telemetry_unique_ids(api_venue_id=options.get_api_venue_id())

topics = ['/stream/%s/%s' % elem for elem in 
    itertools.product(unique_ids, options.get_telemetry_types())]

topic_chunks = chunks(topics, options.get_streams_per_telegraf())

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
cfg.append_key_value('urls', ["%s:%d" % (options.get_influx_url(), options.get_influx_port())])
cfg.append_key_value('database', company_id)
cfg.append_key_value('username', options.get_influx_username())
cfg.append_key_value('password', options.get_influx_password())
cfg.append_key_value('precision', 's') # default
cfg.append_key_value('timeout', '5s') # default
cfg.append_key_value('retention_policy', 'stream_rp')

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastSingleClick')
cfg.append_key_value('out_field_name', 'singleClick')
cfg.append_key_value('tag_key', 'topic')
cfg.append_key_value('threshold', 30)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastSingleClick')
cfg.append_key_value('out_field_name', 'singleClick')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('threshold', 30)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastDoubleTap')
cfg.append_key_value('out_field_name', 'doubleTap')
cfg.append_key_value('tag_key', 'topic')
cfg.append_key_value('threshold', 30)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastDoubleTap')
cfg.append_key_value('out_field_name', 'doubleTap')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('threshold', 30)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastThreshold')
cfg.append_key_value('out_field_name', 'threshold')
cfg.append_key_value('tag_key', 'topic')
cfg.append_key_value('threshold', 30)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastThreshold')
cfg.append_key_value('out_field_name', 'threshold')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('threshold', 30)

cfg.append_section_name('processors.tagrename', True)
cfg.append_section_name('processors.tagrename.renames', inner=True)
cfg.append_key_value('topic', 'trackingId')

cfg.append_section_name('processors.regex', True)
cfg.append_section_name('processors.regex.tags', True)
cfg.append_key_value('key', 'topic')
cfg.append_key_value('pattern', '^\\\\/stream\\\\/([a-zA-Z0-9]+)\\\\/[a-z]+$')
cfg.append_key_value('replacement', '${1}')
cfg.append_section_name('processors.regex.tags', True)
cfg.append_key_value('key', 'trackingId')
cfg.append_key_value('pattern', '^\\\\/stream\\\\/([a-zA-Z0-9]+)\\\\/[a-z]+$')
cfg.append_key_value('replacement', '${1}')

cfg.append_section_name('processors.override', True)
cfg.append_key_value('name_override', 'telemetry')

cfg.append_section_name('inputs.mqtt_consumer', True)
cfg.append_key_value('servers', [options.get_mqtt_url()])
cfg.append_key_value('qos', 0) # default
cfg.append_key_value('connection_timeout', '300s')
cfg.append_key_value('persistent_session', False)
cfg.append_key_value('username', 'telegraf') # default
cfg.append_key_value('password', options.get_api_key())
cfg.append_key_value('data_format', 'json') # default
cfg.append_key_value('tag_keys', ['sourceId'])

idx = 0

for topic_chunk in topic_chunks:
    current = TelegrafConfigFormatter(cfg)
    current.append_key_value('topics', topic_chunk)

    f = open("telegraf.stream.conf.%d" % idx, "w")
    f.write(current.to_string())
    f.close()
    idx = idx + 1
