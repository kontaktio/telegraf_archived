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
        parser.add_argument("--influxdb-url", dest="influxdb_url", required=True)
        parser.add_argument("--influxdb-port", dest="influxdb_port", default=8086, type=int)
        parser.add_argument("--influxdb-username", dest="influxdb_username", required=True)
        parser.add_argument("--influxdb-password", dest="influxdb_password", required=True)
        parser.add_argument("--influxdb-dbname", dest="influxdb_dbname", default="telemetry")
        parser.add_argument("--config-file", dest="config_file")

        self.args = vars(parser.parse_args(args=args))
        if self.args['config_file'] is not None:
            self._load_from_file()

    def _load_from_file(self):
        parser = ConfigParser.ConfigParser()
        parser.read(self.args['config_file'])
        values = parser.options(self.CONFIG_FILE_SECTION)
        for k in values:
            self.args[k] = parser.get(self.CONFIG_FILE_SECTION, k)

    def get_influx_url(self):
        return self.args['influxdb_url']

    def get_influx_port(self):
        return self.args['influxdb_port']

    def get_influx_username(self):
        return self.args['influxdb_username']

    def get_influx_password(self):
        return self.args['influxdb_password']

    def get_influx_dbname(self):
        return self.args['influxdb_dbname']


options = Options(sys.argv[1:])

cfg = TelegrafConfigFormatter()

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

cfg.append_section_name('inputs.socket_listener', True)
cfg.append_key_value('service_address', 'unix:///tmp/telegraf.sock')
cfg.append_key_value('data_format', 'json')
cfg.append_key_value('tag_keys', ['deviceAddress', 'companyId'])

cfg.append_section_name('outputs.influxdb', True)
cfg.append_key_value('urls', ["%s:%d" % (options.get_influx_url(), options.get_influx_port())])
cfg.append_key_value('database', options.get_influx_dbname())
cfg.append_key_value('username', options.get_influx_username())
cfg.append_key_value('password', options.get_influx_password())
cfg.append_key_value('precision', 's')
cfg.append_key_value('timeout', '5s')
cfg.append_key_value('retention_policy', 'stream_rp')
cfg.append_key_value('tagexclude', ['topic'])
cfg.append_key_value('fielddrop', ['lastDoubleTap', 'lastSingleClick', 'lastThreshold', 'loggingEnabled', 'sensitivity',
                                   'utcTime', 'secondsSinceThreshold', 'secondsSinceDoubleTap'])

cfg.append_section_name('processors.kioparser', True)
cfg.append_key_value('order', 0)

cfg.append_section_name('processors.rename', True)
cfg.append_key_value('order', 1)
cfg.append_section_name('processors.rename.replace', True, True)
cfg.append_key_value('tag', 'deviceAddress')
cfg.append_key_value('dest', 'trackingId')

cfg.append_section_name('processors.clickdetect', True)
cfg.append_key_value('field_name', 'clickId')
cfg.append_key_value('out_field_name', 'singleClick')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('order', 2)

cfg.append_section_name('processors.clickdetect', True)
cfg.append_key_value('field_name', 'movementId')
cfg.append_key_value('out_field_name', 'movement')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('order', 3)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastSingleClick')
cfg.append_key_value('out_field_name', 'singleClick')
cfg.append_key_value('threshold', 60)
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('order', 4)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastDoubleTap')
cfg.append_key_value('out_field_name', 'doubleTap')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('threshold', 60)
cfg.append_key_value('order', 5)

cfg.append_section_name('processors.lastcalc', True)
cfg.append_key_value('field_name', 'lastThreshold')
cfg.append_key_value('out_field_name', 'threshold')
cfg.append_key_value('tag_key', 'trackingId')
cfg.append_key_value('threshold', 60)
cfg.append_key_value('order', 6)

cfg.append_section_name('processors.override', True)
cfg.append_key_value('name_override', 'telemetry')


f = open("telegraf.events.conf", "w")
formatter = TelegrafConfigFormatter(cfg)
f.write(formatter.to_string())
f.close()