# Swap Input Plugin

The swap plugin collects system swap metrics.

For more information on what swap memory is, read [All about Linux swap
space](https://www.linux.com/news/all-about-linux-swap-space).

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Read metrics about swap memory usage
[[inputs.swap]]
  # no configuration
```

## Metrics
### Metrics:

- swap
  - fields:
    - free (int, bytes): free swap memory
    - total (int, bytes): total swap memory
    - used (int, bytes): used swap memory
    - used_percent (float, percent): percentage of swap memory used
    - in (int, bytes): data swapped in since last boot calculated from page number
    - out (int, bytes): data swapped out since last boot calculated from page number

## Example Output

```shell
### Example Output:

```
swap total=20855394304i,used_percent=45.43883523785713,used=9476448256i,free=1715331072i 1511894782000000000
```
