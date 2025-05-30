# CPU Input Plugin

The `cpu` plugin gather metrics on the system CPUs.

## macOS Support

The [gopsutil][1] library, which is used to collect CPU data, does not support
gathering CPU metrics without CGO on macOS. The user will see a "not
implemented" message in this case. Builds provided by InfluxData do not build
with CGO.

Users can use the builds provided by [Homebrew][2], which build with CGO, to
produce CPU metrics.

[1]: https://github.com/shirou/gopsutil/blob/master/cpu/cpu_darwin_nocgo.go
[2]: https://formulae.brew.sh/formula/telegraf

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Read metrics about cpu usage
[[inputs.cpu]]
  ## Whether to report per-cpu stats or not
  percpu = true
  ## Whether to report total system cpu stats or not
  totalcpu = true
  ## If true, collect raw CPU time metrics
  collect_cpu_time = false
  ## If true, compute and report the sum of all non-idle CPU states
  report_active = false
  ## If true and the info is available then add core_id and physical_id tags
  core_tags = false
```
[[inputs.cpu]]
  ## Whether to report per-cpu stats or not
  percpu = true
  ## Whether to report total system cpu stats or not
  totalcpu = true
  ## If true, collect raw CPU time metrics
  collect_cpu_time = false
  ## If true, compute and report the sum of all non-idle CPU states
  report_active = false
  ## If true and the info is available then add core_id and physical_id tags
  core_tags = false
```

## Metrics

On Linux, consult `man proc` for details on the meanings of these values.

- cpu
  - tags:
    - cpu (CPU ID or `cpu-total`)
  - fields:
    - time_user (float)
    - time_system (float)
    - time_idle (float)
    - time_active (float)
    - time_nice (float)
    - time_iowait (float)
    - time_irq (float)
    - time_softirq (float)
    - time_steal (float)
    - time_guest (float)
    - time_guest_nice (float)
    - usage_user (float, percent)
    - usage_system (float, percent)
    - usage_idle (float, percent)
    - usage_active (float)
    - usage_nice (float, percent)
    - usage_iowait (float, percent)
    - usage_irq (float, percent)
    - usage_softirq (float, percent)
    - usage_steal (float, percent)
    - usage_guest (float, percent)
    - usage_guest_nice (float, percent)

## Troubleshooting

On Linux systems the `/proc/stat` file is used to gather CPU times.
Percentages are based on the last 2 samples.
Tags core_id and physical_id are read from `/proc/cpuinfo` on Linux systems

## Example Output

