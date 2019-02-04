import argparse
import ConfigParser
import sys

from api_client import ApiClient
from kapacitor_client import KapacitorClient

DWELL_TIME_TASK_NAME = "dwellTimeReport_%s"


class Options(object):
    CONFIG_FILE_SECTION = 'default'

    def __init__(self, args):
        parser = argparse.ArgumentParser(
            description='Create Kapacitor tasks for reports module',
            formatter_class=argparse.ArgumentDefaultsHelpFormatter)
        parser.add_argument('--api-key', dest='api_key', required=True)
        parser.add_argument('--kapacitor-url', dest='kapacitor_url', required=True)
        parser.add_argument('--kapacitor-user', dest='kapacitor_user', required=True)
        parser.add_argument('--kapacitor-pass', dest='kapacitor_pass', required=True)
        parser.add_argument('--config-file', dest='config_file')
        parser.add_argument('--api-url', dest='api_url', default='https://testapi.kontakt.io/')
        parser.add_argument('--debug', dest='debug', default=False, type=bool)
        parser.add_argument('--log-file', dest='log_file', default='~/telegraf.log')

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

    def get_log_file(self):
        return self.args['log_file']

    def get_kapacitor_url(self):
        return self.args['kapacitor_url']

    def get_kapacitor_user(self):
        return self.args['kapacitor_user']

    def get_kapacitor_pass(self):
        return self.args['kapacitor_pass']


options = Options(sys.argv[1:])
api_client = ApiClient(options.get_api_url(), options.get_api_key())

company_id = api_client.get_company_id()
kapacitor_client = KapacitorClient(options, company_id, 'stream_rp')

dwell_time_task_name = DWELL_TIME_TASK_NAME % company_id
kapacitor_client.remove_task(dwell_time_task_name)
result = kapacitor_client.create_task(dwell_time_task_name, 'dwellTime-tpl', {
    'database': {
        'value': company_id,
        'type': 'string'
    }
})
if len(result['error']) > 0:
    raise(Exception(result['error']))
