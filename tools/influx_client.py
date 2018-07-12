from influxdb import InfluxDBClient


class InfluxClient:
    READ_PRIVILEGE = 'read'
    CONTINUOUS_QUERY_FMT = """
CREATE CONTINUOUS QUERY "telemetry_{0}_cq" ON "{4}"
RESAMPLE EVERY {3}
BEGIN
    SELECT 
        mean("batteryLevel") AS "batteryLevel", 
        mean("lightLevel") AS "lightLevel", 
        mean("rssi") AS "rssi", 
        mean("sensitivity") AS "sensitivity", 
        sum("singleClick") AS "singleClick", 
        sum("threshold") AS "threshold", 
        sum("doubleTap") AS "doubleTap", 
        mean("temperature") AS "temperature", 
        mean("x") AS "x", 
        mean("y") AS "y", 
        mean("z") AS "z"
    INTO 
        "{1}"."telemetry_{0}"
    FROM 
       "{2}"."telemetry"
    GROUP BY time({0}), *
END
"""

    def __init__(self, address, port, user_name, password):
        self._client = InfluxDBClient(
            host=address,
            port=port,
            username=user_name, 
            password=password)

    def create_database(self, database_name):
        print "Creating database %s" % database_name
        self._client.create_database(database_name)
        
    def create_user(self, user_name, password, database_name=None):
        print "Creating user %s" % user_name
        self._client.create_user(user_name, password)
        if database_name is not None:
            self._client.grant_privilege(self.READ_PRIVILEGE, database_name, user_name)

    def create_retention_policy(self, database_name, policy_name, duration):
        print "Creating retention policy %s with duration %s on database %s" % (policy_name, duration, database_name)
        self._client.create_retention_policy(policy_name, duration, 1, database=database_name)

    def create_continuous_query(self, database_name, aggregation_time, retention_policy, source_retention_policy, resample_time):
        q = self.CONTINUOUS_QUERY_FMT.format(aggregation_time, retention_policy, source_retention_policy, resample_time, database_name)
        print "Executing query %s" % q
        self._client.query(q, database=database_name)
    