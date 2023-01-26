# Memory Input Plugin
# Mem Input Plugin

The mem plugin collects system memory metrics.

For a more complete explanation of the difference between *used* and
*actual_used* RAM, see [Linux ate my ram](http://www.linuxatemyram.com/).

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Read metrics about memory usage
[[inputs.mem]]
  # no configuration
```

## Metrics
### Metrics:

Available fields are dependent on platform.

- mem
  - fields:
    - active (integer, Darwin, FreeBSD, Linux, OpenBSD)
    - available (integer)
    - available_percent (float)
    - buffered (integer, FreeBSD, Linux)
    - cached (integer, FreeBSD, Linux, OpenBSD)
    - commit_limit (integer, Linux)
    - committed_as (integer, Linux)
    - dirty (integer, Linux)
    - free (integer, Darwin, FreeBSD, Linux, OpenBSD)
    - high_free (integer, Linux)
    - high_total (integer, Linux)
    - huge_pages_free (integer, Linux)
    - huge_page_size (integer, Linux)
    - huge_pages_total (integer, Linux)
    - inactive (integer, Darwin, FreeBSD, Linux, OpenBSD)
    - laundry (integer, FreeBSD)
    - low_free (integer, Linux)
    - low_total (integer, Linux)
    - mapped (integer, Linux)
    - page_tables (integer, Linux)
    - shared (integer, Linux)
    - slab (integer, Linux)
    - sreclaimable (integer, Linux)
    - sunreclaim (integer, Linux)
    - swap_cached (integer, Linux)
    - swap_free (integer, Linux)
    - swap_total (integer, Linux)
    - total (integer)
    - used (integer)
    - used_percent (float)
    - vmalloc_chunk (integer, Linux)
    - vmalloc_total (integer, Linux)
    - vmalloc_used (integer, Linux)
    - wired (integer, Darwin, FreeBSD, OpenBSD)
    - write_back (integer, Linux)
    - write_back_tmp (integer, Linux)

## Example Output

```shell
mem active=9299595264i,available=16818249728i,available_percent=80.41654254645131,buffered=2383761408i,cached=13316689920i,commit_limit=14751920128i,committed_as=11781156864i,dirty=122880i,free=1877688320i,high_free=0i,high_total=0i,huge_page_size=2097152i,huge_pages_free=0i,huge_pages_total=0i,inactive=7549939712i,low_free=0i,low_total=0i,mapped=416763904i,page_tables=19787776i,shared=670679040i,slab=2081071104i,sreclaimable=1923395584i,sunreclaim=157675520i,swap_cached=1302528i,swap_free=4286128128i,swap_total=4294963200i,total=20913917952i,used=3335778304i,used_percent=15.95004011996231,vmalloc_chunk=0i,vmalloc_total=35184372087808i,vmalloc_used=0i,wired=0i,write_back=0i,write_back_tmp=0i 1574712869000000000
    - active (integer)
    - available (integer)
    - buffered (integer)
    - cached (integer)
    - free (integer)
    - inactive (integer)
    - slab (integer)
    - total (integer)
    - used (integer)
    - available_percent (float)
    - used_percent (float)
    - wired (integer)
    - commit_limit (integer)
    - committed_as (integer)
    - dirty (integer)
    - high_free (integer)
    - high_total (integer)
    - huge_page_size (integer)
    - huge_pages_free (integer)
    - huge_pages_total (integer)
    - low_free (integer)
    - low_total (integer)
    - mapped (integer)
    - page_tables (integer)
    - shared (integer)
    - swap_cached (integer)
    - swap_free (integer)
    - swap_total (integer)
    - vmalloc_chunk (integer)
    - vmalloc_total (integer)
    - vmalloc_used (integer)
    - write_back (integer)
    - write_back_tmp (integer)

### Example Output:
```
mem active=11347566592i,available=18705133568i,available_percent=89.4288960571006,buffered=1976709120i,cached=13975572480i,commit_limit=14753067008i,committed_as=2872422400i,dirty=87461888i,free=1352400896i,high_free=0i,high_total=0i,huge_page_size=2097152i,huge_pages_free=0i,huge_pages_total=0i,inactive=6201593856i,low_free=0i,low_total=0i,mapped=310427648i,page_tables=14397440i,shared=200781824i,slab=1937526784i,swap_cached=0i,swap_free=4294963200i,swap_total=4294963200i,total=20916207616i,used=3611525120i,used_percent=17.26663449848977,vmalloc_chunk=0i,vmalloc_total=35184372087808i,vmalloc_used=0i,wired=0i,write_back=0i,write_back_tmp=0i 1536704085000000000
```
