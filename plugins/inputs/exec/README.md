# Exec Input Plugin

The `exec` plugin executes all the `commands` in parallel on every interval and
parses metrics from their output in any one of the accepted [Input Data
Formats](../../../docs/DATA_FORMATS_INPUT.md).

This plugin can be used to poll for custom metrics from any source.

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Read metrics from one or more commands that can output to stdout
[[inputs.exec]]
  ## Commands array
  commands = [
    "/tmp/test.sh",
    "/usr/bin/mycollector --foo=bar",
    "/tmp/collect_*.sh"
  ]

  ## Environment variables
  ## Array of "key=value" pairs to pass as environment variables
  ## e.g. "KEY=value", "USERNAME=John Doe",
  ## "LD_LIBRARY_PATH=/opt/custom/lib64:/usr/local/libs"
  # environment = []

  ## Timeout for each command to complete.
  timeout = "5s"

  ## measurement name suffix (for separating different commands)
  name_suffix = "_mycollector"

  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
```

Glob patterns in the `command` option are matched on every run, so adding new
scripts that match the pattern will cause them to be picked up immediately.

## Example

This script produces static values, since no timestamp is specified the values
are at the current time. Ensure that int values are followed with `i` for proper
parsing.

### Example:

This script produces static values, since no timestamp is specified the values are at the current time.
```sh
#!/bin/sh
echo 'example,tag1=a,tag2=b i=42i,j=43i,k=44i'
```

It can be paired with the following configuration and will be run at the
`interval` of the agent.

It can be paired with the following configuration and will be run at the `interval` of the agent.
```toml
[[inputs.exec]]
  commands = ["sh /tmp/test.sh"]
  timeout = "5s"
  data_format = "influx"
```

## Common Issues

### My script works when I run it by hand, but not when Telegraf is running as a service

This may be related to the Telegraf service running as a different user. The
official packages run Telegraf as the `telegraf` user and group on Linux
systems.

### With a PowerShell on Windows, the output of the script appears to be truncated

You may need to set a variable in your script to increase the number of columns
available for output:

```shell
$host.UI.RawUI.BufferSize = new-object System.Management.Automation.Host.Size(1024,50)
```
### Common Issues:

#### Q: My script works when I run it by hand, but not when Telegraf is running as a service.

This may be related to the Telegraf service running as a different user.  The
official packages run Telegraf as the `telegraf` user and group on Linux
systems.
