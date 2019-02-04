#!/usr/bin/env sh
REGION=us-east-1
TASKDEF=taskdef-accept.json
CREDENTIALS=influx-credentials.json


eval $(aws ecr get-login --region us-east-1 --no-include-email)
aws s3 cp s3://kontakt-telegraf-config/"${CREDENTIALS}" . --source-region us-east-1

INFLUX_USERNAME=$(cat ${CREDENTIALS} | jq .username | tr -d '"')
INFLUX_PASSWORD=$(cat ${CREDENTIALS} | jq .password | tr -d '"')

LABEL_REPLACE_PATTERN="s;%LABEL%;$LABEL;g"
APIKEY_REPLACE_PATTERN="s;%API_KEY%;$API_KEY;g"
INFLUXUSERNAME_REPLACE_PATTERN="s;%INFLUXDB_USERNAME%;$INFLUX_USERNAME;g"
INFLUXPASS_REPLACE_PATTERN="s;%INFLUXDB_PASSWORD%;$INFLUX_PASSWORD;g"
VENUEID_REPLACE_PATTERN="s;%VENUE_ID%;$VENUE_ID;g"

aws s3 cp s3://kontakt-telegraf-config/"${TASKDEF}" . --source-region us-east-1
export TASKDEF_TMP=$(cat ${TASKDEF} | tr -d '\r\n' | tr -d '\t')
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $LABEL_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $APIKEY_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $INFLUXUSERNAME_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $INFLUXPASS_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $VENUEID_REPLACE_PATTERN)

#sed -e "s;%LABEL%;${LABEL};g" $TASKDEF | dd of="${TASKDEF}.tmp"
#sed -e "s;%API_KEY%;${API_KEY};g" "${TASKDEF}.tmp" | dd of="${TASKDEF}.tmp"
#sed -e "s;%INFLUXDB_USERNAME%;${INFLUX_USERNAME};g" "${TASKDEF}.tmp" | dd of="${TASKDEF}.tmp"
#sed -e "s;%INFLUXDB_PASSWORD%;${INFLUX_PASSWORD};g" "${TASKDEF}.tmp" | dd of="${TASKDEF}.tmp"
#sed -e "s;%VENUE_ID%;${VENUE_ID};g" "${TASKDEF}.tmp" | dd of="${TASKDEF}.tmp"

echo $TASKDEF_TMP | dd of="${TASKDEF}.tmp"

TASK_DEFINITION_REVISION=$(aws ecs register-task-definition --cli-input-json "file://${TASKDEF}.tmp" --region $REGION | jq .taskDefinition.revision | tr -d '"')
aws ecs create-service --cluster "accept-ecs-cluster" --service-name "$LABEL-telegraf-service" --task-definition "$LABEL:$TASK_DEFINITION_REVISION" --desired-count 1 --region $REGION
