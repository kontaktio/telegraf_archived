#!/usr/bin/env bash
REGION=eu-west-1

eval $(aws ecr get-login --region us-east-1 --no-include-email)
aws s3 cp s3://kontakt-telegraf-config/build-prod/env . --source-region us-east-1
source ./env

API_KEY=$(curl -s -X POST "https://api.kontakt.io/manager/impersonate" \
	-H "Api-Key: ${SERVICE_API_KEY}" \
    -H "Accept: application/vnd.com.kontakt+json;version=10" \
    -H "Content-Type: application/x-www-form-urlencoded; c \
harset=utf-8" \
	--data-urlencode "email=$EMAIL" | jq .apiKey | tr -d '"')
EMAIL_CLEAN=`echo $EMAIL | sed 's/@/-at-/g' | tr -cd '[:alnum:]-'`
LABEL="telegraf-$EMAIL_CLEAN"

echo "
{
  \"family\": \"$LABEL\",
  \"taskRoleArn\": \"arn:aws:iam::712852996757:role/telegraf-role\",
  \"containerDefinitions\": [
    {
      \"image\": \"712852996757.dkr.ecr.eu-west-1.amazonaws.com/telegraf-deployment:latest\",
      \"name\": \"$LABEL-container\",
      \"memoryReservation\": 256,
      \"cpu\": 0,
      \"portMappings\": [],
      \"logConfiguration\": {
        \"logDriver\": \"awslogs\",
        \"options\": {
          \"awslogs-region\": \"$REGION\",
          \"awslogs-stream-prefix\": \"$LABEL\",
          \"awslogs-group\": \"/ecs/telegraf-deployment\"
        }
      },
      \"environment\": [
        {
          \"name\": \"API_KEY\",
          \"value\": \"$API_KEY\"
        },
        {
          \"name\": \"VENUE_ID\",
          \"value\": \"$VENUE_ID\"
        }
      ]
    }
  ]
}
" > taskdef.json

aws ecs register-task-definition --cli-input-json "file://taskdef.json" --region $REGION
aws ecs create-service --cluster "prod-ecs-cluster" --service-name "$LABEL" --region $REGION --task-definition "$LABEL" --desired-count 1 --placement-strategy type=spread,field=instanceId
