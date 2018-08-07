import requests
import json

class KapacitorClient():
    TASKS = '%s/kapacitor/v1/tasks'
    TASK = '%s/kapacitor/v1/tasks/%s'

    def __init__(self, address, database, retention_policy):
        self._address = address
        self._database = database
        self._retention_policy = retention_policy

    def list_tasks(self):
        return requests.get(self.TASKS % self._address).json()
    
    def get_task_info(self, task_id):
        return requests.get(self.TASK % (self._address, task_id)).json()

    def create_task(self, id, template, parameters={}):
        payload = {
            'id': id,
            'type': 'stream',
            'dbrps': [{'db': self._database, 'rp': self._retention_policy}],
            'vars': parameters,
            'template-id': template,
            'status': 'enabled'
        }

        result = requests.post(self.TASKS % self._address, json=payload)
        return result.json()
