# Nginx Plus API Input Plugin

Nginx Plus is a commercial version of the open source web server Nginx. The use
this plugin you will need a license. For more information about the differences
between Nginx (F/OSS) and Nginx Plus, see the Nginx [documentation][diff-doc].

[diff-doc]: https://www.nginx.com/blog/whats-difference-nginx-foss-nginx-plus/

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Read Nginx Plus API advanced status information
[[inputs.nginx_plus_api]]
  ## An array of Nginx API URIs to gather stats.
  urls = ["http://localhost/api"]
  # Nginx API version, default: 3
  # api_version = 3

  # HTTP response timeout (default: 5s)
  response_timeout = "5s"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
```
# Read Nginx Plus API advanced status information
[[inputs.nginx_plus_api]]
  ## An array of Nginx API URIs to gather stats.
  urls = ["http://localhost/api"]
  # Nginx API version, default: 3
  # api_version = 3

  # HTTP response timeout (default: 5s)
  response_timeout = "5s"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
```

## Migration from Nginx Plus (Status) input plugin
```

### Migration from Nginx Plus (Status) input plugin

| Nginx Plus                      | Nginx Plus API                       |
|---------------------------------|--------------------------------------|
| nginx_plus_processes            | nginx_plus_api_processes             |
| nginx_plus_connections          | nginx_plus_api_connections           |
| nginx_plus_ssl                  | nginx_plus_api_ssl                   |
| nginx_plus_requests             | nginx_plus_api_http_requests         |
| nginx_plus_zone                 | nginx_plus_api_http_server_zones     |
| nginx_plus_upstream             | nginx_plus_api_http_upstreams        |
| nginx_plus_upstream_peer        | nginx_plus_api_http_upstream_peers   |
| nginx_plus_cache                | nginx_plus_api_http_caches           |
| nginx_plus_stream_upstream      | nginx_plus_api_stream_upstreams      |
| nginx_plus_stream_upstream_peer | nginx_plus_api_stream_upstream_peers |
| nginx.stream.zone               | nginx_plus_api_stream_server_zones   |

## Measurements by API version

| Measurement                          | API version (api_version) |
|--------------------------------------|---------------------------|
| nginx_plus_api_processes             | >= 3                      |
| nginx_plus_api_connections           | >= 3                      |
| nginx_plus_api_ssl                   | >= 3                      |
| nginx_plus_api_slabs_pages           | >= 3                      |
| nginx_plus_api_slabs_slots           | >= 3                      |
| nginx_plus_api_http_requests         | >= 3                      |
| nginx_plus_api_http_server_zones     | >= 3                      |
| nginx_plus_api_http_upstreams        | >= 3                      |
| nginx_plus_api_http_upstream_peers   | >= 3                      |
| nginx_plus_api_http_caches           | >= 3                      |
| nginx_plus_api_stream_upstreams      | >= 3                      |
| nginx_plus_api_stream_upstream_peers | >= 3                      |
| nginx_plus_api_stream_server_zones   | >= 3                      |
| nginx_plus_api_http_location_zones   | >= 5                      |
| nginx_plus_api_resolver_zones        | >= 5                      |
| nginx_plus_api_http_limit_reqs       | >= 6                      |

## Metrics
### Measurements & Fields:

- nginx_plus_api_processes
  - respawned
- nginx_plus_api_connections
  - accepted
  - dropped
  - active
  - idle
- nginx_plus_api_slabs_pages
  - used
  - free
- nginx_plus_api_slabs_slots
  - used
  - free
  - reqs
  - fails
- nginx_plus_api_ssl
  - handshakes
  - handshakes_failed
  - session_reuses
- nginx_plus_api_http_requests
  - total
  - current
- nginx_plus_api_http_server_zones
  - processing
  - requests
  - responses_1xx
  - responses_2xx
  - responses_3xx
  - responses_4xx
  - responses_5xx
  - responses_total
  - received
  - sent
  - discarded
- nginx_plus_api_http_upstreams
  - keepalive
  - zombies
- nginx_plus_api_http_upstream_peers
  - requests
  - unavail
  - healthchecks_checks
  - header_time
  - state
  - response_time
  - active
  - healthchecks_last_passed
  - weight
  - responses_1xx
  - responses_2xx
  - responses_3xx
  - responses_4xx
  - responses_5xx
  - received
  - healthchecks_fails
  - healthchecks_unhealthy
  - backup
  - responses_total
  - sent
  - fails
  - downtime
- nginx_plus_api_http_caches
  - size
  - max_size
  - cold
  - hit_responses
  - hit_bytes
  - stale_responses
  - stale_bytes
  - updating_responses
  - updating_bytes
  - revalidated_responses
  - revalidated_bytes
  - miss_responses
  - miss_bytes
  - miss_responses_written
  - miss_bytes_written
  - expired_responses
  - expired_bytes
  - expired_responses_written
  - expired_bytes_written
  - bypass_responses
  - bypass_bytes
  - bypass_responses_written
  - bypass_bytes_written
- nginx_plus_api_stream_upstreams
  - zombies
- nginx_plus_api_stream_upstream_peers
  - unavail
  - healthchecks_checks
  - healthchecks_fails
  - healthchecks_unhealthy
  - healthchecks_last_passed
  - response_time
  - state
  - active
  - weight
  - received
  - backup
  - sent
  - fails
  - downtime
- nginx_plus_api_stream_server_zones
  - processing
  - connections
  - received
  - sent
- nginx_plus_api_location_zones
  - requests
  - responses_1xx
  - responses_2xx
  - responses_3xx
  - responses_4xx
  - responses_5xx
  - responses_total
  - received
  - sent
  - discarded
- nginx_plus_api_resolver_zones
  - name
  - srv
  - addr
  - noerror
  - formerr
  - servfail
  - nxdomain
  - notimp
  - refused
  - timedout
  - unknown
- nginx_plus_api_http_limit_reqs
  - passed
  - delayed
  - rejected
  - delayed_dry_run
  - rejected_dry_run

### Tags


### Tags:

- nginx_plus_api_processes, nginx_plus_api_connections, nginx_plus_api_ssl, nginx_plus_api_http_requests
  - source
  - port

- nginx_plus_api_http_upstreams, nginx_plus_api_stream_upstreams
  - upstream
  - source
  - port

- nginx_plus_api_http_server_zones, nginx_plus_api_upstream_server_zones, nginx_plus_api_http_location_zones, nginx_plus_api_resolver_zones, nginx_plus_api_slabs_pages
- nginx_plus_api_http_server_zones, nginx_plus_api_upstream_server_zones
  - source
  - port
  - zone

- nginx_plus_api_slabs_slots
  - source
  - port
  - zone
  - slot

- nginx_plus_api_upstream_peers, nginx_plus_api_stream_upstream_peers
  - id
  - upstream
  - source
  - port
  - upstream_address

- nginx_plus_api_http_caches
  - source
  - port

- nginx_plus_api_http_limit_reqs
  - source
  - port
  - limit

## Example Output

Using this configuration:

```toml
### Example Output:

