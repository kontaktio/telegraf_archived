from influxdb import InfluxDBClient


class InfluxClient:
    READ_PRIVILEGE = 'read'

    def __init__(self, address, user_name, password):
        self._client = InfluxDBClient(
            host=address.replace('http://', ''), 
            username=user_name, 
            password=password)

    def create_database(self, database_name):
        self._client.create_database(database_name)
        
    def create_user(self, user_name, password, database_name=None):
        self._client.create_user(user_name, password)
        if database_name is not None:
            self._client.grant_privilege(self.READ_PRIVILEGE, database_name, user_name)
    