```shell
cpu,cpu=cpu0,host=loaner time_active=202224.15999999992,time_guest=30250.35,time_guest_nice=0,time_idle=1527035.04,time_iowait=1352,time_irq=0,time_nice=169.28,time_softirq=6281.4,time_steal=0,time_system=40097.14,time_user=154324.34 1568760922000000000
cpu,cpu=cpu0,host=loaner usage_active=31.249999981810106,usage_guest=2.083333333080696,usage_guest_nice=0,usage_idle=68.7500000181899,usage_iowait=0,usage_irq=0,usage_nice=0,usage_softirq=0,usage_steal=0,usage_system=4.166666666161392,usage_user=25.000000002273737 1568760922000000000
cpu,cpu=cpu1,host=loaner time_active=201890.02000000002,time_guest=30508.41,time_guest_nice=0,time_idle=264641.18,time_iowait=210.44,time_irq=0,time_nice=181.75,time_softirq=4537.88,time_steal=0,time_system=39480.7,time_user=157479.25 1568760922000000000
cpu,cpu=cpu1,host=loaner usage_active=12.500000010610771,usage_guest=2.0833333328280585,usage_guest_nice=0,usage_idle=87.49999998938922,usage_iowait=0,usage_irq=0,usage_nice=0,usage_softirq=2.0833333332070145,usage_steal=0,usage_system=4.166666665656117,usage_user=4.166666666414029 1568760922000000000
cpu,cpu=cpu2,host=loaner time_active=201382.78999999998,time_guest=30325.8,time_guest_nice=0,time_idle=264686.63,time_iowait=202.77,time_irq=0,time_nice=162.81,time_softirq=3378.34,time_steal=0,time_system=39270.59,time_user=158368.28 1568760922000000000
cpu,cpu=cpu2,host=loaner usage_active=15.999999993480742,usage_guest=1.9999999999126885,usage_guest_nice=0,usage_idle=84.00000000651926,usage_iowait=0,usage_irq=0,usage_nice=0,usage_softirq=2.0000000002764864,usage_steal=0,usage_system=3.999999999825377,usage_user=7.999999998923158 1568760922000000000
cpu,cpu=cpu3,host=loaner time_active=198953.51000000007,time_guest=30344.43,time_guest_nice=0,time_idle=265504.09,time_iowait=187.64,time_irq=0,time_nice=197.47,time_softirq=2301.47,time_steal=0,time_system=39313.73,time_user=156953.2 1568760922000000000
cpu,cpu=cpu3,host=loaner usage_active=10.41666667424579,usage_guest=0,usage_guest_nice=0,usage_idle=89.58333332575421,usage_iowait=0,usage_irq=0,usage_nice=0,usage_softirq=0,usage_steal=0,usage_system=4.166666666666667,usage_user=6.249999998484175 1568760922000000000
cpu,cpu=cpu-total,host=loaner time_active=804450.5299999998,time_guest=121429,time_guest_nice=0,time_idle=2321866.96,time_iowait=1952.86,time_irq=0,time_nice=711.32,time_softirq=16499.1,time_steal=0,time_system=158162.17,time_user=627125.08 1568760922000000000
cpu,cpu=cpu-total,host=loaner usage_active=17.616580305880305,usage_guest=1.036269430422946,usage_guest_nice=0,usage_idle=82.3834196941197,usage_iowait=0,usage_irq=0,usage_nice=0,usage_softirq=1.0362694300459534,usage_steal=0,usage_system=4.145077721691784,usage_user=11.398963731636465 1568760922000000000
```
  ## If true, collect raw CPU time metrics.
  collect_cpu_time = false
  ## If true, compute and report the sum of all non-idle CPU states.
  report_active = false
```

#### Description

The CPU plugin collects standard CPU metrics as defined in `man proc`. All
architectures do not support all of these metrics.

```
cpu  3357 0 4313 1362393
    The amount of time, measured in units of USER_HZ (1/100ths of a second on
    most architectures, use sysconf(_SC_CLK_TCK) to obtain the right value),
    that the system spent in various states:

    user   (1) Time spent in user mode.

    nice   (2) Time spent in user mode with low priority (nice).

    system (3) Time spent in system mode.

    idle   (4) Time spent in the idle task.  This value should be USER_HZ times
    the second entry in the /proc/uptime pseudo-file.

    iowait (since Linux 2.5.41)
           (5) Time waiting for I/O to complete.

    irq (since Linux 2.6.0-test4)
           (6) Time servicing interrupts.

    softirq (since Linux 2.6.0-test4)
           (7) Time servicing softirqs.

    steal (since Linux 2.6.11)
           (8) Stolen time, which is the time spent in other operating systems
           when running in a virtualized environment

    guest (since Linux 2.6.24)
           (9) Time spent running a virtual CPU for guest operating systems
           under the control of the Linux kernel.

    guest_nice (since Linux 2.6.33)
           (10) Time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel).
```

# Measurements:
### CPU Time measurements:

Meta:
- units: CPU Time
- tags: `cpu=<cpuN> or <cpu-total>`

Measurement names:
- cpu_time_user
- cpu_time_system
- cpu_time_idle
- cpu_time_active (must be explicitly enabled by setting `report_active = true`)
- cpu_time_nice
- cpu_time_iowait
- cpu_time_irq
- cpu_time_softirq
- cpu_time_steal
- cpu_time_guest
- cpu_time_guest_nice

### CPU Usage Percent Measurements:

Meta:
- units: percent (out of 100)
- tags: `cpu=<cpuN> or <cpu-total>`

Measurement names:
- cpu_usage_user
- cpu_usage_system
- cpu_usage_idle
- cpu_usage_active (must be explicitly enabled by setting `report_active = true`)
- cpu_usage_nice
- cpu_usage_iowait
- cpu_usage_irq
- cpu_usage_softirq
- cpu_usage_steal
- cpu_usage_guest
- cpu_usage_guest_nice
