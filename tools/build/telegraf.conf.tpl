[agent]
  interval="5s"
  round_interval=false
  metric_buffer_limit=1000
  flush_buffer_when_full=true
  collection_jitter="0s"
  flush_interval="5s"
  flush_jitter=0
  omit_hostname=true
  debug={{ env.Getenv "DEBUG" "false" }}

[[inputs.http_listener_v2]]
  service_address = ":{{ env.Getenv "SERVICE_PORT" "8080" }}"
  path = "/event/collect"
  methods = ["POST"]
  data_format = "kontakt"
  header_tags = ["Api-Key"]

{{ if env.Getenv "HEALTH_PORT" "" -}}
[[outputs.health]]
  ## Address and port to listen on.
  service_address = "http://:{{ env.Getenv "HEALTH_PORT" 8080 }}"

  ## The maximum duration for reading the entire request.
  read_timeout = "{{ env.Getenv "HEALTH_READ_TIMEOUT" "5s" }}"
  ## The maximum duration for writing the entire response.
  write_timeout = "{{ env.Getenv "HEALTH_WRITE_TIMEOUT" "5s" }}"

  ## Username and password to accept for HTTP basic authentication.
  # basic_username = "{{ env.Getenv "HEALTH_BASIC_USERNAME" "health" }}"
  # basic_password = "{{ env.Getenv "HEALTH_BASIC_PASSWORD" "health" }}"

  ## Allowed CA certificates for client certificates.
  # tls_allowed_cacerts = ["/etc/telegraf/clientca.pem"]

  ## TLS server certificate and private key.
  # tls_cert = "{{ env.Getenv "HEALTH_TLS_CERT" "/etc/telegraf/cert.pem" }}"
  # tls_key  = "{{ env.Getenv "HEALTH_TLS_KEY" "/etc/telegraf/key.pem" }}"

{{ end -}}

{{ if env.Getenv "TELEMETRY_DATABASE_ENDPOINT" "" -}}
[[outputs.influxdb]]
  urls      = [ "{{ .Env.TELEMETRY_DATABASE_ENDPOINT }}" ]
  database  = "{{ env.Getenv "TELEMETRY_DATABASE_NAME" "telemetry" }}"
  username  = "{{ .Env.TELEMETRY_DATABASE_USERNAME }}"
  password  = "{{ .Env.TELEMETRY_DATABASE_PASSWORD }}"
  precision = "{{ env.Getenv "TELEMETRY_DATABASE_PRECISION" "s" }}"
  timeout   = "{{ env.Getenv "TELEMETRY_DATABASE_TIMEOUT" "5s" }}"
  retention_policy="stream_rp"
  taginclude = ["companyId", "trackingId"]
  fielddrop=["lastDoubleTap","lastSingleClick","lastThreshold","loggingEnabled","sensitivity","utcTime","secondsSinceThreshold","secondsSinceDoubleTap"]

{{ end -}}

{{ if env.Getenv "LOCATION_DATABASE_ENDPOINT" -}}
[[outputs.influxdb]]
  urls             = [ "{{ .Env.LOCATION_DATABASE_ENDPOINT }}" ]
  database         = "{{ env.Getenv "LOCATION_DATABASE_NAME" "le" }}"
  username         = "{{ .Env.LOCATION_DATABASE_USERNAME }}"
  password         = "{{ .Env.LOCATION_DATABASE_PASSWORD }}"
  precision        = "{{ env.Getenv "LOCATION_DATABASE_PRECISION" "s" }}"
  timeout          = "{{ env.Getenv "LOCATION_DATABASE_TIMEOUT" "10s" }}"
  retention_policy = "stream_rp"
  taginclude       = ["companyId", "trackingId", "sourceId"]
  fielddrop        = ["sourceId", "lastDoubleTap","lastSingleClick","lastThreshold","loggingEnabled","sensitivity","utcTime","secondsSinceThreshold","secondsSinceDoubleTap"]
{{ end -}}

{{ if env.Getenv "LOCATION_ENGINE_ENDPOINT" -}}
[[outputs.http]]
  url         = "{{ .Env.LOCATION_ENGINE_ENDPOINT }}"
  method      = "{{ env.Getenv "LOCATION_ENGINE_METHOD" "POST" }}"
  data_format = "{{ env.Getenv "LOCATION_ENGINE_DATA_FORMAT" "flat_json" }}"

  [outputs.http.headers]
    Content-Type = "application/json"
{{ end -}}

[[processors.rename]]
  order = 0
  [[processors.rename.replace]]
    tag  = "deviceAddress"
    dest = "trackingId"

{{ if env.Getenv "AUTH_API_ENDPOINT" -}}
[[processors.kontaktauth]]
  api_address = "{{ .Env.AUTH_API_ENDPOINT }}"
  order       = 1
{{ end -}}

[[processors.kioparser]]
  order = 2

[[processors.clickdetect]]
  field_name     = "clickId"
  order          = 3
  out_field_name = "singleClick"
  tag_key        = "trackingId"

[[processors.clickdetect]]
  field_name     = "movementId"
  order          = 4
  out_field_name = "movement"
  tag_key        = "trackingId"

[[processors.lastcalc]]
  field_name     = "lastSingleClick"
  order          = 5
  out_field_name = "singleClick"
  tag_key        = "trackingId"

[[processors.lastcalc]]
  field_name     = "lastThreshold"
  order          = 6
  out_field_name = "threshold"
  tag_key        = "trackingId"

[[processors.lastcalc]]
  field_name     = "lastDoubleTap"
  order          = 7
  out_field_name = "tap"
  tag_key        = "trackingId"

[[processors.regex]]
  order = 8
  [[processors.regex.fields]]
    key = "sourceId"
    pattern = "^(.*)$"
    replacement = "${1}"
    ## adding new field key
    result_key = "tSourceId"

[[processors.converter]]
  order = 9
  [processors.converter.fields]
    tag = ["tSourceId"]

[[processors.rename]]
  order = 10
  [[processors.rename.replace]]
    tag = "tSourceId"
    dest = "sourceId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
{{ end -}}