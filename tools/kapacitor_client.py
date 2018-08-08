import requests
import json

class KapacitorClient():
    TASKS = '%s/kapacitor/v1/tasks'
    TASK = '%s/kapacitor/v1/tasks/%s'

    def __init__(self, options, database, retention_policy):
        self._address = options.get_kapacitor_url()
        self._username = options.get_kapacitor_user()
        self._password = options.get_kapacitor_pass()
        self._database = database
        self._retention_policy = retention_policy

    def list_tasks(self):
        return requests.get(self.TASKS % self._address, auth=self._auth()).json()
    
    def get_task_info(self, task_id):
        return requests.get(self.TASK % (self._address, task_id), auth=self._auth()).json()

    def create_task(self, id, template, parameters={}):
        payload = {
            'id': id,
            'type': 'stream',
            'dbrps': [{'db': self._database, 'rp': self._retention_policy}],
            'vars': parameters,
            'template-id': template,
            'status': 'enabled'
        }

        result = requests.post(self.TASKS % self._address, json=payload, auth=self._auth())
        return result.json()

    def _auth(self):
        return (self._username, self._password)
