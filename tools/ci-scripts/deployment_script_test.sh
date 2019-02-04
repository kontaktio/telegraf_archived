#!/usr/bin/env sh
REGION=us-east-1
TASKDEF=taskdef-test.json
INFLUX_CREDENTIALS=influx-credentials.json
KAPACITOR_CREDENTIALS=kapacitor-credentials.json


eval $(aws ecr get-login --region us-east-1 --no-include-email)
aws s3 cp s3://kontakt-telegraf-config/"${INFLUX_CREDENTIALS}" . --source-region us-east-1
aws s3 cp s3://kontakt-telegraf-config/"${KAPACITOR_CREDENTIALS}" . --source-region us-east-1

INFLUX_USERNAME=$(cat ${INFLUX_CREDENTIALS} | jq .username | tr -d '"')
INFLUX_PASSWORD=$(cat ${INFLUX_CREDENTIALS} | jq .password | tr -d '"')

KAPACITOR_USER=$(cat ${KAPACITOR_CREDENTIALS} | jq .username | tr -d '"')
KAPACITOR_PASS=$(cat ${KAPACITOR_CREDENTIALS} | jq .password | tr -d '"')

LABEL_REPLACE_PATTERN="s;%LABEL%;$LABEL;g"
APIKEY_REPLACE_PATTERN="s;%API_KEY%;$API_KEY;g"
INFLUXUSERNAME_REPLACE_PATTERN="s;%INFLUXDB_USERNAME%;$INFLUX_USERNAME;g"
INFLUXPASS_REPLACE_PATTERN="s;%INFLUXDB_PASSWORD%;$INFLUX_PASSWORD;g"
VENUEID_REPLACE_PATTERN="s;%VENUE_ID%;$VENUE_ID;g"
KAPACITOR_USER_REPLACE_PATTERN="s;%KAPACITOR_USER%;$KAPACITOR_USER;g"
KAPACITOR_PASS_REPLACE_PATTERN="s;%KAPACITOR_PASS%;$KAPACITOR_PASS;g"

aws s3 cp s3://kontakt-telegraf-config/"${TASKDEF}" . --source-region us-east-1
export TASKDEF_TMP=$(cat ${TASKDEF} | tr -d '\r\n' | tr -d '\t')
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $LABEL_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $APIKEY_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $INFLUXUSERNAME_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $INFLUXPASS_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $VENUEID_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $KAPACITOR_USER_REPLACE_PATTERN)
export TASKDEF_TMP=$(echo $TASKDEF_TMP | sed -e $KAPACITOR_PASS_REPLACE_PATTERN)

echo $TASKDEF_TMP | dd of="${TASKDEF}.tmp"

TASK_DEFINITION_REVISION=$(aws ecs register-task-definition --cli-input-json "file://${TASKDEF}.tmp" --region $REGION | jq .taskDefinition.revision | tr -d '"')
aws ecs create-service --cluster "test-ecs-cluster" --service-name "$LABEL-telegraf-service" --task-definition "$LABEL:$TASK_DEFINITION_REVISION" --desired-count 1 --region $REGION
