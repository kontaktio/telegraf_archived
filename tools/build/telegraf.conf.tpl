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

{{ if env.Getenv "DEBUG" "" -}}
[[outputs.file]]
  files       = ["stdout", "/tmp/raw_metrics.out"]
  data_format = "influx"
{{ end -}}

{{ if env.Getenv "DEBUG" "" -}}
[[outputs.file]]
  files       = ["stdout", "/tmp/filtered_metrics.out"]
  taginclude  = ["companyId", "trackingId", "sourceId"]
  fielddrop   = ["sourceId", "lastDoubleTap","lastSingleClick","lastThreshold","loggingEnabled","sensitivity","utcTime","secondsSinceThreshold","secondsSinceDoubleTap"]
  data_format = "influx"
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

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 0
{{ end -}}

[[processors.rename]]
  order = 1
  [[processors.rename.replace]]
    tag  = "deviceAddress"
    dest = "trackingId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 2
{{ end -}}

{{ if env.Getenv "AUTH_API_ENDPOINT" -}}
[[processors.kontaktauth]]
  api_address = "{{ .Env.AUTH_API_ENDPOINT }}"
  order       = 3
{{ end -}}

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 4
{{ end -}}

[[processors.kioparser]]
  order = 5

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 6
{{ end -}}

[[processors.clickdetect]]
  order          = 7
  field_name     = "clickId"
  out_field_name = "singleClick"
  tag_key        = "trackingId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 8
{{ end -}}

[[processors.clickdetect]]
  field_name     = "movementId"
  order          = 9
  out_field_name = "movement"
  tag_key        = "trackingId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 10
{{ end -}}

[[processors.lastcalc]]
  field_name     = "lastSingleClick"
  order          = 11
  out_field_name = "singleClick"
  tag_key        = "trackingId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 12
{{ end -}}

[[processors.lastcalc]]
  field_name     = "lastThreshold"
  order          = 13
  out_field_name = "threshold"
  tag_key        = "trackingId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 14
{{ end -}}

[[processors.lastcalc]]
  field_name     = "lastDoubleTap"
  order          = 15
  out_field_name = "tap"
  tag_key        = "trackingId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 16
{{ end -}}

[[processors.regex]]
  order = 17
  [[processors.regex.fields]]
    key = "sourceId"
    pattern = "^(.*)$"
    replacement = "${1}"
    ## adding new field key
    result_key = "tSourceId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 18
{{ end -}}

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order =19
{{ end -}}

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 20
{{ end -}}

[[processors.converter]]
  order = 21
  [processors.converter.fields]
    tag = ["tSourceId"]

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 22
{{ end -}}

[[processors.rename]]
  order = 23
  [[processors.rename.replace]]
    tag = "tSourceId"
    dest = "sourceId"

{{ if env.Getenv "DEBUG" "" -}}
[[processors.printer]]
  order = 24
{{ end -}}