Using this configuration:
```
[[inputs.nginx_plus_api]]
  ## An array of Nginx Plus API URIs to gather stats.
  urls = ["http://localhost/api"]
```

When run with:

```sh
```
./telegraf -config telegraf.conf -input-filter nginx_plus_api -test
```

It produces:

```text
> nginx_plus_api_processes,port=80,source=demo.nginx.com respawned=0i 1570696321000000000
> nginx_plus_api_connections,port=80,source=demo.nginx.com accepted=68998606i,active=7i,dropped=0i,idle=57i 1570696322000000000
> nginx_plus_api_slabs_pages,port=80,source=demo.nginx.com,zone=hg.nginx.org used=1i,free=503i 1570696322000000000
> nginx_plus_api_slabs_pages,port=80,source=demo.nginx.com,zone=trac.nginx.org used=3i,free=500i 1570696322000000000
> nginx_plus_api_slabs_slots,port=80,source=demo.nginx.com,zone=hg.nginx.org,slot=8 used=1i,free=503i,reqs=10i,fails=0i 1570696322000000000
> nginx_plus_api_slabs_slots,port=80,source=demo.nginx.com,zone=hg.nginx.org,slot=16 used=3i,free=500i,reqs=1024i,fails=0i 1570696322000000000
> nginx_plus_api_slabs_slots,port=80,source=demo.nginx.com,zone=trac.nginx.org,slot=8 used=1i,free=503i,reqs=10i,fails=0i 1570696322000000000
> nginx_plus_api_slabs_slots,port=80,source=demo.nginx.com,zone=trac.nginx.org,slot=16 used=0i,free=1520i,reqs=0i,fails=1i 1570696322000000000
> nginx_plus_api_ssl,port=80,source=demo.nginx.com handshakes=9398978i,handshakes_failed=289353i,session_reuses=1004389i 1570696322000000000
> nginx_plus_api_http_requests,port=80,source=demo.nginx.com current=51i,total=264649353i 1570696322000000000
> nginx_plus_api_http_server_zones,port=80,source=demo.nginx.com,zone=hg.nginx.org discarded=5i,processing=0i,received=24123604i,requests=60138i,responses_1xx=0i,responses_2xx=59353i,responses_3xx=531i,responses_4xx=249i,responses_5xx=0i,responses_total=60133i,sent=830165221i 1570696322000000000
> nginx_plus_api_http_server_zones,port=80,source=demo.nginx.com,zone=trac.nginx.org discarded=250i,processing=0i,received=2184618i,requests=12404i,responses_1xx=0i,responses_2xx=8579i,responses_3xx=2513i,responses_4xx=583i,responses_5xx=479i,responses_total=12154i,sent=139384159i 1570696322000000000
> nginx_plus_api_http_server_zones,port=80,source=demo.nginx.com,zone=lxr.nginx.org discarded=1i,processing=0i,received=1011701i,requests=4523i,responses_1xx=0i,responses_2xx=4332i,responses_3xx=28i,responses_4xx=39i,responses_5xx=123i,responses_total=4522i,sent=72631354i 1570696322000000000
> nginx_plus_api_http_upstreams,port=80,source=demo.nginx.com,upstream=trac-backend keepalive=0i,zombies=0i 1570696322000000000
> nginx_plus_api_http_upstream_peers,id=0,port=80,source=demo.nginx.com,upstream=trac-backend,upstream_address=10.0.0.1:8080 active=0i,backup=false,downtime=0i,fails=0i,header_time=235i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=88581178i,requests=3180i,response_time=235i,responses_1xx=0i,responses_2xx=3168i,responses_3xx=5i,responses_4xx=6i,responses_5xx=0i,responses_total=3179i,sent=1321720i,state="up",unavail=0i,weight=1i 1570696322000000000
> nginx_plus_api_http_upstream_peers,id=1,port=80,source=demo.nginx.com,upstream=trac-backend,upstream_address=10.0.0.1:8081 active=0i,backup=true,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696322000000000
> nginx_plus_api_http_upstreams,port=80,source=demo.nginx.com,upstream=hg-backend keepalive=0i,zombies=0i 1570696322000000000
> nginx_plus_api_http_upstream_peers,id=0,port=80,source=demo.nginx.com,upstream=hg-backend,upstream_address=10.0.0.1:8088 active=0i,backup=false,downtime=0i,fails=0i,header_time=22i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=909402572i,requests=18514i,response_time=88i,responses_1xx=0i,responses_2xx=17799i,responses_3xx=531i,responses_4xx=179i,responses_5xx=0i,responses_total=18509i,sent=10608107i,state="up",unavail=0i,weight=5i 1570696322000000000
> nginx_plus_api_http_upstream_peers,id=1,port=80,source=demo.nginx.com,upstream=hg-backend,upstream_address=10.0.0.1:8089 active=0i,backup=true,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696322000000000
> nginx_plus_api_http_upstreams,port=80,source=demo.nginx.com,upstream=lxr-backend keepalive=0i,zombies=0i 1570696322000000000
> nginx_plus_api_http_upstream_peers,id=0,port=80,source=demo.nginx.com,upstream=lxr-backend,upstream_address=unix:/tmp/cgi.sock active=0i,backup=false,downtime=0i,fails=123i,header_time=91i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=71782888i,requests=4354i,response_time=91i,responses_1xx=0i,responses_2xx=4230i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=4230i,sent=3088656i,state="up",unavail=0i,weight=1i 1570696322000000000
> nginx_plus_api_http_upstream_peers,id=1,port=80,source=demo.nginx.com,upstream=lxr-backend,upstream_address=unix:/tmp/cgib.sock active=0i,backup=true,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,max_conns=42i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696322000000000
> nginx_plus_api_http_upstreams,port=80,source=demo.nginx.com,upstream=demo-backend keepalive=0i,zombies=0i 1570696322000000000
> nginx_plus_api_http_upstream_peers,id=0,port=80,source=demo.nginx.com,upstream=demo-backend,upstream_address=10.0.0.2:15431 active=0i,backup=false,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696322000000000
> nginx_plus_api_http_caches,cache=http_cache,port=80,source=demo.nginx.com bypass_bytes=0i,bypass_bytes_written=0i,bypass_responses=0i,bypass_responses_written=0i,cold=false,expired_bytes=381518640i,expired_bytes_written=363449785i,expired_responses=42114i,expired_responses_written=39954i,hit_bytes=6321885979i,hit_responses=596730i,max_size=536870912i,miss_bytes=48512185i,miss_bytes_written=155600i,miss_responses=6052i,miss_responses_written=136i,revalidated_bytes=0i,revalidated_responses=0i,size=765952i,stale_bytes=0i,stale_responses=0i,updating_bytes=0i,updating_responses=0i 1570696323000000000
> nginx_plus_api_stream_server_zones,port=80,source=demo.nginx.com,zone=postgresql_loadbalancer connections=0i,processing=0i,received=0i,sent=0i 1570696323000000000
> nginx_plus_api_stream_server_zones,port=80,source=demo.nginx.com,zone=dns_loadbalancer connections=0i,processing=0i,received=0i,sent=0i 1570696323000000000
> nginx_plus_api_stream_upstreams,port=80,source=demo.nginx.com,upstream=postgresql_backends zombies=0i 1570696323000000000
> nginx_plus_api_stream_upstream_peers,id=0,port=80,source=demo.nginx.com,upstream=postgresql_backends,upstream_address=10.0.0.2:15432 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696323000000000
> nginx_plus_api_stream_upstream_peers,id=1,port=80,source=demo.nginx.com,upstream=postgresql_backends,upstream_address=10.0.0.2:15433 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696323000000000
> nginx_plus_api_stream_upstream_peers,id=2,port=80,source=demo.nginx.com,upstream=postgresql_backends,upstream_address=10.0.0.2:15434 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696323000000000
> nginx_plus_api_stream_upstream_peers,id=3,port=80,source=demo.nginx.com,upstream=postgresql_backends,upstream_address=10.0.0.2:15435 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1570696323000000000
> nginx_plus_api_stream_upstreams,port=80,source=demo.nginx.com,upstream=dns_udp_backends zombies=0i 1570696323000000000
> nginx_plus_api_stream_upstream_peers,id=0,port=80,source=demo.nginx.com,upstream=dns_udp_backends,upstream_address=10.0.0.5:53 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="up",unavail=0i,weight=2i 1570696323000000000
> nginx_plus_api_stream_upstream_peers,id=1,port=80,source=demo.nginx.com,upstream=dns_udp_backends,upstream_address=10.0.0.2:53 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="up",unavail=0i,weight=1i 1570696323000000000
> nginx_plus_api_stream_upstream_peers,id=2,port=80,source=demo.nginx.com,upstream=dns_udp_backends,upstream_address=10.0.0.7:53 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1570696323000000000
> nginx_plus_api_stream_upstreams,port=80,source=demo.nginx.com,upstream=unused_tcp_backends zombies=0i 1570696323000000000
> nginx_plus_api_http_location_zones,port=80,source=demo.nginx.com,zone=swagger discarded=0i,received=1622i,requests=8i,responses_1xx=0i,responses_2xx=7i,responses_3xx=0i,responses_4xx=1i,responses_5xx=0i,responses_total=8i,sent=638333i 1570696323000000000
> nginx_plus_api_http_location_zones,port=80,source=demo.nginx.com,zone=api-calls discarded=64i,received=337530181i,requests=1726513i,responses_1xx=0i,responses_2xx=1726428i,responses_3xx=0i,responses_4xx=21i,responses_5xx=0i,responses_total=1726449i,sent=1902577668i 1570696323000000000
> nginx_plus_api_resolver_zones,port=80,source=demo.nginx.com,zone=resolver1 addr=0i,formerr=0i,name=0i,noerror=0i,notimp=0i,nxdomain=0i,refused=0i,servfail=0i,srv=0i,timedout=0i,unknown=0i 1570696324000000000
> nginx_plus_api_http_limit_reqs,port=80,source=demo.nginx.com,limit=limit_1 delayed=0i,delayed_dry_run=0i,passed=6i,rejected=9i,rejected_dry_run=0i 1570696322000000000
> nginx_plus_api_http_limit_reqs,port=80,source=demo.nginx.com,limit=limit_2 delayed=13i,delayed_dry_run=3i,passed=6i,rejected=1i,rejected_dry_run=31i 1570696322000000000
```
> nginx_plus_api_processes,host=localhost,port=80,source=localhost respawned=0i 1539163505000000000
> nginx_plus_api_connections,host=localhost,port=80,source=localhost accepted=120890747i,active=6i,dropped=0i,idle=67i 1539163505000000000
> nginx_plus_api_ssl,host=localhost,port=80,source=localhost handshakes=2983938i,handshakes_failed=54350i,session_reuses=2485267i 1539163506000000000
> nginx_plus_api_http_requests,host=localhost,port=80,source=localhost current=12i,total=175270198i 1539163506000000000
> nginx_plus_api_http_server_zones,host=localhost,port=80,source=localhost,zone=hg.nginx.org discarded=45i,processing=0i,received=35723884i,requests=134102i,responses_1xx=0i,responses_2xx=96890i,responses_3xx=6892i,responses_4xx=30270i,responses_5xx=5i,responses_total=134057i,sent=3681826618i 1539163506000000000
> nginx_plus_api_http_server_zones,host=localhost,port=80,source=localhost,zone=trac.nginx.org discarded=4034i,processing=9i,received=282399663i,requests=336129i,responses_1xx=0i,responses_2xx=101264i,responses_3xx=25454i,responses_4xx=68961i,responses_5xx=136407i,responses_total=332086i,sent=2346677493i 1539163506000000000
> nginx_plus_api_http_server_zones,host=localhost,port=80,source=localhost,zone=lxr.nginx.org discarded=4i,processing=1i,received=7223569i,requests=29661i,responses_1xx=0i,responses_2xx=28584i,responses_3xx=73i,responses_4xx=390i,responses_5xx=609i,responses_total=29656i,sent=5811238975i 1539163506000000000
> nginx_plus_api_http_upstreams,host=localhost,port=80,source=localhost,upstream=trac-backend keepalive=0i,zombies=0i 1539163506000000000
> nginx_plus_api_http_upstream_peers,host=localhost,id=0,port=80,source=localhost,upstream=trac-backend,upstream_address=10.0.0.1:8080 active=0i,backup=false,downtime=53870i,fails=5i,header_time=421i,healthchecks_checks=17275i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=1885213684i,requests=88476i,response_time=423i,responses_1xx=0i,responses_2xx=50997i,responses_3xx=205i,responses_4xx=34344i,responses_5xx=2076i,responses_total=87622i,sent=189938404i,state="up",unavail=5i,weight=1i 1539163506000000000
> nginx_plus_api_http_upstream_peers,host=localhost,id=1,port=80,source=localhost,upstream=trac-backend,upstream_address=10.0.0.1:8081 active=0i,backup=true,downtime=173957231i,fails=0i,healthchecks_checks=17394i,healthchecks_fails=17394i,healthchecks_last_passed=false,healthchecks_unhealthy=1i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="unhealthy",unavail=0i,weight=1i 1539163506000000000
> nginx_plus_api_http_upstreams,host=localhost,port=80,source=localhost,upstream=hg-backend keepalive=0i,zombies=0i 1539163506000000000
> nginx_plus_api_http_upstream_peers,host=localhost,id=0,port=80,source=localhost,upstream=hg-backend,upstream_address=10.0.0.1:8088 active=0i,backup=false,downtime=0i,fails=0i,header_time=22i,healthchecks_checks=17319i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=3724240605i,requests=89563i,response_time=44i,responses_1xx=0i,responses_2xx=81996i,responses_3xx=6886i,responses_4xx=639i,responses_5xx=5i,responses_total=89526i,sent=31597952i,state="up",unavail=0i,weight=5i 1539163506000000000
> nginx_plus_api_http_upstream_peers,host=localhost,id=1,port=80,source=localhost,upstream=hg-backend,upstream_address=10.0.0.1:8089 active=0i,backup=true,downtime=173957231i,fails=0i,healthchecks_checks=17394i,healthchecks_fails=17394i,healthchecks_last_passed=false,healthchecks_unhealthy=1i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="unhealthy",unavail=0i,weight=1i 1539163506000000000
> nginx_plus_api_http_upstreams,host=localhost,port=80,source=localhost,upstream=lxr-backend keepalive=0i,zombies=0i 1539163506000000000
> nginx_plus_api_http_upstream_peers,host=localhost,id=0,port=80,source=localhost,upstream=lxr-backend,upstream_address=unix:/tmp/cgi.sock active=0i,backup=false,downtime=0i,fails=609i,header_time=111i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=6220215064i,requests=28278i,response_time=172i,responses_1xx=0i,responses_2xx=27665i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=27665i,sent=21337016i,state="up",unavail=0i,weight=1i 1539163506000000000
> nginx_plus_api_http_upstream_peers,host=localhost,id=1,port=80,source=localhost,upstream=lxr-backend,upstream_address=unix:/tmp/cgib.sock active=0i,backup=true,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,max_conns=42i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="up",unavail=0i,weight=1i 1539163506000000000
> nginx_plus_api_http_upstreams,host=localhost,port=80,source=localhost,upstream=demo-backend keepalive=0i,zombies=0i 1539163506000000000
> nginx_plus_api_http_upstream_peers,host=localhost,id=0,port=80,source=localhost,upstream=demo-backend,upstream_address=10.0.0.2:15431 active=0i,backup=false,downtime=0i,fails=0i,healthchecks_checks=173640i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=0i,requests=0i,responses_1xx=0i,responses_2xx=0i,responses_3xx=0i,responses_4xx=0i,responses_5xx=0i,responses_total=0i,sent=0i,state="up",unavail=0i,weight=1i 1539163506000000000
> nginx_plus_api_http_caches,cache=http_cache,host=localhost,port=80,source=localhost bypass_bytes=0i,bypass_bytes_written=0i,bypass_responses=0i,bypass_responses_written=0i,cold=false,expired_bytes=133671410i,expired_bytes_written=129210272i,expired_responses=15721i,expired_responses_written=15213i,hit_bytes=2459840828i,hit_responses=231195i,max_size=536870912i,miss_bytes=18742246i,miss_bytes_written=85199i,miss_responses=2816i,miss_responses_written=69i,revalidated_bytes=0i,revalidated_responses=0i,size=774144i,stale_bytes=0i,stale_responses=0i,updating_bytes=0i,updating_responses=0i 1539163506000000000
> nginx_plus_api_stream_server_zones,host=localhost,port=80,source=localhost,zone=postgresql_loadbalancer connections=173639i,processing=0i,received=17884817i,sent=33685966i 1539163506000000000
> nginx_plus_api_stream_server_zones,host=localhost,port=80,source=localhost,zone=dns_loadbalancer connections=97255i,processing=0i,received=2699082i,sent=16566552i 1539163506000000000
> nginx_plus_api_stream_upstreams,host=localhost,port=80,source=localhost,upstream=postgresql_backends zombies=0i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=0,port=80,source=localhost,upstream=postgresql_backends,upstream_address=10.0.0.2:15432 active=0i,backup=false,connect_time=4i,connections=57880i,downtime=0i,fails=0i,first_byte_time=10i,healthchecks_checks=34781i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=11228720i,response_time=10i,sent=5961640i,state="up",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=1,port=80,source=localhost,upstream=postgresql_backends,upstream_address=10.0.0.2:15433 active=0i,backup=false,connect_time=3i,connections=57880i,downtime=0i,fails=0i,first_byte_time=9i,healthchecks_checks=34781i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=11228720i,response_time=10i,sent=5961640i,state="up",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=2,port=80,source=localhost,upstream=postgresql_backends,upstream_address=10.0.0.2:15434 active=0i,backup=false,connect_time=2i,connections=57879i,downtime=0i,fails=0i,first_byte_time=9i,healthchecks_checks=34781i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=11228526i,response_time=9i,sent=5961537i,state="up",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=3,port=80,source=localhost,upstream=postgresql_backends,upstream_address=10.0.0.2:15435 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstreams,host=localhost,port=80,source=localhost,upstream=dns_udp_backends zombies=0i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=0,port=80,source=localhost,upstream=dns_udp_backends,upstream_address=10.0.0.5:53 active=0i,backup=false,connect_time=0i,connections=64837i,downtime=0i,fails=0i,first_byte_time=17i,healthchecks_checks=34761i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=10996616i,response_time=17i,sent=1791693i,state="up",unavail=0i,weight=2i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=1,port=80,source=localhost,upstream=dns_udp_backends,upstream_address=10.0.0.2:53 active=0i,backup=false,connect_time=0i,connections=32418i,downtime=0i,fails=0i,first_byte_time=17i,healthchecks_checks=34761i,healthchecks_fails=0i,healthchecks_last_passed=true,healthchecks_unhealthy=0i,received=5569936i,response_time=17i,sent=907389i,state="up",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=2,port=80,source=localhost,upstream=dns_udp_backends,upstream_address=10.0.0.7:53 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstreams,host=localhost,port=80,source=localhost,upstream=unused_tcp_backends zombies=0i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=1,port=80,source=localhost,upstream=unused_tcp_backends,upstream_address=95.211.80.227:80 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=2,port=80,source=localhost,upstream=unused_tcp_backends,upstream_address=206.251.255.63:80 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=3,port=80,source=localhost,upstream=unused_tcp_backends,upstream_address=[2001:1af8:4060:a004:21::e3]:80 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1539163507000000000
> nginx_plus_api_stream_upstream_peers,host=localhost,id=4,port=80,source=localhost,upstream=unused_tcp_backends,upstream_address=[2606:7100:1:69::3f]:80 active=0i,backup=false,connections=0i,downtime=0i,fails=0i,healthchecks_checks=0i,healthchecks_fails=0i,healthchecks_unhealthy=0i,received=0i,sent=0i,state="down",unavail=0i,weight=1i 1539163507000000000
```

### Reference material

- [api documentation](http://demo.nginx.com/swagger-ui/#/)
- [nginx_api_module documentation](http://nginx.org/en/docs/http/ngx_http_api_module.html)
[api documentation](http://demo.nginx.com/swagger-ui/#/)
