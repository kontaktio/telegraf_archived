# Teamspeak 3 Input Plugin

This plugin uses the Teamspeak 3 ServerQuery interface of the Teamspeak server
to collect statistics of one or more virtual servers. If you are querying an
external Teamspeak server, make sure to add the host which is running Telegraf
to query_ip_allowlist.txt in the Teamspeak Server directory. For information
about how to configure the server take a look the [Teamspeak 3 ServerQuery
Manual][1].

[1]: http://media.teamspeak.com/ts3_literature/TeamSpeak%203%20Server%20Query%20Manual.pdf

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Reads metrics from a Teamspeak 3 Server via ServerQuery
[[inputs.teamspeak]]
  ## Server address for Teamspeak 3 ServerQuery
  # server = "127.0.0.1:10011"
  ## Username for ServerQuery
  username = "serverqueryuser"
  ## Password for ServerQuery
  password = "secret"
  ## Nickname of the ServerQuery client
  nickname = "telegraf"
  ## Array of virtual servers
  # virtual_servers = [1]
```
# Reads metrics from a Teamspeak 3 Server via ServerQuery
[[inputs.teamspeak]]
  ## Server address for Teamspeak 3 ServerQuery
  # server = "127.0.0.1:10011"
  ## Username for ServerQuery
  username = "serverqueryuser"
  ## Password for ServerQuery
  password = "secret"
  ## Nickname of the ServerQuery client
  nickname = "telegraf"
  ## Array of virtual servers
  # virtual_servers = [1]
```

## Metrics

- teamspeak
  - uptime
  - clients_online
  - total_ping
  - total_packet_loss
  - packets_sent_total
  - packets_received_total
  - bytes_sent_total
  - bytes_received_total
  - query_clients_online

### Tags

- The following tags are used:
  - virtual_server
  - name

## Example Output

```shell
teamspeak,virtual_server=1,name=LeopoldsServer,host=vm01 bytes_received_total=29638202639i,uptime=13567846i,total_ping=26.89,total_packet_loss=0,packets_sent_total=415821252i,packets_received_total=237069900i,bytes_sent_total=55309568252i,clients_online=11i,query_clients_online=1i 1507406561000000000
```
### Measurements:

- teamspeak
    - uptime
    - clients_online
    - total_ping
    - total_packet_loss
    - packets_sent_total
    - packets_received_total
    - bytes_sent_total
    - bytes_received_total

### Tags:

- The following tags are used:
    - virtual_server
    - name

### Example output:

```
teamspeak,virtual_server=1,name=LeopoldsServer,host=vm01 bytes_received_total=29638202639i,uptime=13567846i,total_ping=26.89,total_packet_loss=0,packets_sent_total=415821252i,packets_received_total=237069900i,bytes_sent_total=55309568252i,clients_online=11i 1507406561000000000
```