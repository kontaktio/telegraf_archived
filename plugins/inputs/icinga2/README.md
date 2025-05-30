# Icinga2 Input Plugin

This plugin gather services & hosts status using Icinga2 Remote API.

The icinga2 plugin uses the icinga2 remote API to gather status on running
services and hosts. You can read Icinga2's documentation for their remote API
[here][1].

[1]: https://docs.icinga.com/icinga2/latest/doc/module/icinga2/chapter/icinga2-api

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Gather Icinga2 status
[[inputs.icinga2]]
  ## Required Icinga2 server address
  # server = "https://localhost:5665"
  
  ## Required Icinga2 object type ("services" or "hosts")
  # object_type = "services"

  ## Credentials for basic HTTP authentication
  # username = "admin"
  # password = "admin"

  ## Maximum time to receive response.
  # response_timeout = "5s"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = true
```

## Measurements & Fields

- All measurements have the following fields:
  - name (string)
  - state_code (int)

## Tags

- All measurements have the following tags:
  - check_command - The short name of the check command
  - display_name - The name of the service or host
  - state - The state: UP/DOWN for hosts, OK/WARNING/CRITICAL/UNKNOWN for services
  - source - The icinga2 host
  - port - The icinga2 port
  - scheme - The icinga2 protocol (http/https)
  - server - The server the check_command is running for

## Sample Queries

```sql
### Measurements & Fields:

- All measurements have the following fields:
    - name (string)
    - state_code (int)

### Tags:

- All measurements have the following tags:
    - check_command
    - display_name
    - state
    - source
    - port
    - scheme

### Sample Queries:

```
SELECT * FROM "icinga2_services" WHERE state_code = 0 AND time > now() - 24h // Service with OK status
SELECT * FROM "icinga2_services" WHERE state_code = 1 AND time > now() - 24h // Service with WARNING status
SELECT * FROM "icinga2_services" WHERE state_code = 2 AND time > now() - 24h // Service with CRITICAL status
SELECT * FROM "icinga2_services" WHERE state_code = 3 AND time > now() - 24h // Service with UNKNOWN status
```

## Example Output

```text
### Example Output:

```
$ ./telegraf -config telegraf.conf -input-filter icinga2 -test
icinga2_hosts,display_name=router-fr.eqx.fr,check_command=hostalive-custom,host=test-vm,source=localhost,port=5665,scheme=https,state=ok name="router-fr.eqx.fr",state=0 1492021603000000000
```
