# Ceph Storage Input Plugin

Collects performance metrics from the MON and OSD nodes in a Ceph storage
cluster.

Ceph has introduced a Telegraf and Influx plugin in the 13.x Mimic release.
The Telegraf module sends to a Telegraf configured with a socket_listener.
[Learn more in their docs](https://docs.ceph.com/en/latest/mgr/telegraf/)

## Admin Socket Stats

This gatherer works by scanning the configured SocketDir for OSD, MON, MDS
and RGW socket files.  When it finds a MON socket, it runs

```shell
ceph --admin-daemon $file perfcounters_dump
```

For OSDs it runs

```shell
ceph --admin-daemon $file perf dump
```

The resulting JSON is parsed and grouped into collections, based on
top-level key. Top-level keys are used as collection tags, and all
sub-keys are flattened. For example:

```json
Collects performance metrics from the MON and OSD nodes in a Ceph storage cluster.

*Admin Socket Stats*

This gatherer works by scanning the configured SocketDir for OSD and MON socket files.  When it finds
a MON socket, it runs **ceph --admin-daemon $file perfcounters_dump**. For OSDs it runs **ceph --admin-daemon $file perf dump**

The resulting JSON is parsed and grouped into collections, based on top-level key.  Top-level keys are
used as collection tags, and all sub-keys are flattened. For example:

```
 {
   "paxos": {
     "refresh": 9363435,
     "refresh_latency": {
       "avgcount": 9363435,
       "sum": 5378.794002000
     }
   }
 }
```

Would be parsed into the following metrics, all of which would be tagged
with `collection=paxos`:

- refresh = 9363435
- refresh_latency.avgcount: 9363435
- refresh_latency.sum: 5378.794002000

## Cluster Stats

This gatherer works by invoking ceph commands against the cluster thus only
requires the ceph client, valid ceph configuration and an access key to
function (the ceph_config and ceph_user configuration variables work in
conjunction to specify these prerequisites). It may be run on any server you
wish which has access to the cluster.  The currently supported commands are:

- ceph status
- ceph df
- ceph osd pool stats

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Collects performance metrics from the MON, OSD, MDS and RGW nodes
# in a Ceph storage cluster.
[[inputs.ceph]]
  ## This is the recommended interval to poll. Too frequent and you
  ## will lose data points due to timeouts during rebalancing and recovery
Would be parsed into the following metrics, all of which would be tagged with collection=paxos:

 - refresh = 9363435
 - refresh_latency.avgcount: 9363435
 - refresh_latency.sum: 5378.794002000


*Cluster Stats*

This gatherer works by invoking ceph commands against the cluster thus only requires the ceph client, valid
ceph configuration and an access key to function (the ceph_config and ceph_user configuration variables work
in conjunction to specify these prerequisites). It may be run on any server you wish which has access to
the cluster.  The currently supported commands are:

* ceph status
* ceph df
* ceph osd pool stats

### Configuration:

```
# Collects performance metrics from the MON and OSD nodes in a Ceph storage cluster.
[[inputs.ceph]]
  ## This is the recommended interval to poll.  Too frequent and you will lose
  ## data points due to timeouts during rebalancing and recovery
  interval = '1m'

  ## All configuration values are optional, defaults are shown below

  ## location of ceph binary
  ceph_binary = "/usr/bin/ceph"

  ## directory in which to look for socket files
  socket_dir = "/var/run/ceph"

  ## prefix of MON and OSD socket files, used to determine socket type
  mon_prefix = "ceph-mon"
  osd_prefix = "ceph-osd"
  mds_prefix = "ceph-mds"
  rgw_prefix = "ceph-client"

  ## suffix used to identify socket files
  socket_suffix = "asok"

  ## Ceph user to authenticate as, ceph will search for the corresponding
  ## keyring e.g. client.admin.keyring in /etc/ceph, or the explicit path
  ## defined in the client section of ceph.conf for example:
  ## Ceph user to authenticate as, ceph will search for the corresponding keyring
  ## e.g. client.admin.keyring in /etc/ceph, or the explicit path defined in the
  ## client section of ceph.conf for example:
  ##
  ##     [client.telegraf]
  ##         keyring = /etc/ceph/client.telegraf.keyring
  ##
  ## Consult the ceph documentation for more detail on keyring generation.
  ceph_user = "client.admin"

  ## Ceph configuration to use to locate the cluster
  ceph_config = "/etc/ceph/ceph.conf"

  ## Whether to gather statistics via the admin socket
  gather_admin_socket_stats = true

  ## Whether to gather statistics via ceph commands, requires ceph_user
  ## and ceph_config to be specified
  gather_cluster_stats = false
```

## Metrics

### Admin Socket

All fields are collected under the **ceph** measurement and stored as
float64s. For a full list of fields, see the sample perf dumps in ceph_test.go.

All admin measurements will have the following tags:

- type: either 'osd', 'mon', 'mds' or 'rgw' to indicate the queried node type
- id: a unique string identifier, parsed from the socket file name for the node
- collection: the top-level key under which these fields were reported.
  Possible values are:
  ## Whether to gather statistics via ceph commands, requires ceph_user and ceph_config
  ## to be specified
  gather_cluster_stats = false
```

### Measurements & Fields:

*Admin Socket Stats*

All fields are collected under the **ceph** measurement and stored as float64s. For a full list of fields, see the sample perf dumps in ceph_test.go.

*Cluster Stats*

* ceph\_osdmap
  * epoch (float)
  * full (boolean)
  * nearfull (boolean)
  * num\_in\_osds (float)
  * num\_osds (float)
  * num\_remremapped\_pgs (float)
  * num\_up\_osds (float)

* ceph\_pgmap
  * bytes\_avail (float)
  * bytes\_total (float)
  * bytes\_used (float)
  * data\_bytes (float)
  * num\_pgs (float)
  * op\_per\_sec (float)
  * read\_bytes\_sec (float)
  * version (float)
  * write\_bytes\_sec (float)
  * recovering\_bytes\_per\_sec (float)
  * recovering\_keys\_per\_sec (float)
  * recovering\_objects\_per\_sec (float)

* ceph\_pgmap\_state
  * count (float)

* ceph\_usage
  * bytes\_used (float)
  * kb\_used (float)
  * max\_avail (float)
  * objects (float)

* ceph\_pool\_usage
  * bytes\_used (float)
  * kb\_used (float)
  * max\_avail (float)
  * objects (float)

* ceph\_pool\_stats
  * op\_per\_sec (float)
  * read\_bytes\_sec (float)
  * write\_bytes\_sec (float)
  * recovering\_object\_per\_sec (float)
  * recovering\_bytes\_per\_sec (float)
  * recovering\_keys\_per\_sec (float)

### Tags:

*Admin Socket Stats*

All measurements will have the following tags:

- type: either 'osd' or 'mon' to indicate which type of node was queried
- id: a unique string identifier, parsed from the socket file name for the node
- collection: the top-level key under which these fields were reported. Possible values are:
  - for MON nodes:
    - cluster
    - leveldb
    - mon
    - paxos
    - throttle-mon_client_bytes
    - throttle-mon_daemon_bytes
    - throttle-msgr_dispatch_throttler-mon
  - for OSD nodes:
    - WBThrottle
    - filestore
    - leveldb
    - mutex-FileJournal::completions_lock
    - mutex-FileJournal::finisher_lock
    - mutex-FileJournal::write_lock
    - mutex-FileJournal::writeq_lock
    - mutex-JOS::ApplyManager::apply_lock
    - mutex-JOS::ApplyManager::com_lock
    - mutex-JOS::SubmitManager::lock
    - mutex-WBThrottle::lock
    - objecter
    - osd
    - recoverystate_perf
    - throttle-filestore_bytes
    - throttle-filestore_ops
    - throttle-msgr_dispatch_throttler-client
    - throttle-msgr_dispatch_throttler-cluster
    - throttle-msgr_dispatch_throttler-hb_back_server
    - throttle-msgr_dispatch_throttler-hb_front_serve
    - throttle-msgr_dispatch_throttler-hbclient
    - throttle-msgr_dispatch_throttler-ms_objecter
    - throttle-objecter_bytes
    - throttle-objecter_ops
    - throttle-osd_client_bytes
    - throttle-osd_client_messages
  - for MDS nodes:
    - AsyncMessenger::Worker-0
    - AsyncMessenger::Worker-1
    - AsyncMessenger::Worker-2
    - finisher-PurgeQueue
    - mds
    - mds_cache
    - mds_log
    - mds_mem
    - mds_server
    - mds_sessions
    - objecter
    - purge_queue
    - throttle-msgr_dispatch_throttler-mds
    - throttle-objecter_bytes
    - throttle-objecter_ops
    - throttle-write_buf_throttle
  - for RGW nodes:
    - AsyncMessenger::Worker-0
    - AsyncMessenger::Worker-1
    - AsyncMessenger::Worker-2
    - cct
    - finisher-radosclient
    - mempool
    - objecter
    - rgw
    - simple-throttler
    - throttle-msgr_dispatch_throttler-radosclient
    - throttle-objecter_bytes
    - throttle-objecter_ops
    - throttle-rgw_async_rados_ops

## Cluster

- ceph_fsmap
  - fields:
    - up (float)
    - in (float)
    - max (float)
    - up_standby (float)

- ceph_health
  - fields:
    - status (string)
    - status_code (int)
    - overall_status (string, exists only in ceph <15)

- ceph_monmap
  - fields:
    - num_mons (float)

- ceph_osdmap
  - fields:
    - epoch (float)
    - full (bool, exists only in ceph <15)
    - nearfull (bool, exists only in ceph <15)
    - num_in_osds (float)
    - num_osds (float)
    - num_remapped_pgs (float)
    - num_up_osds (float)

- ceph_pgmap
  - fields:
    - bytes_avail (float)
    - bytes_total (float)
    - bytes_used (float)
    - data_bytes (float)
    - degraded_objects (float)
    - degraded_ratio (float)
    - degraded_total (float)
    - inactive_pgs_ratio (float)
    - num_bytes_recovered (float)
    - num_keys_recovered (float)
    - num_objects (float)
    - num_objects_recovered (float)
    - num_pgs (float)
    - num_pools (float)
    - op_per_sec (float, exists only in ceph <10)
    - read_bytes_sec (float)
    - read_op_per_sec (float)
    - recovering_bytes_per_sec (float)
    - recovering_keys_per_sec (float)
    - recovering_objects_per_sec (float)
    - version (float)
    - write_bytes_sec (float)
    - write_op_per_sec (float)

- ceph_pgmap_state
  - tags:
    - state
  - fields:
    - count (float)

- ceph_usage
  - fields:
    - num_osd (float)
    - num_per_pool_omap_osds (float)
    - num_per_pool_osds (float)
    - total_avail (float, exists only in ceph <0.84)
    - total_avail_bytes (float)
    - total_bytes (float)
    - total_space (float, exists only in ceph <0.84)
    - total_used (float, exists only in ceph <0.84)
    - total_used_bytes (float)
    - total_used_raw_bytes (float)
    - total_used_raw_ratio (float)

- ceph_deviceclass_usage
  - tags:
    - class
  - fields:
    - total_avail_bytes (float)
    - total_bytes (float)
    - total_used_bytes (float)
    - total_used_raw_bytes (float)
    - total_used_raw_ratio (float)

- ceph_pool_usage
  - tags:
    - name
  - fields:
    - bytes_used (float)
    - kb_used (float)
    - max_avail (float)
    - objects (float)
    - percent_used (float)
    - stored (float)

- ceph_pool_stats
  - tags:
    - name
  - fields:
    - degraded_objects (float)
    - degraded_ratio (float)
    - degraded_total (float)
    - num_bytes_recovered (float)
    - num_keys_recovered (float)
    - num_objects_recovered (float)
    - op_per_sec (float, exists only in ceph <10)
    - read_bytes_sec (float)
    - read_op_per_sec (float)
    - recovering_bytes_per_sec (float)
    - recovering_keys_per_sec (float)
    - recovering_objects_per_sec (float)
    - write_bytes_sec (float)
    - write_op_per_sec (float)

## Example Output

Below is an example of a cluster stats:

```text
ceph_fsmap,host=ceph in=1,max=1,up=1,up_standby=2 1646782035000000000
ceph_health,host=ceph status="HEALTH_OK",status_code=2 1646782035000000000
ceph_monmap,host=ceph num_mons=3 1646782035000000000
ceph_osdmap,host=ceph epoch=10560,num_in_osds=6,num_osds=6,num_remapped_pgs=0,num_up_osds=6 1646782035000000000
ceph_pgmap,host=ceph bytes_avail=7863124942848,bytes_total=14882929901568,bytes_used=7019804958720,data_bytes=2411111520818,degraded_objects=0,degraded_ratio=0,degraded_total=0,inactive_pgs_ratio=0,num_bytes_recovered=0,num_keys_recovered=0,num_objects=973030,num_objects_recovered=0,num_pgs=233,num_pools=6,read_bytes_sec=7334,read_op_per_sec=2,recovering_bytes_per_sec=0,recovering_keys_per_sec=0,recovering_objects_per_sec=0,version=0,write_bytes_sec=13113085,write_op_per_sec=355 1646782035000000000
ceph_pgmap_state,host=ceph,state=active+clean count=233 1646782035000000000
ceph_usage,host=ceph num_osds=6,num_per_pool_omap_osds=6,num_per_pool_osds=6,total_avail_bytes=7863124942848,total_bytes=14882929901568,total_used_bytes=7019804958720,total_used_raw_bytes=7019804958720,total_used_raw_ratio=0.47166821360588074 1646782035000000000
ceph_deviceclass_usage,class=hdd,host=ceph total_avail_bytes=6078650843136,total_bytes=12002349023232,total_used_bytes=5923698180096,total_used_raw_bytes=5923698180096,total_used_raw_ratio=0.49354490637779236 1646782035000000000
ceph_deviceclass_usage,class=ssd,host=ceph total_avail_bytes=1784474099712,total_bytes=2880580878336,total_used_bytes=1096106778624,total_used_raw_bytes=1096106778624,total_used_raw_ratio=0.3805158734321594 1646782035000000000
ceph_pool_usage,host=ceph,name=Foo bytes_used=2019483848658,kb_used=1972152196,max_avail=1826022621184,objects=161029,percent_used=0.26935243606567383,stored=672915064134 1646782035000000000
ceph_pool_usage,host=ceph,name=Bar_metadata bytes_used=4370899787,kb_used=4268457,max_avail=546501918720,objects=89702,percent_used=0.002658897778019309,stored=1456936576 1646782035000000000
ceph_pool_usage,host=ceph,name=Bar_data bytes_used=3893328740352,kb_used=3802078848,max_avail=1826022621184,objects=518396,percent_used=0.41544806957244873,stored=1292214337536 1646782035000000000
ceph_pool_usage,host=ceph,name=device_health_metrics bytes_used=85289044,kb_used=83291,max_avail=3396406870016,objects=9,percent_used=0.000012555617104226258,stored=42644520 1646782035000000000
ceph_pool_usage,host=ceph,name=Foo_Fast bytes_used=597511814461,kb_used=583507632,max_avail=546501918720,objects=67014,percent_used=0.2671019732952118,stored=199093853972 1646782035000000000
ceph_pool_usage,host=ceph,name=Bar_data_fast bytes_used=490009280512,kb_used=478524688,max_avail=546501918720,objects=136880,percent_used=0.23010368645191193,stored=163047325696 1646782035000000000
ceph_pool_stats,host=ceph,name=Foo degraded_objects=0,degraded_ratio=0,degraded_total=0,num_bytes_recovered=0,num_keys_recovered=0,num_objects_recovered=0,read_bytes_sec=0,read_op_per_sec=0,recovering_bytes_per_sec=0,recovering_keys_per_sec=0,recovering_objects_per_sec=0,write_bytes_sec=27720,write_op_per_sec=4 1646782036000000000
ceph_pool_stats,host=ceph,name=Bar_metadata degraded_objects=0,degraded_ratio=0,degraded_total=0,num_bytes_recovered=0,num_keys_recovered=0,num_objects_recovered=0,read_bytes_sec=9638,read_op_per_sec=3,recovering_bytes_per_sec=0,recovering_keys_per_sec=0,recovering_objects_per_sec=0,write_bytes_sec=11802778,write_op_per_sec=60 1646782036000000000
ceph_pool_stats,host=ceph,name=Bar_data degraded_objects=0,degraded_ratio=0,degraded_total=0,num_bytes_recovered=0,num_keys_recovered=0,num_objects_recovered=0,read_bytes_sec=0,read_op_per_sec=0,recovering_bytes_per_sec=0,recovering_keys_per_sec=0,recovering_objects_per_sec=0,write_bytes_sec=0,write_op_per_sec=104 1646782036000000000
ceph_pool_stats,host=ceph,name=device_health_metrics degraded_objects=0,degraded_ratio=0,degraded_total=0,num_bytes_recovered=0,num_keys_recovered=0,num_objects_recovered=0,read_bytes_sec=0,read_op_per_sec=0,recovering_bytes_per_sec=0,recovering_keys_per_sec=0,recovering_objects_per_sec=0,write_bytes_sec=0,write_op_per_sec=0 1646782036000000000
ceph_pool_stats,host=ceph,name=Foo_Fast degraded_objects=0,degraded_ratio=0,degraded_total=0,num_bytes_recovered=0,num_keys_recovered=0,num_objects_recovered=0,read_bytes_sec=0,read_op_per_sec=0,recovering_bytes_per_sec=0,recovering_keys_per_sec=0,recovering_objects_per_sec=0,write_bytes_sec=11173,write_op_per_sec=1 1646782036000000000
ceph_pool_stats,host=ceph,name=Bar_data_fast degraded_objects=0,degraded_ratio=0,degraded_total=0,num_bytes_recovered=0,num_keys_recovered=0,num_objects_recovered=0,read_bytes_sec=0,read_op_per_sec=0,recovering_bytes_per_sec=0,recovering_keys_per_sec=0,recovering_objects_per_sec=0,write_bytes_sec=2155404,write_op_per_sec=262 1646782036000000000
```

Below is an example of admin socket stats:

```text
ceph,collection=cct,host=stefanmon1,id=stefanmon1,type=monitor total_workers=0,unhealthy_workers=0 1587117563000000000
ceph,collection=mempool,host=stefanmon1,id=stefanmon1,type=monitor bloom_filter_bytes=0,bloom_filter_items=0,bluefs_bytes=0,bluefs_items=0,bluestore_alloc_bytes=0,bluestore_alloc_items=0,bluestore_cache_data_bytes=0,bluestore_cache_data_items=0,bluestore_cache_onode_bytes=0,bluestore_cache_onode_items=0,bluestore_cache_other_bytes=0,bluestore_cache_other_items=0,bluestore_fsck_bytes=0,bluestore_fsck_items=0,bluestore_txc_bytes=0,bluestore_txc_items=0,bluestore_writing_bytes=0,bluestore_writing_deferred_bytes=0,bluestore_writing_deferred_items=0,bluestore_writing_items=0,buffer_anon_bytes=719152,buffer_anon_items=192,buffer_meta_bytes=352,buffer_meta_items=4,mds_co_bytes=0,mds_co_items=0,osd_bytes=0,osd_items=0,osd_mapbl_bytes=0,osd_mapbl_items=0,osd_pglog_bytes=0,osd_pglog_items=0,osdmap_bytes=15872,osdmap_items=138,osdmap_mapping_bytes=63112,osdmap_mapping_items=7626,pgmap_bytes=38680,pgmap_items=477,unittest_1_bytes=0,unittest_1_items=0,unittest_2_bytes=0,unittest_2_items=0 1587117563000000000
ceph,collection=throttle-mon_client_bytes,host=stefanmon1,id=stefanmon1,type=monitor get=1041157,get_or_fail_fail=0,get_or_fail_success=1041157,get_started=0,get_sum=64928901,max=104857600,put=1041157,put_sum=64928901,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117563000000000
ceph,collection=throttle-msgr_dispatch_throttler-mon,host=stefanmon1,id=stefanmon1,type=monitor get=12695426,get_or_fail_fail=0,get_or_fail_success=12695426,get_started=0,get_sum=42542216884,max=104857600,put=12695426,put_sum=42542216884,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117563000000000
ceph,collection=finisher-mon_finisher,host=stefanmon1,id=stefanmon1,type=monitor complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117563000000000
ceph,collection=finisher-monstore,host=stefanmon1,id=stefanmon1,type=monitor complete_latency.avgcount=1609831,complete_latency.avgtime=0.015857621,complete_latency.sum=25528.09131035,queue_len=0 1587117563000000000
ceph,collection=mon,host=stefanmon1,id=stefanmon1,type=monitor election_call=25,election_lose=0,election_win=22,num_elections=94,num_sessions=3,session_add=174679,session_rm=439316,session_trim=137 1587117563000000000
ceph,collection=throttle-mon_daemon_bytes,host=stefanmon1,id=stefanmon1,type=monitor get=72697,get_or_fail_fail=0,get_or_fail_success=72697,get_started=0,get_sum=32261199,max=419430400,put=72697,put_sum=32261199,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117563000000000
ceph,collection=rocksdb,host=stefanmon1,id=stefanmon1,type=monitor compact=1,compact_queue_len=0,compact_queue_merge=1,compact_range=19126,get=62449211,get_latency.avgcount=62449211,get_latency.avgtime=0.000022216,get_latency.sum=1387.371811726,rocksdb_write_delay_time.avgcount=0,rocksdb_write_delay_time.avgtime=0,rocksdb_write_delay_time.sum=0,rocksdb_write_memtable_time.avgcount=0,rocksdb_write_memtable_time.avgtime=0,rocksdb_write_memtable_time.sum=0,rocksdb_write_pre_and_post_time.avgcount=0,rocksdb_write_pre_and_post_time.avgtime=0,rocksdb_write_pre_and_post_time.sum=0,rocksdb_write_wal_time.avgcount=0,rocksdb_write_wal_time.avgtime=0,rocksdb_write_wal_time.sum=0,submit_latency.avgcount=0,submit_latency.avgtime=0,submit_latency.sum=0,submit_sync_latency.avgcount=3219961,submit_sync_latency.avgtime=0.007532173,submit_sync_latency.sum=24253.303584224,submit_transaction=0,submit_transaction_sync=3219961 1587117563000000000
ceph,collection=AsyncMessenger::Worker-0,host=stefanmon1,id=stefanmon1,type=monitor msgr_active_connections=148317,msgr_created_connections=162806,msgr_recv_bytes=11557888328,msgr_recv_messages=5113369,msgr_running_fast_dispatch_time=0,msgr_running_recv_time=868.377161686,msgr_running_send_time=1626.525392721,msgr_running_total_time=4222.235694322,msgr_send_bytes=91516226816,msgr_send_messages=6973706 1587117563000000000
ceph,collection=AsyncMessenger::Worker-2,host=stefanmon1,id=stefanmon1,type=monitor msgr_active_connections=146396,msgr_created_connections=159788,msgr_recv_bytes=2162802496,msgr_recv_messages=689168,msgr_running_fast_dispatch_time=0,msgr_running_recv_time=164.148550562,msgr_running_send_time=153.462890368,msgr_running_total_time=644.188791379,msgr_send_bytes=7422484152,msgr_send_messages=749381 1587117563000000000
ceph,collection=cluster,host=stefanmon1,id=stefanmon1,type=monitor num_bytes=5055,num_mon=3,num_mon_quorum=3,num_object=245,num_object_degraded=0,num_object_misplaced=0,num_object_unfound=0,num_osd=9,num_osd_in=8,num_osd_up=8,num_pg=504,num_pg_active=504,num_pg_active_clean=504,num_pg_peering=0,num_pool=17,osd_bytes=858959904768,osd_bytes_avail=849889787904,osd_bytes_used=9070116864,osd_epoch=203 1587117563000000000
ceph,collection=paxos,host=stefanmon1,id=stefanmon1,type=monitor accept_timeout=1,begin=1609847,begin_bytes.avgcount=1609847,begin_bytes.sum=41408662074,begin_keys.avgcount=1609847,begin_keys.sum=4829541,begin_latency.avgcount=1609847,begin_latency.avgtime=0.007213392,begin_latency.sum=11612.457661116,collect=0,collect_bytes.avgcount=0,collect_bytes.sum=0,collect_keys.avgcount=0,collect_keys.sum=0,collect_latency.avgcount=0,collect_latency.avgtime=0,collect_latency.sum=0,collect_timeout=1,collect_uncommitted=17,commit=1609831,commit_bytes.avgcount=1609831,commit_bytes.sum=41087428442,commit_keys.avgcount=1609831,commit_keys.sum=11637931,commit_latency.avgcount=1609831,commit_latency.avgtime=0.006236333,commit_latency.sum=10039.442388355,lease_ack_timeout=0,lease_timeout=0,new_pn=33,new_pn_latency.avgcount=33,new_pn_latency.avgtime=3.844272773,new_pn_latency.sum=126.86100151,refresh=1609856,refresh_latency.avgcount=1609856,refresh_latency.avgtime=0.005900486,refresh_latency.sum=9498.932866761,restart=109,share_state=2,share_state_bytes.avgcount=2,share_state_bytes.sum=39612,share_state_keys.avgcount=2,share_state_keys.sum=2,start_leader=22,start_peon=0,store_state=14,store_state_bytes.avgcount=14,store_state_bytes.sum=51908281,store_state_keys.avgcount=14,store_state_keys.sum=7016,store_state_latency.avgcount=14,store_state_latency.avgtime=11.668377665,store_state_latency.sum=163.357287311 1587117563000000000
ceph,collection=throttle-msgr_dispatch_throttler-mon-mgrc,host=stefanmon1,id=stefanmon1,type=monitor get=13225,get_or_fail_fail=0,get_or_fail_success=13225,get_started=0,get_sum=158700,max=104857600,put=13225,put_sum=158700,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117563000000000
ceph,collection=AsyncMessenger::Worker-1,host=stefanmon1,id=stefanmon1,type=monitor msgr_active_connections=147680,msgr_created_connections=162374,msgr_recv_bytes=29781706740,msgr_recv_messages=7170733,msgr_running_fast_dispatch_time=0,msgr_running_recv_time=1728.559151358,msgr_running_send_time=2086.681244508,msgr_running_total_time=6084.532916585,msgr_send_bytes=94062125718,msgr_send_messages=9161564 1587117563000000000
ceph,collection=throttle-msgr_dispatch_throttler-cluster,host=stefanosd1,id=0,type=osd get=281745,get_or_fail_fail=0,get_or_fail_success=281745,get_started=0,get_sum=446024457,max=104857600,put=281745,put_sum=446024457,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-bluestore_throttle_bytes,host=stefanosd1,id=0,type=osd get=275707,get_or_fail_fail=0,get_or_fail_success=0,get_started=275707,get_sum=185073179842,max=67108864,put=268870,put_sum=185073179842,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_front_server,host=stefanosd1,id=0,type=osd get=2606982,get_or_fail_fail=0,get_or_fail_success=2606982,get_started=0,get_sum=5224391928,max=104857600,put=2606982,put_sum=5224391928,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=rocksdb,host=stefanosd1,id=0,type=osd compact=0,compact_queue_len=0,compact_queue_merge=0,compact_range=0,get=1570,get_latency.avgcount=1570,get_latency.avgtime=0.000051233,get_latency.sum=0.080436788,rocksdb_write_delay_time.avgcount=0,rocksdb_write_delay_time.avgtime=0,rocksdb_write_delay_time.sum=0,rocksdb_write_memtable_time.avgcount=0,rocksdb_write_memtable_time.avgtime=0,rocksdb_write_memtable_time.sum=0,rocksdb_write_pre_and_post_time.avgcount=0,rocksdb_write_pre_and_post_time.avgtime=0,rocksdb_write_pre_and_post_time.sum=0,rocksdb_write_wal_time.avgcount=0,rocksdb_write_wal_time.avgtime=0,rocksdb_write_wal_time.sum=0,submit_latency.avgcount=275707,submit_latency.avgtime=0.000174936,submit_latency.sum=48.231345334,submit_sync_latency.avgcount=268870,submit_sync_latency.avgtime=0.006097313,submit_sync_latency.sum=1639.384555624,submit_transaction=275707,submit_transaction_sync=268870 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_back_server,host=stefanosd1,id=0,type=osd get=2606982,get_or_fail_fail=0,get_or_fail_success=2606982,get_started=0,get_sum=5224391928,max=104857600,put=2606982,put_sum=5224391928,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-objecter_bytes,host=stefanosd1,id=0,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_back_client,host=stefanosd1,id=0,type=osd get=2610285,get_or_fail_fail=0,get_or_fail_success=2610285,get_started=0,get_sum=5231011140,max=104857600,put=2610285,put_sum=5231011140,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-1,host=stefanosd1,id=0,type=osd msgr_active_connections=2093,msgr_created_connections=29142,msgr_recv_bytes=7214238199,msgr_recv_messages=3928206,msgr_running_fast_dispatch_time=171.289615064,msgr_running_recv_time=278.531155966,msgr_running_send_time=489.482588813,msgr_running_total_time=1134.004853662,msgr_send_bytes=9814725232,msgr_send_messages=3814927 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-client,host=stefanosd1,id=0,type=osd get=488206,get_or_fail_fail=0,get_or_fail_success=488206,get_started=0,get_sum=104085134,max=104857600,put=488206,put_sum=104085134,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=finisher-defered_finisher,host=stefanosd1,id=0,type=osd complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117698000000000
ceph,collection=recoverystate_perf,host=stefanosd1,id=0,type=osd activating_latency.avgcount=87,activating_latency.avgtime=0.114348341,activating_latency.sum=9.948305683,active_latency.avgcount=25,active_latency.avgtime=1790.961574431,active_latency.sum=44774.039360795,backfilling_latency.avgcount=0,backfilling_latency.avgtime=0,backfilling_latency.sum=0,clean_latency.avgcount=25,clean_latency.avgtime=1790.830827794,clean_latency.sum=44770.770694867,down_latency.avgcount=0,down_latency.avgtime=0,down_latency.sum=0,getinfo_latency.avgcount=141,getinfo_latency.avgtime=0.446233476,getinfo_latency.sum=62.918920183,getlog_latency.avgcount=87,getlog_latency.avgtime=0.007708069,getlog_latency.sum=0.670602073,getmissing_latency.avgcount=87,getmissing_latency.avgtime=0.000077594,getmissing_latency.sum=0.006750701,incomplete_latency.avgcount=0,incomplete_latency.avgtime=0,incomplete_latency.sum=0,initial_latency.avgcount=166,initial_latency.avgtime=0.001313715,initial_latency.sum=0.218076764,notbackfilling_latency.avgcount=0,notbackfilling_latency.avgtime=0,notbackfilling_latency.sum=0,notrecovering_latency.avgcount=0,notrecovering_latency.avgtime=0,notrecovering_latency.sum=0,peering_latency.avgcount=141,peering_latency.avgtime=0.948324273,peering_latency.sum=133.713722563,primary_latency.avgcount=79,primary_latency.avgtime=567.706192991,primary_latency.sum=44848.78924634,recovered_latency.avgcount=87,recovered_latency.avgtime=0.000378284,recovered_latency.sum=0.032910791,recovering_latency.avgcount=2,recovering_latency.avgtime=0.338242008,recovering_latency.sum=0.676484017,replicaactive_latency.avgcount=23,replicaactive_latency.avgtime=1790.893991295,replicaactive_latency.sum=41190.561799786,repnotrecovering_latency.avgcount=25,repnotrecovering_latency.avgtime=1647.627024984,repnotrecovering_latency.sum=41190.675624616,reprecovering_latency.avgcount=2,reprecovering_latency.avgtime=0.311884638,reprecovering_latency.sum=0.623769276,repwaitbackfillreserved_latency.avgcount=0,repwaitbackfillreserved_latency.avgtime=0,repwaitbackfillreserved_latency.sum=0,repwaitrecoveryreserved_latency.avgcount=2,repwaitrecoveryreserved_latency.avgtime=0.000462873,repwaitrecoveryreserved_latency.sum=0.000925746,reset_latency.avgcount=372,reset_latency.avgtime=0.125056393,reset_latency.sum=46.520978537,start_latency.avgcount=372,start_latency.avgtime=0.000109397,start_latency.sum=0.040695881,started_latency.avgcount=206,started_latency.avgtime=418.299777245,started_latency.sum=86169.754112641,stray_latency.avgcount=231,stray_latency.avgtime=0.98203205,stray_latency.sum=226.849403565,waitactingchange_latency.avgcount=0,waitactingchange_latency.avgtime=0,waitactingchange_latency.sum=0,waitlocalbackfillreserved_latency.avgcount=0,waitlocalbackfillreserved_latency.avgtime=0,waitlocalbackfillreserved_latency.sum=0,waitlocalrecoveryreserved_latency.avgcount=2,waitlocalrecoveryreserved_latency.avgtime=0.002802377,waitlocalrecoveryreserved_latency.sum=0.005604755,waitremotebackfillreserved_latency.avgcount=0,waitremotebackfillreserved_latency.avgtime=0,waitremotebackfillreserved_latency.sum=0,waitremoterecoveryreserved_latency.avgcount=2,waitremoterecoveryreserved_latency.avgtime=0.012855439,waitremoterecoveryreserved_latency.sum=0.025710878,waitupthru_latency.avgcount=87,waitupthru_latency.avgtime=0.805727895,waitupthru_latency.sum=70.09832695 1587117698000000000
ceph,collection=cct,host=stefanosd1,id=0,type=osd total_workers=6,unhealthy_workers=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_front_client,host=stefanosd1,id=0,type=osd get=2610285,get_or_fail_fail=0,get_or_fail_success=2610285,get_started=0,get_sum=5231011140,max=104857600,put=2610285,put_sum=5231011140,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=bluefs,host=stefanosd1,id=0,type=osd bytes_written_slow=0,bytes_written_sst=9018781,bytes_written_wal=831081573,db_total_bytes=4294967296,db_used_bytes=434110464,files_written_sst=3,files_written_wal=2,gift_bytes=0,log_bytes=134291456,log_compactions=1,logged_bytes=1101668352,max_bytes_db=1234173952,max_bytes_slow=0,max_bytes_wal=0,num_files=11,reclaim_bytes=0,slow_total_bytes=0,slow_used_bytes=0,wal_total_bytes=0,wal_used_bytes=0 1587117698000000000
ceph,collection=mempool,host=stefanosd1,id=0,type=osd bloom_filter_bytes=0,bloom_filter_items=0,bluefs_bytes=10600,bluefs_items=458,bluestore_alloc_bytes=230288,bluestore_alloc_items=28786,bluestore_cache_data_bytes=622592,bluestore_cache_data_items=43,bluestore_cache_onode_bytes=249280,bluestore_cache_onode_items=380,bluestore_cache_other_bytes=192678,bluestore_cache_other_items=20199,bluestore_fsck_bytes=0,bluestore_fsck_items=0,bluestore_txc_bytes=8272,bluestore_txc_items=11,bluestore_writing_bytes=0,bluestore_writing_deferred_bytes=670130,bluestore_writing_deferred_items=176,bluestore_writing_items=0,buffer_anon_bytes=2412465,buffer_anon_items=297,buffer_meta_bytes=5896,buffer_meta_items=67,mds_co_bytes=0,mds_co_items=0,osd_bytes=2124800,osd_items=166,osd_mapbl_bytes=155152,osd_mapbl_items=10,osd_pglog_bytes=3214704,osd_pglog_items=6288,osdmap_bytes=710892,osdmap_items=4426,osdmap_mapping_bytes=0,osdmap_mapping_items=0,pgmap_bytes=0,pgmap_items=0,unittest_1_bytes=0,unittest_1_items=0,unittest_2_bytes=0,unittest_2_items=0 1587117698000000000
ceph,collection=osd,host=stefanosd1,id=0,type=osd agent_evict=0,agent_flush=0,agent_skip=0,agent_wake=0,cached_crc=0,cached_crc_adjusted=0,copyfrom=0,heartbeat_to_peers=7,loadavg=11,map_message_epoch_dups=21,map_message_epochs=40,map_messages=31,messages_delayed_for_map=0,missed_crc=0,numpg=166,numpg_primary=62,numpg_removing=0,numpg_replica=104,numpg_stray=0,object_ctx_cache_hit=476529,object_ctx_cache_total=476536,op=476525,op_before_dequeue_op_lat.avgcount=755708,op_before_dequeue_op_lat.avgtime=0.000205759,op_before_dequeue_op_lat.sum=155.493843473,op_before_queue_op_lat.avgcount=755702,op_before_queue_op_lat.avgtime=0.000047877,op_before_queue_op_lat.sum=36.181069552,op_cache_hit=0,op_in_bytes=0,op_latency.avgcount=476525,op_latency.avgtime=0.000365956,op_latency.sum=174.387387878,op_out_bytes=10882,op_prepare_latency.avgcount=476527,op_prepare_latency.avgtime=0.000205307,op_prepare_latency.sum=97.834380034,op_process_latency.avgcount=476525,op_process_latency.avgtime=0.000139616,op_process_latency.sum=66.530847665,op_r=476521,op_r_latency.avgcount=476521,op_r_latency.avgtime=0.00036559,op_r_latency.sum=174.21148267,op_r_out_bytes=10882,op_r_prepare_latency.avgcount=476523,op_r_prepare_latency.avgtime=0.000205302,op_r_prepare_latency.sum=97.831473175,op_r_process_latency.avgcount=476521,op_r_process_latency.avgtime=0.000139396,op_r_process_latency.sum=66.425498624,op_rw=2,op_rw_in_bytes=0,op_rw_latency.avgcount=2,op_rw_latency.avgtime=0.048818975,op_rw_latency.sum=0.097637951,op_rw_out_bytes=0,op_rw_prepare_latency.avgcount=2,op_rw_prepare_latency.avgtime=0.000467887,op_rw_prepare_latency.sum=0.000935775,op_rw_process_latency.avgcount=2,op_rw_process_latency.avgtime=0.013741256,op_rw_process_latency.sum=0.027482512,op_w=2,op_w_in_bytes=0,op_w_latency.avgcount=2,op_w_latency.avgtime=0.039133628,op_w_latency.sum=0.078267257,op_w_prepare_latency.avgcount=2,op_w_prepare_latency.avgtime=0.000985542,op_w_prepare_latency.sum=0.001971084,op_w_process_latency.avgcount=2,op_w_process_latency.avgtime=0.038933264,op_w_process_latency.sum=0.077866529,op_wip=0,osd_map_bl_cache_hit=22,osd_map_bl_cache_miss=40,osd_map_cache_hit=4570,osd_map_cache_miss=15,osd_map_cache_miss_low=0,osd_map_cache_miss_low_avg.avgcount=0,osd_map_cache_miss_low_avg.sum=0,osd_pg_biginfo=2050,osd_pg_fastinfo=265780,osd_pg_info=274542,osd_tier_flush_lat.avgcount=0,osd_tier_flush_lat.avgtime=0,osd_tier_flush_lat.sum=0,osd_tier_promote_lat.avgcount=0,osd_tier_promote_lat.avgtime=0,osd_tier_promote_lat.sum=0,osd_tier_r_lat.avgcount=0,osd_tier_r_lat.avgtime=0,osd_tier_r_lat.sum=0,pull=0,push=2,push_out_bytes=10,recovery_bytes=10,recovery_ops=2,stat_bytes=107369988096,stat_bytes_avail=106271539200,stat_bytes_used=1098448896,subop=253554,subop_in_bytes=168644225,subop_latency.avgcount=253554,subop_latency.avgtime=0.0073036,subop_latency.sum=1851.857230388,subop_pull=0,subop_pull_latency.avgcount=0,subop_pull_latency.avgtime=0,subop_pull_latency.sum=0,subop_push=0,subop_push_in_bytes=0,subop_push_latency.avgcount=0,subop_push_latency.avgtime=0,subop_push_latency.sum=0,subop_w=253554,subop_w_in_bytes=168644225,subop_w_latency.avgcount=253554,subop_w_latency.avgtime=0.0073036,subop_w_latency.sum=1851.857230388,tier_clean=0,tier_delay=0,tier_dirty=0,tier_evict=0,tier_flush=0,tier_flush_fail=0,tier_promote=0,tier_proxy_read=0,tier_proxy_write=0,tier_try_flush=0,tier_try_flush_fail=0,tier_whiteout=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-ms_objecter,host=stefanosd1,id=0,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-2,host=stefanosd1,id=0,type=osd msgr_active_connections=2055,msgr_created_connections=27411,msgr_recv_bytes=6431950009,msgr_recv_messages=3552443,msgr_running_fast_dispatch_time=162.271664213,msgr_running_recv_time=254.307853033,msgr_running_send_time=503.037285799,msgr_running_total_time=1130.21070681,msgr_send_bytes=10865436237,msgr_send_messages=3523374 1587117698000000000
ceph,collection=bluestore,host=stefanosd1,id=0,type=osd bluestore_allocated=24641536,bluestore_blob_split=0,bluestore_blobs=88,bluestore_buffer_bytes=622592,bluestore_buffer_hit_bytes=160578,bluestore_buffer_miss_bytes=540236,bluestore_buffers=43,bluestore_compressed=0,bluestore_compressed_allocated=0,bluestore_compressed_original=0,bluestore_extent_compress=0,bluestore_extents=88,bluestore_fragmentation_micros=1,bluestore_gc_merged=0,bluestore_onode_hits=532102,bluestore_onode_misses=388,bluestore_onode_reshard=0,bluestore_onode_shard_hits=0,bluestore_onode_shard_misses=0,bluestore_onodes=380,bluestore_read_eio=0,bluestore_reads_with_retries=0,bluestore_stored=1987856,bluestore_txc=275707,bluestore_write_big=0,bluestore_write_big_blobs=0,bluestore_write_big_bytes=0,bluestore_write_small=60,bluestore_write_small_bytes=343843,bluestore_write_small_deferred=22,bluestore_write_small_new=38,bluestore_write_small_pre_read=22,bluestore_write_small_unused=0,commit_lat.avgcount=275707,commit_lat.avgtime=0.00699778,commit_lat.sum=1929.337103334,compress_lat.avgcount=0,compress_lat.avgtime=0,compress_lat.sum=0,compress_rejected_count=0,compress_success_count=0,csum_lat.avgcount=67,csum_lat.avgtime=0.000032601,csum_lat.sum=0.002184323,decompress_lat.avgcount=0,decompress_lat.avgtime=0,decompress_lat.sum=0,deferred_write_bytes=0,deferred_write_ops=0,kv_commit_lat.avgcount=268870,kv_commit_lat.avgtime=0.006365428,kv_commit_lat.sum=1711.472749866,kv_final_lat.avgcount=268867,kv_final_lat.avgtime=0.000043227,kv_final_lat.sum=11.622427109,kv_flush_lat.avgcount=268870,kv_flush_lat.avgtime=0.000000223,kv_flush_lat.sum=0.060141588,kv_sync_lat.avgcount=268870,kv_sync_lat.avgtime=0.006365652,kv_sync_lat.sum=1711.532891454,omap_lower_bound_lat.avgcount=2,omap_lower_bound_lat.avgtime=0.000006524,omap_lower_bound_lat.sum=0.000013048,omap_next_lat.avgcount=6704,omap_next_lat.avgtime=0.000004721,omap_next_lat.sum=0.031654097,omap_seek_to_first_lat.avgcount=323,omap_seek_to_first_lat.avgtime=0.00000522,omap_seek_to_first_lat.sum=0.00168614,omap_upper_bound_lat.avgcount=4,omap_upper_bound_lat.avgtime=0.000013086,omap_upper_bound_lat.sum=0.000052344,read_lat.avgcount=227,read_lat.avgtime=0.000699457,read_lat.sum=0.158776879,read_onode_meta_lat.avgcount=311,read_onode_meta_lat.avgtime=0.000072207,read_onode_meta_lat.sum=0.022456667,read_wait_aio_lat.avgcount=84,read_wait_aio_lat.avgtime=0.001556141,read_wait_aio_lat.sum=0.130715885,state_aio_wait_lat.avgcount=275707,state_aio_wait_lat.avgtime=0.000000345,state_aio_wait_lat.sum=0.095246457,state_deferred_aio_wait_lat.avgcount=0,state_deferred_aio_wait_lat.avgtime=0,state_deferred_aio_wait_lat.sum=0,state_deferred_cleanup_lat.avgcount=0,state_deferred_cleanup_lat.avgtime=0,state_deferred_cleanup_lat.sum=0,state_deferred_queued_lat.avgcount=0,state_deferred_queued_lat.avgtime=0,state_deferred_queued_lat.sum=0,state_done_lat.avgcount=275696,state_done_lat.avgtime=0.00000286,state_done_lat.sum=0.788700007,state_finishing_lat.avgcount=275696,state_finishing_lat.avgtime=0.000000302,state_finishing_lat.sum=0.083437168,state_io_done_lat.avgcount=275707,state_io_done_lat.avgtime=0.000001041,state_io_done_lat.sum=0.287025147,state_kv_commiting_lat.avgcount=275707,state_kv_commiting_lat.avgtime=0.006424459,state_kv_commiting_lat.sum=1771.268407864,state_kv_done_lat.avgcount=275707,state_kv_done_lat.avgtime=0.000001627,state_kv_done_lat.sum=0.448805853,state_kv_queued_lat.avgcount=275707,state_kv_queued_lat.avgtime=0.000488565,state_kv_queued_lat.sum=134.7009424,state_prepare_lat.avgcount=275707,state_prepare_lat.avgtime=0.000082464,state_prepare_lat.sum=22.736065534,submit_lat.avgcount=275707,submit_lat.avgtime=0.000120236,submit_lat.sum=33.149934412,throttle_lat.avgcount=275707,throttle_lat.avgtime=0.000001571,throttle_lat.sum=0.433185935,write_pad_bytes=151773,write_penalty_read_ops=0 1587117698000000000
ceph,collection=finisher-objecter-finisher-0,host=stefanosd1,id=0,type=osd complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117698000000000
ceph,collection=objecter,host=stefanosd1,id=0,type=osd command_active=0,command_resend=0,command_send=0,linger_active=0,linger_ping=0,linger_resend=0,linger_send=0,map_epoch=203,map_full=0,map_inc=19,omap_del=0,omap_rd=0,omap_wr=0,op=0,op_active=0,op_laggy=0,op_pg=0,op_r=0,op_reply=0,op_resend=0,op_rmw=0,op_send=0,op_send_bytes=0,op_w=0,osd_laggy=0,osd_session_close=0,osd_session_open=0,osd_sessions=0,osdop_append=0,osdop_call=0,osdop_clonerange=0,osdop_cmpxattr=0,osdop_create=0,osdop_delete=0,osdop_getxattr=0,osdop_mapext=0,osdop_notify=0,osdop_other=0,osdop_pgls=0,osdop_pgls_filter=0,osdop_read=0,osdop_resetxattrs=0,osdop_rmxattr=0,osdop_setxattr=0,osdop_sparse_read=0,osdop_src_cmpxattr=0,osdop_stat=0,osdop_truncate=0,osdop_watch=0,osdop_write=0,osdop_writefull=0,osdop_writesame=0,osdop_zero=0,poolop_active=0,poolop_resend=0,poolop_send=0,poolstat_active=0,poolstat_resend=0,poolstat_send=0,statfs_active=0,statfs_resend=0,statfs_send=0 1587117698000000000
ceph,collection=finisher-commit_finisher,host=stefanosd1,id=0,type=osd complete_latency.avgcount=11,complete_latency.avgtime=0.003447516,complete_latency.sum=0.037922681,queue_len=0 1587117698000000000
ceph,collection=throttle-objecter_ops,host=stefanosd1,id=0,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=1024,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-0,host=stefanosd1,id=0,type=osd msgr_active_connections=2128,msgr_created_connections=33685,msgr_recv_bytes=8679123051,msgr_recv_messages=4200356,msgr_running_fast_dispatch_time=151.889337454,msgr_running_recv_time=297.632294886,msgr_running_send_time=599.20020523,msgr_running_total_time=1321.361931202,msgr_send_bytes=11716202897,msgr_send_messages=4347418 1587117698000000000
ceph,collection=throttle-osd_client_bytes,host=stefanosd1,id=0,type=osd get=476554,get_or_fail_fail=0,get_or_fail_success=476554,get_started=0,get_sum=103413728,max=524288000,put=476587,put_sum=103413728,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-bluestore_throttle_deferred_bytes,host=stefanosd1,id=0,type=osd get=11,get_or_fail_fail=0,get_or_fail_success=11,get_started=0,get_sum=7723117,max=201326592,put=0,put_sum=0,take=0,take_sum=0,val=7723117,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-cluster,host=stefanosd1,id=1,type=osd get=860895,get_or_fail_fail=0,get_or_fail_success=860895,get_started=0,get_sum=596482256,max=104857600,put=860895,put_sum=596482256,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-objecter_ops,host=stefanosd1,id=1,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=1024,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-objecter_bytes,host=stefanosd1,id=1,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=finisher-defered_finisher,host=stefanosd1,id=1,type=osd complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117698000000000
ceph,collection=osd,host=stefanosd1,id=1,type=osd agent_evict=0,agent_flush=0,agent_skip=0,agent_wake=0,cached_crc=0,cached_crc_adjusted=0,copyfrom=0,heartbeat_to_peers=7,loadavg=11,map_message_epoch_dups=29,map_message_epochs=50,map_messages=39,messages_delayed_for_map=0,missed_crc=0,numpg=188,numpg_primary=71,numpg_removing=0,numpg_replica=117,numpg_stray=0,object_ctx_cache_hit=1349777,object_ctx_cache_total=2934118,op=1319230,op_before_dequeue_op_lat.avgcount=3792053,op_before_dequeue_op_lat.avgtime=0.000405802,op_before_dequeue_op_lat.sum=1538.826381623,op_before_queue_op_lat.avgcount=3778690,op_before_queue_op_lat.avgtime=0.000033273,op_before_queue_op_lat.sum=125.731131596,op_cache_hit=0,op_in_bytes=0,op_latency.avgcount=1319230,op_latency.avgtime=0.002858138,op_latency.sum=3770.541581676,op_out_bytes=1789210,op_prepare_latency.avgcount=1336472,op_prepare_latency.avgtime=0.000279458,op_prepare_latency.sum=373.488913339,op_process_latency.avgcount=1319230,op_process_latency.avgtime=0.002666408,op_process_latency.sum=3517.606407526,op_r=1075394,op_r_latency.avgcount=1075394,op_r_latency.avgtime=0.000303779,op_r_latency.sum=326.682443032,op_r_out_bytes=1789210,op_r_prepare_latency.avgcount=1075394,op_r_prepare_latency.avgtime=0.000171228,op_r_prepare_latency.sum=184.138580631,op_r_process_latency.avgcount=1075394,op_r_process_latency.avgtime=0.00011609,op_r_process_latency.sum=124.842894319,op_rw=243832,op_rw_in_bytes=0,op_rw_latency.avgcount=243832,op_rw_latency.avgtime=0.014123636,op_rw_latency.sum=3443.79445124,op_rw_out_bytes=0,op_rw_prepare_latency.avgcount=261072,op_rw_prepare_latency.avgtime=0.000725265,op_rw_prepare_latency.sum=189.346543463,op_rw_process_latency.avgcount=243832,op_rw_process_latency.avgtime=0.013914089,op_rw_process_latency.sum=3392.700241086,op_w=4,op_w_in_bytes=0,op_w_latency.avgcount=4,op_w_latency.avgtime=0.016171851,op_w_latency.sum=0.064687404,op_w_prepare_latency.avgcount=6,op_w_prepare_latency.avgtime=0.00063154,op_w_prepare_latency.sum=0.003789245,op_w_process_latency.avgcount=4,op_w_process_latency.avgtime=0.01581803,op_w_process_latency.sum=0.063272121,op_wip=0,osd_map_bl_cache_hit=36,osd_map_bl_cache_miss=40,osd_map_cache_hit=5404,osd_map_cache_miss=14,osd_map_cache_miss_low=0,osd_map_cache_miss_low_avg.avgcount=0,osd_map_cache_miss_low_avg.sum=0,osd_pg_biginfo=2333,osd_pg_fastinfo=576157,osd_pg_info=591751,osd_tier_flush_lat.avgcount=0,osd_tier_flush_lat.avgtime=0,osd_tier_flush_lat.sum=0,osd_tier_promote_lat.avgcount=0,osd_tier_promote_lat.avgtime=0,osd_tier_promote_lat.sum=0,osd_tier_r_lat.avgcount=0,osd_tier_r_lat.avgtime=0,osd_tier_r_lat.sum=0,pull=0,push=22,push_out_bytes=0,recovery_bytes=0,recovery_ops=21,stat_bytes=107369988096,stat_bytes_avail=106271997952,stat_bytes_used=1097990144,subop=306946,subop_in_bytes=204236742,subop_latency.avgcount=306946,subop_latency.avgtime=0.006744881,subop_latency.sum=2070.314452989,subop_pull=0,subop_pull_latency.avgcount=0,subop_pull_latency.avgtime=0,subop_pull_latency.sum=0,subop_push=0,subop_push_in_bytes=0,subop_push_latency.avgcount=0,subop_push_latency.avgtime=0,subop_push_latency.sum=0,subop_w=306946,subop_w_in_bytes=204236742,subop_w_latency.avgcount=306946,subop_w_latency.avgtime=0.006744881,subop_w_latency.sum=2070.314452989,tier_clean=0,tier_delay=0,tier_dirty=8,tier_evict=0,tier_flush=0,tier_flush_fail=0,tier_promote=0,tier_proxy_read=0,tier_proxy_write=0,tier_try_flush=0,tier_try_flush_fail=0,tier_whiteout=0 1587117698000000000
ceph,collection=objecter,host=stefanosd1,id=1,type=osd command_active=0,command_resend=0,command_send=0,linger_active=0,linger_ping=0,linger_resend=0,linger_send=0,map_epoch=203,map_full=0,map_inc=19,omap_del=0,omap_rd=0,omap_wr=0,op=0,op_active=0,op_laggy=0,op_pg=0,op_r=0,op_reply=0,op_resend=0,op_rmw=0,op_send=0,op_send_bytes=0,op_w=0,osd_laggy=0,osd_session_close=0,osd_session_open=0,osd_sessions=0,osdop_append=0,osdop_call=0,osdop_clonerange=0,osdop_cmpxattr=0,osdop_create=0,osdop_delete=0,osdop_getxattr=0,osdop_mapext=0,osdop_notify=0,osdop_other=0,osdop_pgls=0,osdop_pgls_filter=0,osdop_read=0,osdop_resetxattrs=0,osdop_rmxattr=0,osdop_setxattr=0,osdop_sparse_read=0,osdop_src_cmpxattr=0,osdop_stat=0,osdop_truncate=0,osdop_watch=0,osdop_write=0,osdop_writefull=0,osdop_writesame=0,osdop_zero=0,poolop_active=0,poolop_resend=0,poolop_send=0,poolstat_active=0,poolstat_resend=0,poolstat_send=0,statfs_active=0,statfs_resend=0,statfs_send=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-0,host=stefanosd1,id=1,type=osd msgr_active_connections=1356,msgr_created_connections=12290,msgr_recv_bytes=8577187219,msgr_recv_messages=6387040,msgr_running_fast_dispatch_time=475.903632306,msgr_running_recv_time=425.937196699,msgr_running_send_time=783.676217521,msgr_running_total_time=1989.242459076,msgr_send_bytes=12583034449,msgr_send_messages=6074344 1587117698000000000
ceph,collection=bluestore,host=stefanosd1,id=1,type=osd bluestore_allocated=24182784,bluestore_blob_split=0,bluestore_blobs=88,bluestore_buffer_bytes=614400,bluestore_buffer_hit_bytes=142047,bluestore_buffer_miss_bytes=541480,bluestore_buffers=41,bluestore_compressed=0,bluestore_compressed_allocated=0,bluestore_compressed_original=0,bluestore_extent_compress=0,bluestore_extents=88,bluestore_fragmentation_micros=1,bluestore_gc_merged=0,bluestore_onode_hits=1403948,bluestore_onode_misses=1584732,bluestore_onode_reshard=0,bluestore_onode_shard_hits=0,bluestore_onode_shard_misses=0,bluestore_onodes=459,bluestore_read_eio=0,bluestore_reads_with_retries=0,bluestore_stored=1985647,bluestore_txc=593150,bluestore_write_big=0,bluestore_write_big_blobs=0,bluestore_write_big_bytes=0,bluestore_write_small=58,bluestore_write_small_bytes=343091,bluestore_write_small_deferred=20,bluestore_write_small_new=38,bluestore_write_small_pre_read=20,bluestore_write_small_unused=0,commit_lat.avgcount=593150,commit_lat.avgtime=0.006514834,commit_lat.sum=3864.274280733,compress_lat.avgcount=0,compress_lat.avgtime=0,compress_lat.sum=0,compress_rejected_count=0,compress_success_count=0,csum_lat.avgcount=60,csum_lat.avgtime=0.000028258,csum_lat.sum=0.001695512,decompress_lat.avgcount=0,decompress_lat.avgtime=0,decompress_lat.sum=0,deferred_write_bytes=0,deferred_write_ops=0,kv_commit_lat.avgcount=578129,kv_commit_lat.avgtime=0.00570707,kv_commit_lat.sum=3299.423186928,kv_final_lat.avgcount=578124,kv_final_lat.avgtime=0.000042752,kv_final_lat.sum=24.716171934,kv_flush_lat.avgcount=578129,kv_flush_lat.avgtime=0.000000209,kv_flush_lat.sum=0.121169044,kv_sync_lat.avgcount=578129,kv_sync_lat.avgtime=0.00570728,kv_sync_lat.sum=3299.544355972,omap_lower_bound_lat.avgcount=22,omap_lower_bound_lat.avgtime=0.000005979,omap_lower_bound_lat.sum=0.000131539,omap_next_lat.avgcount=13248,omap_next_lat.avgtime=0.000004836,omap_next_lat.sum=0.064077797,omap_seek_to_first_lat.avgcount=525,omap_seek_to_first_lat.avgtime=0.000004906,omap_seek_to_first_lat.sum=0.002575786,omap_upper_bound_lat.avgcount=0,omap_upper_bound_lat.avgtime=0,omap_upper_bound_lat.sum=0,read_lat.avgcount=406,read_lat.avgtime=0.000383254,read_lat.sum=0.155601529,read_onode_meta_lat.avgcount=483,read_onode_meta_lat.avgtime=0.000008805,read_onode_meta_lat.sum=0.004252832,read_wait_aio_lat.avgcount=77,read_wait_aio_lat.avgtime=0.001907361,read_wait_aio_lat.sum=0.146866799,state_aio_wait_lat.avgcount=593150,state_aio_wait_lat.avgtime=0.000000388,state_aio_wait_lat.sum=0.230498048,state_deferred_aio_wait_lat.avgcount=0,state_deferred_aio_wait_lat.avgtime=0,state_deferred_aio_wait_lat.sum=0,state_deferred_cleanup_lat.avgcount=0,state_deferred_cleanup_lat.avgtime=0,state_deferred_cleanup_lat.sum=0,state_deferred_queued_lat.avgcount=0,state_deferred_queued_lat.avgtime=0,state_deferred_queued_lat.sum=0,state_done_lat.avgcount=593140,state_done_lat.avgtime=0.000003048,state_done_lat.sum=1.80789161,state_finishing_lat.avgcount=593140,state_finishing_lat.avgtime=0.000000325,state_finishing_lat.sum=0.192952339,state_io_done_lat.avgcount=593150,state_io_done_lat.avgtime=0.000001202,state_io_done_lat.sum=0.713333116,state_kv_commiting_lat.avgcount=593150,state_kv_commiting_lat.avgtime=0.005788541,state_kv_commiting_lat.sum=3433.473378536,state_kv_done_lat.avgcount=593150,state_kv_done_lat.avgtime=0.000001472,state_kv_done_lat.sum=0.873559611,state_kv_queued_lat.avgcount=593150,state_kv_queued_lat.avgtime=0.000634215,state_kv_queued_lat.sum=376.18491577,state_prepare_lat.avgcount=593150,state_prepare_lat.avgtime=0.000089694,state_prepare_lat.sum=53.202464675,submit_lat.avgcount=593150,submit_lat.avgtime=0.000127856,submit_lat.sum=75.83816759,throttle_lat.avgcount=593150,throttle_lat.avgtime=0.000001726,throttle_lat.sum=1.023832181,write_pad_bytes=144333,write_penalty_read_ops=0 1587117698000000000
ceph,collection=throttle-osd_client_bytes,host=stefanosd1,id=1,type=osd get=2920772,get_or_fail_fail=0,get_or_fail_success=2920772,get_started=0,get_sum=739935873,max=524288000,put=4888498,put_sum=739935873,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_front_client,host=stefanosd1,id=1,type=osd get=2605442,get_or_fail_fail=0,get_or_fail_success=2605442,get_started=0,get_sum=5221305768,max=104857600,put=2605442,put_sum=5221305768,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-2,host=stefanosd1,id=1,type=osd msgr_active_connections=1375,msgr_created_connections=12689,msgr_recv_bytes=6393440855,msgr_recv_messages=3260458,msgr_running_fast_dispatch_time=120.622437418,msgr_running_recv_time=225.24709441,msgr_running_send_time=499.150587343,msgr_running_total_time=1043.340296846,msgr_send_bytes=11134862571,msgr_send_messages=3450760 1587117698000000000
ceph,collection=bluefs,host=stefanosd1,id=1,type=osd bytes_written_slow=0,bytes_written_sst=19824993,bytes_written_wal=1788507023,db_total_bytes=4294967296,db_used_bytes=522190848,files_written_sst=4,files_written_wal=2,gift_bytes=0,log_bytes=1056768,log_compactions=2,logged_bytes=1933271040,max_bytes_db=1483735040,max_bytes_slow=0,max_bytes_wal=0,num_files=12,reclaim_bytes=0,slow_total_bytes=0,slow_used_bytes=0,wal_total_bytes=0,wal_used_bytes=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_back_client,host=stefanosd1,id=1,type=osd get=2605442,get_or_fail_fail=0,get_or_fail_success=2605442,get_started=0,get_sum=5221305768,max=104857600,put=2605442,put_sum=5221305768,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-bluestore_throttle_deferred_bytes,host=stefanosd1,id=1,type=osd get=10,get_or_fail_fail=0,get_or_fail_success=10,get_started=0,get_sum=7052009,max=201326592,put=0,put_sum=0,take=0,take_sum=0,val=7052009,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=rocksdb,host=stefanosd1,id=1,type=osd compact=0,compact_queue_len=0,compact_queue_merge=0,compact_range=0,get=1586061,get_latency.avgcount=1586061,get_latency.avgtime=0.000083009,get_latency.sum=131.658296684,rocksdb_write_delay_time.avgcount=0,rocksdb_write_delay_time.avgtime=0,rocksdb_write_delay_time.sum=0,rocksdb_write_memtable_time.avgcount=0,rocksdb_write_memtable_time.avgtime=0,rocksdb_write_memtable_time.sum=0,rocksdb_write_pre_and_post_time.avgcount=0,rocksdb_write_pre_and_post_time.avgtime=0,rocksdb_write_pre_and_post_time.sum=0,rocksdb_write_wal_time.avgcount=0,rocksdb_write_wal_time.avgtime=0,rocksdb_write_wal_time.sum=0,submit_latency.avgcount=593150,submit_latency.avgtime=0.000172072,submit_latency.sum=102.064900673,submit_sync_latency.avgcount=578129,submit_sync_latency.avgtime=0.005447017,submit_sync_latency.sum=3149.078822012,submit_transaction=593150,submit_transaction_sync=578129 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_back_server,host=stefanosd1,id=1,type=osd get=2607669,get_or_fail_fail=0,get_or_fail_success=2607669,get_started=0,get_sum=5225768676,max=104857600,put=2607669,put_sum=5225768676,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=recoverystate_perf,host=stefanosd1,id=1,type=osd activating_latency.avgcount=104,activating_latency.avgtime=0.071646485,activating_latency.sum=7.451234493,active_latency.avgcount=33,active_latency.avgtime=1734.369034268,active_latency.sum=57234.178130859,backfilling_latency.avgcount=1,backfilling_latency.avgtime=2.598401698,backfilling_latency.sum=2.598401698,clean_latency.avgcount=33,clean_latency.avgtime=1734.213467342,clean_latency.sum=57229.044422292,down_latency.avgcount=0,down_latency.avgtime=0,down_latency.sum=0,getinfo_latency.avgcount=167,getinfo_latency.avgtime=0.373444627,getinfo_latency.sum=62.365252849,getlog_latency.avgcount=105,getlog_latency.avgtime=0.003575062,getlog_latency.sum=0.375381569,getmissing_latency.avgcount=104,getmissing_latency.avgtime=0.000157091,getmissing_latency.sum=0.016337565,incomplete_latency.avgcount=0,incomplete_latency.avgtime=0,incomplete_latency.sum=0,initial_latency.avgcount=188,initial_latency.avgtime=0.001833512,initial_latency.sum=0.344700343,notbackfilling_latency.avgcount=0,notbackfilling_latency.avgtime=0,notbackfilling_latency.sum=0,notrecovering_latency.avgcount=0,notrecovering_latency.avgtime=0,notrecovering_latency.sum=0,peering_latency.avgcount=167,peering_latency.avgtime=1.501818082,peering_latency.sum=250.803619796,primary_latency.avgcount=97,primary_latency.avgtime=591.344286378,primary_latency.sum=57360.395778762,recovered_latency.avgcount=104,recovered_latency.avgtime=0.000291138,recovered_latency.sum=0.030278433,recovering_latency.avgcount=2,recovering_latency.avgtime=0.142378096,recovering_latency.sum=0.284756192,replicaactive_latency.avgcount=32,replicaactive_latency.avgtime=1788.474901442,replicaactive_latency.sum=57231.196846165,repnotrecovering_latency.avgcount=34,repnotrecovering_latency.avgtime=1683.273587087,repnotrecovering_latency.sum=57231.301960987,reprecovering_latency.avgcount=2,reprecovering_latency.avgtime=0.418094818,reprecovering_latency.sum=0.836189637,repwaitbackfillreserved_latency.avgcount=0,repwaitbackfillreserved_latency.avgtime=0,repwaitbackfillreserved_latency.sum=0,repwaitrecoveryreserved_latency.avgcount=2,repwaitrecoveryreserved_latency.avgtime=0.000588413,repwaitrecoveryreserved_latency.sum=0.001176827,reset_latency.avgcount=433,reset_latency.avgtime=0.15669689,reset_latency.sum=67.849753631,start_latency.avgcount=433,start_latency.avgtime=0.000412707,start_latency.sum=0.178702508,started_latency.avgcount=245,started_latency.avgtime=468.419544137,started_latency.sum=114762.788313581,stray_latency.avgcount=266,stray_latency.avgtime=1.489291271,stray_latency.sum=396.151478238,waitactingchange_latency.avgcount=1,waitactingchange_latency.avgtime=0.982689906,waitactingchange_latency.sum=0.982689906,waitlocalbackfillreserved_latency.avgcount=1,waitlocalbackfillreserved_latency.avgtime=0.000542092,waitlocalbackfillreserved_latency.sum=0.000542092,waitlocalrecoveryreserved_latency.avgcount=2,waitlocalrecoveryreserved_latency.avgtime=0.00391669,waitlocalrecoveryreserved_latency.sum=0.007833381,waitremotebackfillreserved_latency.avgcount=1,waitremotebackfillreserved_latency.avgtime=0.003110409,waitremotebackfillreserved_latency.sum=0.003110409,waitremoterecoveryreserved_latency.avgcount=2,waitremoterecoveryreserved_latency.avgtime=0.012229338,waitremoterecoveryreserved_latency.sum=0.024458677,waitupthru_latency.avgcount=104,waitupthru_latency.avgtime=1.807608905,waitupthru_latency.sum=187.991326197 1587117698000000000
ceph,collection=AsyncMessenger::Worker-1,host=stefanosd1,id=1,type=osd msgr_active_connections=1289,msgr_created_connections=9469,msgr_recv_bytes=8348149800,msgr_recv_messages=5048791,msgr_running_fast_dispatch_time=313.754567889,msgr_running_recv_time=372.054833029,msgr_running_send_time=694.900405016,msgr_running_total_time=1656.294769387,msgr_send_bytes=11550148208,msgr_send_messages=5175962 1587117698000000000
ceph,collection=throttle-bluestore_throttle_bytes,host=stefanosd1,id=1,type=osd get=593150,get_or_fail_fail=0,get_or_fail_success=0,get_started=593150,get_sum=398147414260,max=67108864,put=578129,put_sum=398147414260,take=0,take_sum=0,val=0,wait.avgcount=29,wait.avgtime=0.000972655,wait.sum=0.028207005 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-ms_objecter,host=stefanosd1,id=1,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=cct,host=stefanosd1,id=1,type=osd total_workers=6,unhealthy_workers=0 1587117698000000000
ceph,collection=mempool,host=stefanosd1,id=1,type=osd bloom_filter_bytes=0,bloom_filter_items=0,bluefs_bytes=13064,bluefs_items=593,bluestore_alloc_bytes=230288,bluestore_alloc_items=28786,bluestore_cache_data_bytes=614400,bluestore_cache_data_items=41,bluestore_cache_onode_bytes=301104,bluestore_cache_onode_items=459,bluestore_cache_other_bytes=230945,bluestore_cache_other_items=26119,bluestore_fsck_bytes=0,bluestore_fsck_items=0,bluestore_txc_bytes=7520,bluestore_txc_items=10,bluestore_writing_bytes=0,bluestore_writing_deferred_bytes=657768,bluestore_writing_deferred_items=172,bluestore_writing_items=0,buffer_anon_bytes=2328515,buffer_anon_items=271,buffer_meta_bytes=5808,buffer_meta_items=66,mds_co_bytes=0,mds_co_items=0,osd_bytes=2406400,osd_items=188,osd_mapbl_bytes=139623,osd_mapbl_items=9,osd_pglog_bytes=6768784,osd_pglog_items=18179,osdmap_bytes=710892,osdmap_items=4426,osdmap_mapping_bytes=0,osdmap_mapping_items=0,pgmap_bytes=0,pgmap_items=0,unittest_1_bytes=0,unittest_1_items=0,unittest_2_bytes=0,unittest_2_items=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-client,host=stefanosd1,id=1,type=osd get=2932513,get_or_fail_fail=0,get_or_fail_success=2932513,get_started=0,get_sum=740620215,max=104857600,put=2932513,put_sum=740620215,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_front_server,host=stefanosd1,id=1,type=osd get=2607669,get_or_fail_fail=0,get_or_fail_success=2607669,get_started=0,get_sum=5225768676,max=104857600,put=2607669,put_sum=5225768676,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=finisher-commit_finisher,host=stefanosd1,id=1,type=osd complete_latency.avgcount=10,complete_latency.avgtime=0.002884646,complete_latency.sum=0.028846469,queue_len=0 1587117698000000000
ceph,collection=finisher-objecter-finisher-0,host=stefanosd1,id=1,type=osd complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117698000000000
ceph,collection=throttle-objecter_bytes,host=stefanosd1,id=2,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=finisher-commit_finisher,host=stefanosd1,id=2,type=osd complete_latency.avgcount=11,complete_latency.avgtime=0.002714416,complete_latency.sum=0.029858583,queue_len=0 1587117698000000000
ceph,collection=finisher-defered_finisher,host=stefanosd1,id=2,type=osd complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117698000000000
ceph,collection=objecter,host=stefanosd1,id=2,type=osd command_active=0,command_resend=0,command_send=0,linger_active=0,linger_ping=0,linger_resend=0,linger_send=0,map_epoch=203,map_full=0,map_inc=19,omap_del=0,omap_rd=0,omap_wr=0,op=0,op_active=0,op_laggy=0,op_pg=0,op_r=0,op_reply=0,op_resend=0,op_rmw=0,op_send=0,op_send_bytes=0,op_w=0,osd_laggy=0,osd_session_close=0,osd_session_open=0,osd_sessions=0,osdop_append=0,osdop_call=0,osdop_clonerange=0,osdop_cmpxattr=0,osdop_create=0,osdop_delete=0,osdop_getxattr=0,osdop_mapext=0,osdop_notify=0,osdop_other=0,osdop_pgls=0,osdop_pgls_filter=0,osdop_read=0,osdop_resetxattrs=0,osdop_rmxattr=0,osdop_setxattr=0,osdop_sparse_read=0,osdop_src_cmpxattr=0,osdop_stat=0,osdop_truncate=0,osdop_watch=0,osdop_write=0,osdop_writefull=0,osdop_writesame=0,osdop_zero=0,poolop_active=0,poolop_resend=0,poolop_send=0,poolstat_active=0,poolstat_resend=0,poolstat_send=0,statfs_active=0,statfs_resend=0,statfs_send=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_back_client,host=stefanosd1,id=2,type=osd get=2607136,get_or_fail_fail=0,get_or_fail_success=2607136,get_started=0,get_sum=5224700544,max=104857600,put=2607136,put_sum=5224700544,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=mempool,host=stefanosd1,id=2,type=osd bloom_filter_bytes=0,bloom_filter_items=0,bluefs_bytes=11624,bluefs_items=522,bluestore_alloc_bytes=230288,bluestore_alloc_items=28786,bluestore_cache_data_bytes=614400,bluestore_cache_data_items=41,bluestore_cache_onode_bytes=228288,bluestore_cache_onode_items=348,bluestore_cache_other_bytes=174158,bluestore_cache_other_items=18527,bluestore_fsck_bytes=0,bluestore_fsck_items=0,bluestore_txc_bytes=8272,bluestore_txc_items=11,bluestore_writing_bytes=0,bluestore_writing_deferred_bytes=670130,bluestore_writing_deferred_items=176,bluestore_writing_items=0,buffer_anon_bytes=2311664,buffer_anon_items=244,buffer_meta_bytes=5456,buffer_meta_items=62,mds_co_bytes=0,mds_co_items=0,osd_bytes=1920000,osd_items=150,osd_mapbl_bytes=155152,osd_mapbl_items=10,osd_pglog_bytes=3393520,osd_pglog_items=9128,osdmap_bytes=710892,osdmap_items=4426,osdmap_mapping_bytes=0,osdmap_mapping_items=0,pgmap_bytes=0,pgmap_items=0,unittest_1_bytes=0,unittest_1_items=0,unittest_2_bytes=0,unittest_2_items=0 1587117698000000000
ceph,collection=osd,host=stefanosd1,id=2,type=osd agent_evict=0,agent_flush=0,agent_skip=0,agent_wake=0,cached_crc=0,cached_crc_adjusted=0,copyfrom=0,heartbeat_to_peers=7,loadavg=11,map_message_epoch_dups=37,map_message_epochs=56,map_messages=37,messages_delayed_for_map=0,missed_crc=0,numpg=150,numpg_primary=59,numpg_removing=0,numpg_replica=91,numpg_stray=0,object_ctx_cache_hit=705923,object_ctx_cache_total=705951,op=690584,op_before_dequeue_op_lat.avgcount=1155697,op_before_dequeue_op_lat.avgtime=0.000217926,op_before_dequeue_op_lat.sum=251.856487141,op_before_queue_op_lat.avgcount=1148445,op_before_queue_op_lat.avgtime=0.000039696,op_before_queue_op_lat.sum=45.589516462,op_cache_hit=0,op_in_bytes=0,op_latency.avgcount=690584,op_latency.avgtime=0.002488685,op_latency.sum=1718.646504654,op_out_bytes=1026000,op_prepare_latency.avgcount=698700,op_prepare_latency.avgtime=0.000300375,op_prepare_latency.sum=209.872029659,op_process_latency.avgcount=690584,op_process_latency.avgtime=0.00230742,op_process_latency.sum=1593.46739165,op_r=548020,op_r_latency.avgcount=548020,op_r_latency.avgtime=0.000298287,op_r_latency.sum=163.467760649,op_r_out_bytes=1026000,op_r_prepare_latency.avgcount=548020,op_r_prepare_latency.avgtime=0.000186359,op_r_prepare_latency.sum=102.128629183,op_r_process_latency.avgcount=548020,op_r_process_latency.avgtime=0.00012716,op_r_process_latency.sum=69.686468884,op_rw=142562,op_rw_in_bytes=0,op_rw_latency.avgcount=142562,op_rw_latency.avgtime=0.010908597,op_rw_latency.sum=1555.151525732,op_rw_out_bytes=0,op_rw_prepare_latency.avgcount=150678,op_rw_prepare_latency.avgtime=0.000715043,op_rw_prepare_latency.sum=107.741399304,op_rw_process_latency.avgcount=142562,op_rw_process_latency.avgtime=0.01068836,op_rw_process_latency.sum=1523.754107887,op_w=2,op_w_in_bytes=0,op_w_latency.avgcount=2,op_w_latency.avgtime=0.013609136,op_w_latency.sum=0.027218273,op_w_prepare_latency.avgcount=2,op_w_prepare_latency.avgtime=0.001000586,op_w_prepare_latency.sum=0.002001172,op_w_process_latency.avgcount=2,op_w_process_latency.avgtime=0.013407439,op_w_process_latency.sum=0.026814879,op_wip=0,osd_map_bl_cache_hit=15,osd_map_bl_cache_miss=41,osd_map_cache_hit=4241,osd_map_cache_miss=14,osd_map_cache_miss_low=0,osd_map_cache_miss_low_avg.avgcount=0,osd_map_cache_miss_low_avg.sum=0,osd_pg_biginfo=1824,osd_pg_fastinfo=285998,osd_pg_info=294869,osd_tier_flush_lat.avgcount=0,osd_tier_flush_lat.avgtime=0,osd_tier_flush_lat.sum=0,osd_tier_promote_lat.avgcount=0,osd_tier_promote_lat.avgtime=0,osd_tier_promote_lat.sum=0,osd_tier_r_lat.avgcount=0,osd_tier_r_lat.avgtime=0,osd_tier_r_lat.sum=0,pull=0,push=1,push_out_bytes=0,recovery_bytes=0,recovery_ops=0,stat_bytes=107369988096,stat_bytes_avail=106271932416,stat_bytes_used=1098055680,subop=134165,subop_in_bytes=89501237,subop_latency.avgcount=134165,subop_latency.avgtime=0.007313523,subop_latency.sum=981.218888627,subop_pull=0,subop_pull_latency.avgcount=0,subop_pull_latency.avgtime=0,subop_pull_latency.sum=0,subop_push=0,subop_push_in_bytes=0,subop_push_latency.avgcount=0,subop_push_latency.avgtime=0,subop_push_latency.sum=0,subop_w=134165,subop_w_in_bytes=89501237,subop_w_latency.avgcount=134165,subop_w_latency.avgtime=0.007313523,subop_w_latency.sum=981.218888627,tier_clean=0,tier_delay=0,tier_dirty=4,tier_evict=0,tier_flush=0,tier_flush_fail=0,tier_promote=0,tier_proxy_read=0,tier_proxy_write=0,tier_try_flush=0,tier_try_flush_fail=0,tier_whiteout=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-1,host=stefanosd1,id=2,type=osd msgr_active_connections=746,msgr_created_connections=15212,msgr_recv_bytes=8633229006,msgr_recv_messages=4284202,msgr_running_fast_dispatch_time=153.820479102,msgr_running_recv_time=282.031655658,msgr_running_send_time=585.444749736,msgr_running_total_time=1231.431789242,msgr_send_bytes=11962769351,msgr_send_messages=4440622 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-ms_objecter,host=stefanosd1,id=2,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_front_client,host=stefanosd1,id=2,type=osd get=2607136,get_or_fail_fail=0,get_or_fail_success=2607136,get_started=0,get_sum=5224700544,max=104857600,put=2607136,put_sum=5224700544,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=bluefs,host=stefanosd1,id=2,type=osd bytes_written_slow=0,bytes_written_sst=9065815,bytes_written_wal=901884611,db_total_bytes=4294967296,db_used_bytes=546308096,files_written_sst=3,files_written_wal=2,gift_bytes=0,log_bytes=225726464,log_compactions=1,logged_bytes=1195945984,max_bytes_db=1234173952,max_bytes_slow=0,max_bytes_wal=0,num_files=11,reclaim_bytes=0,slow_total_bytes=0,slow_used_bytes=0,wal_total_bytes=0,wal_used_bytes=0 1587117698000000000
ceph,collection=recoverystate_perf,host=stefanosd1,id=2,type=osd activating_latency.avgcount=88,activating_latency.avgtime=0.086149065,activating_latency.sum=7.581117751,active_latency.avgcount=29,active_latency.avgtime=1790.849396082,active_latency.sum=51934.632486379,backfilling_latency.avgcount=0,backfilling_latency.avgtime=0,backfilling_latency.sum=0,clean_latency.avgcount=29,clean_latency.avgtime=1790.754765195,clean_latency.sum=51931.888190683,down_latency.avgcount=0,down_latency.avgtime=0,down_latency.sum=0,getinfo_latency.avgcount=134,getinfo_latency.avgtime=0.427567953,getinfo_latency.sum=57.294105786,getlog_latency.avgcount=88,getlog_latency.avgtime=0.011810192,getlog_latency.sum=1.03929697,getmissing_latency.avgcount=88,getmissing_latency.avgtime=0.000104598,getmissing_latency.sum=0.009204673,incomplete_latency.avgcount=0,incomplete_latency.avgtime=0,incomplete_latency.sum=0,initial_latency.avgcount=150,initial_latency.avgtime=0.001251361,initial_latency.sum=0.187704197,notbackfilling_latency.avgcount=0,notbackfilling_latency.avgtime=0,notbackfilling_latency.sum=0,notrecovering_latency.avgcount=0,notrecovering_latency.avgtime=0,notrecovering_latency.sum=0,peering_latency.avgcount=134,peering_latency.avgtime=0.998405763,peering_latency.sum=133.786372331,primary_latency.avgcount=75,primary_latency.avgtime=693.473306562,primary_latency.sum=52010.497992212,recovered_latency.avgcount=88,recovered_latency.avgtime=0.000609715,recovered_latency.sum=0.053654964,recovering_latency.avgcount=1,recovering_latency.avgtime=0.100713031,recovering_latency.sum=0.100713031,replicaactive_latency.avgcount=21,replicaactive_latency.avgtime=1790.852354921,replicaactive_latency.sum=37607.89945336,repnotrecovering_latency.avgcount=21,repnotrecovering_latency.avgtime=1790.852315529,repnotrecovering_latency.sum=37607.898626121,reprecovering_latency.avgcount=0,reprecovering_latency.avgtime=0,reprecovering_latency.sum=0,repwaitbackfillreserved_latency.avgcount=0,repwaitbackfillreserved_latency.avgtime=0,repwaitbackfillreserved_latency.sum=0,repwaitrecoveryreserved_latency.avgcount=0,repwaitrecoveryreserved_latency.avgtime=0,repwaitrecoveryreserved_latency.sum=0,reset_latency.avgcount=346,reset_latency.avgtime=0.126826803,reset_latency.sum=43.882073917,start_latency.avgcount=346,start_latency.avgtime=0.000233277,start_latency.sum=0.080713962,started_latency.avgcount=196,started_latency.avgtime=457.885378797,started_latency.sum=89745.534244237,stray_latency.avgcount=212,stray_latency.avgtime=1.013774396,stray_latency.sum=214.920172121,waitactingchange_latency.avgcount=0,waitactingchange_latency.avgtime=0,waitactingchange_latency.sum=0,waitlocalbackfillreserved_latency.avgcount=0,waitlocalbackfillreserved_latency.avgtime=0,waitlocalbackfillreserved_latency.sum=0,waitlocalrecoveryreserved_latency.avgcount=1,waitlocalrecoveryreserved_latency.avgtime=0.001572379,waitlocalrecoveryreserved_latency.sum=0.001572379,waitremotebackfillreserved_latency.avgcount=0,waitremotebackfillreserved_latency.avgtime=0,waitremotebackfillreserved_latency.sum=0,waitremoterecoveryreserved_latency.avgcount=1,waitremoterecoveryreserved_latency.avgtime=0.012729633,waitremoterecoveryreserved_latency.sum=0.012729633,waitupthru_latency.avgcount=88,waitupthru_latency.avgtime=0.857137729,waitupthru_latency.sum=75.428120205 1587117698000000000
ceph,collection=throttle-objecter_ops,host=stefanosd1,id=2,type=osd get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=1024,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=bluestore,host=stefanosd1,id=2,type=osd bluestore_allocated=24248320,bluestore_blob_split=0,bluestore_blobs=83,bluestore_buffer_bytes=614400,bluestore_buffer_hit_bytes=161362,bluestore_buffer_miss_bytes=534799,bluestore_buffers=41,bluestore_compressed=0,bluestore_compressed_allocated=0,bluestore_compressed_original=0,bluestore_extent_compress=0,bluestore_extents=83,bluestore_fragmentation_micros=1,bluestore_gc_merged=0,bluestore_onode_hits=723852,bluestore_onode_misses=364,bluestore_onode_reshard=0,bluestore_onode_shard_hits=0,bluestore_onode_shard_misses=0,bluestore_onodes=348,bluestore_read_eio=0,bluestore_reads_with_retries=0,bluestore_stored=1984402,bluestore_txc=295997,bluestore_write_big=0,bluestore_write_big_blobs=0,bluestore_write_big_bytes=0,bluestore_write_small=60,bluestore_write_small_bytes=343843,bluestore_write_small_deferred=22,bluestore_write_small_new=38,bluestore_write_small_pre_read=22,bluestore_write_small_unused=0,commit_lat.avgcount=295997,commit_lat.avgtime=0.006994931,commit_lat.sum=2070.478673619,compress_lat.avgcount=0,compress_lat.avgtime=0,compress_lat.sum=0,compress_rejected_count=0,compress_success_count=0,csum_lat.avgcount=47,csum_lat.avgtime=0.000034434,csum_lat.sum=0.001618423,decompress_lat.avgcount=0,decompress_lat.avgtime=0,decompress_lat.sum=0,deferred_write_bytes=0,deferred_write_ops=0,kv_commit_lat.avgcount=291889,kv_commit_lat.avgtime=0.006347015,kv_commit_lat.sum=1852.624108527,kv_final_lat.avgcount=291885,kv_final_lat.avgtime=0.00004358,kv_final_lat.sum=12.720529751,kv_flush_lat.avgcount=291889,kv_flush_lat.avgtime=0.000000211,kv_flush_lat.sum=0.061636079,kv_sync_lat.avgcount=291889,kv_sync_lat.avgtime=0.006347227,kv_sync_lat.sum=1852.685744606,omap_lower_bound_lat.avgcount=1,omap_lower_bound_lat.avgtime=0.000004482,omap_lower_bound_lat.sum=0.000004482,omap_next_lat.avgcount=6933,omap_next_lat.avgtime=0.000003956,omap_next_lat.sum=0.027427456,omap_seek_to_first_lat.avgcount=309,omap_seek_to_first_lat.avgtime=0.000005879,omap_seek_to_first_lat.sum=0.001816658,omap_upper_bound_lat.avgcount=0,omap_upper_bound_lat.avgtime=0,omap_upper_bound_lat.sum=0,read_lat.avgcount=229,read_lat.avgtime=0.000394981,read_lat.sum=0.090450704,read_onode_meta_lat.avgcount=295,read_onode_meta_lat.avgtime=0.000016832,read_onode_meta_lat.sum=0.004965516,read_wait_aio_lat.avgcount=66,read_wait_aio_lat.avgtime=0.001237841,read_wait_aio_lat.sum=0.081697561,state_aio_wait_lat.avgcount=295997,state_aio_wait_lat.avgtime=0.000000357,state_aio_wait_lat.sum=0.105827433,state_deferred_aio_wait_lat.avgcount=0,state_deferred_aio_wait_lat.avgtime=0,state_deferred_aio_wait_lat.sum=0,state_deferred_cleanup_lat.avgcount=0,state_deferred_cleanup_lat.avgtime=0,state_deferred_cleanup_lat.sum=0,state_deferred_queued_lat.avgcount=0,state_deferred_queued_lat.avgtime=0,state_deferred_queued_lat.sum=0,state_done_lat.avgcount=295986,state_done_lat.avgtime=0.000003017,state_done_lat.sum=0.893199127,state_finishing_lat.avgcount=295986,state_finishing_lat.avgtime=0.000000306,state_finishing_lat.sum=0.090792683,state_io_done_lat.avgcount=295997,state_io_done_lat.avgtime=0.000001066,state_io_done_lat.sum=0.315577655,state_kv_commiting_lat.avgcount=295997,state_kv_commiting_lat.avgtime=0.006423586,state_kv_commiting_lat.sum=1901.362268572,state_kv_done_lat.avgcount=295997,state_kv_done_lat.avgtime=0.00000155,state_kv_done_lat.sum=0.458963064,state_kv_queued_lat.avgcount=295997,state_kv_queued_lat.avgtime=0.000477234,state_kv_queued_lat.sum=141.260101773,state_prepare_lat.avgcount=295997,state_prepare_lat.avgtime=0.000091806,state_prepare_lat.sum=27.174436583,submit_lat.avgcount=295997,submit_lat.avgtime=0.000135729,submit_lat.sum=40.17557682,throttle_lat.avgcount=295997,throttle_lat.avgtime=0.000002734,throttle_lat.sum=0.809479837,write_pad_bytes=151773,write_penalty_read_ops=0 1587117698000000000
ceph,collection=throttle-bluestore_throttle_bytes,host=stefanosd1,id=2,type=osd get=295997,get_or_fail_fail=0,get_or_fail_success=0,get_started=295997,get_sum=198686579299,max=67108864,put=291889,put_sum=198686579299,take=0,take_sum=0,val=0,wait.avgcount=83,wait.avgtime=0.003670612,wait.sum=0.304660858 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-cluster,host=stefanosd1,id=2,type=osd get=452060,get_or_fail_fail=0,get_or_fail_success=452060,get_started=0,get_sum=269934345,max=104857600,put=452060,put_sum=269934345,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-bluestore_throttle_deferred_bytes,host=stefanosd1,id=2,type=osd get=11,get_or_fail_fail=0,get_or_fail_success=11,get_started=0,get_sum=7723117,max=201326592,put=0,put_sum=0,take=0,take_sum=0,val=7723117,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_front_server,host=stefanosd1,id=2,type=osd get=2607433,get_or_fail_fail=0,get_or_fail_success=2607433,get_started=0,get_sum=5225295732,max=104857600,put=2607433,put_sum=5225295732,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=finisher-objecter-finisher-0,host=stefanosd1,id=2,type=osd complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117698000000000
ceph,collection=cct,host=stefanosd1,id=2,type=osd total_workers=6,unhealthy_workers=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-2,host=stefanosd1,id=2,type=osd msgr_active_connections=670,msgr_created_connections=13455,msgr_recv_bytes=6334605563,msgr_recv_messages=3287843,msgr_running_fast_dispatch_time=137.016615819,msgr_running_recv_time=240.687997039,msgr_running_send_time=471.710658466,msgr_running_total_time=1034.029109337,msgr_send_bytes=9753423475,msgr_send_messages=3439611 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-client,host=stefanosd1,id=2,type=osd get=710355,get_or_fail_fail=0,get_or_fail_success=710355,get_started=0,get_sum=166306283,max=104857600,put=710355,put_sum=166306283,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=throttle-msgr_dispatch_throttler-hb_back_server,host=stefanosd1,id=2,type=osd get=2607433,get_or_fail_fail=0,get_or_fail_success=2607433,get_started=0,get_sum=5225295732,max=104857600,put=2607433,put_sum=5225295732,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=AsyncMessenger::Worker-0,host=stefanosd1,id=2,type=osd msgr_active_connections=705,msgr_created_connections=17953,msgr_recv_bytes=7261438733,msgr_recv_messages=4496034,msgr_running_fast_dispatch_time=254.716476808,msgr_running_recv_time=272.196741555,msgr_running_send_time=571.102924903,msgr_running_total_time=1338.461077493,msgr_send_bytes=10772250508,msgr_send_messages=4192781 1587117698000000000
ceph,collection=rocksdb,host=stefanosd1,id=2,type=osd compact=0,compact_queue_len=0,compact_queue_merge=0,compact_range=0,get=1424,get_latency.avgcount=1424,get_latency.avgtime=0.000030752,get_latency.sum=0.043792142,rocksdb_write_delay_time.avgcount=0,rocksdb_write_delay_time.avgtime=0,rocksdb_write_delay_time.sum=0,rocksdb_write_memtable_time.avgcount=0,rocksdb_write_memtable_time.avgtime=0,rocksdb_write_memtable_time.sum=0,rocksdb_write_pre_and_post_time.avgcount=0,rocksdb_write_pre_and_post_time.avgtime=0,rocksdb_write_pre_and_post_time.sum=0,rocksdb_write_wal_time.avgcount=0,rocksdb_write_wal_time.avgtime=0,rocksdb_write_wal_time.sum=0,submit_latency.avgcount=295997,submit_latency.avgtime=0.000173137,submit_latency.sum=51.248072285,submit_sync_latency.avgcount=291889,submit_sync_latency.avgtime=0.006094397,submit_sync_latency.sum=1778.887521449,submit_transaction=295997,submit_transaction_sync=291889 1587117698000000000
ceph,collection=throttle-osd_client_bytes,host=stefanosd1,id=2,type=osd get=698701,get_or_fail_fail=0,get_or_fail_success=698701,get_started=0,get_sum=165630172,max=524288000,put=920880,put_sum=165630172,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117698000000000
ceph,collection=mds_sessions,host=stefanmds1,id=stefanmds1,type=mds average_load=0,avg_session_uptime=0,session_add=0,session_count=0,session_remove=0,sessions_open=0,sessions_stale=0,total_load=0 1587117476000000000
ceph,collection=mempool,host=stefanmds1,id=stefanmds1,type=mds bloom_filter_bytes=0,bloom_filter_items=0,bluefs_bytes=0,bluefs_items=0,bluestore_alloc_bytes=0,bluestore_alloc_items=0,bluestore_cache_data_bytes=0,bluestore_cache_data_items=0,bluestore_cache_onode_bytes=0,bluestore_cache_onode_items=0,bluestore_cache_other_bytes=0,bluestore_cache_other_items=0,bluestore_fsck_bytes=0,bluestore_fsck_items=0,bluestore_txc_bytes=0,bluestore_txc_items=0,bluestore_writing_bytes=0,bluestore_writing_deferred_bytes=0,bluestore_writing_deferred_items=0,bluestore_writing_items=0,buffer_anon_bytes=132069,buffer_anon_items=82,buffer_meta_bytes=0,buffer_meta_items=0,mds_co_bytes=44208,mds_co_items=154,osd_bytes=0,osd_items=0,osd_mapbl_bytes=0,osd_mapbl_items=0,osd_pglog_bytes=0,osd_pglog_items=0,osdmap_bytes=16952,osdmap_items=139,osdmap_mapping_bytes=0,osdmap_mapping_items=0,pgmap_bytes=0,pgmap_items=0,unittest_1_bytes=0,unittest_1_items=0,unittest_2_bytes=0,unittest_2_items=0 1587117476000000000
ceph,collection=objecter,host=stefanmds1,id=stefanmds1,type=mds command_active=0,command_resend=0,command_send=0,linger_active=0,linger_ping=0,linger_resend=0,linger_send=0,map_epoch=203,map_full=0,map_inc=1,omap_del=0,omap_rd=28,omap_wr=1,op=33,op_active=0,op_laggy=0,op_pg=0,op_r=26,op_reply=33,op_resend=2,op_rmw=0,op_send=35,op_send_bytes=364,op_w=7,osd_laggy=0,osd_session_close=91462,osd_session_open=91468,osd_sessions=6,osdop_append=0,osdop_call=0,osdop_clonerange=0,osdop_cmpxattr=0,osdop_create=0,osdop_delete=5,osdop_getxattr=14,osdop_mapext=0,osdop_notify=0,osdop_other=0,osdop_pgls=0,osdop_pgls_filter=0,osdop_read=8,osdop_resetxattrs=0,osdop_rmxattr=0,osdop_setxattr=0,osdop_sparse_read=0,osdop_src_cmpxattr=0,osdop_stat=2,osdop_truncate=0,osdop_watch=0,osdop_write=0,osdop_writefull=0,osdop_writesame=0,osdop_zero=1,poolop_active=0,poolop_resend=0,poolop_send=0,poolstat_active=0,poolstat_resend=0,poolstat_send=0,statfs_active=0,statfs_resend=0,statfs_send=0 1587117476000000000
ceph,collection=cct,host=stefanmds1,id=stefanmds1,type=mds total_workers=1,unhealthy_workers=0 1587117476000000000
ceph,collection=mds_server,host=stefanmds1,id=stefanmds1,type=mds cap_revoke_eviction=0,dispatch_client_request=0,dispatch_server_request=0,handle_client_request=0,handle_client_session=0,handle_slave_request=0,req_create_latency.avgcount=0,req_create_latency.avgtime=0,req_create_latency.sum=0,req_getattr_latency.avgcount=0,req_getattr_latency.avgtime=0,req_getattr_latency.sum=0,req_getfilelock_latency.avgcount=0,req_getfilelock_latency.avgtime=0,req_getfilelock_latency.sum=0,req_link_latency.avgcount=0,req_link_latency.avgtime=0,req_link_latency.sum=0,req_lookup_latency.avgcount=0,req_lookup_latency.avgtime=0,req_lookup_latency.sum=0,req_lookuphash_latency.avgcount=0,req_lookuphash_latency.avgtime=0,req_lookuphash_latency.sum=0,req_lookupino_latency.avgcount=0,req_lookupino_latency.avgtime=0,req_lookupino_latency.sum=0,req_lookupname_latency.avgcount=0,req_lookupname_latency.avgtime=0,req_lookupname_latency.sum=0,req_lookupparent_latency.avgcount=0,req_lookupparent_latency.avgtime=0,req_lookupparent_latency.sum=0,req_lookupsnap_latency.avgcount=0,req_lookupsnap_latency.avgtime=0,req_lookupsnap_latency.sum=0,req_lssnap_latency.avgcount=0,req_lssnap_latency.avgtime=0,req_lssnap_latency.sum=0,req_mkdir_latency.avgcount=0,req_mkdir_latency.avgtime=0,req_mkdir_latency.sum=0,req_mknod_latency.avgcount=0,req_mknod_latency.avgtime=0,req_mknod_latency.sum=0,req_mksnap_latency.avgcount=0,req_mksnap_latency.avgtime=0,req_mksnap_latency.sum=0,req_open_latency.avgcount=0,req_open_latency.avgtime=0,req_open_latency.sum=0,req_readdir_latency.avgcount=0,req_readdir_latency.avgtime=0,req_readdir_latency.sum=0,req_rename_latency.avgcount=0,req_rename_latency.avgtime=0,req_rename_latency.sum=0,req_renamesnap_latency.avgcount=0,req_renamesnap_latency.avgtime=0,req_renamesnap_latency.sum=0,req_rmdir_latency.avgcount=0,req_rmdir_latency.avgtime=0,req_rmdir_latency.sum=0,req_rmsnap_latency.avgcount=0,req_rmsnap_latency.avgtime=0,req_rmsnap_latency.sum=0,req_rmxattr_latency.avgcount=0,req_rmxattr_latency.avgtime=0,req_rmxattr_latency.sum=0,req_setattr_latency.avgcount=0,req_setattr_latency.avgtime=0,req_setattr_latency.sum=0,req_setdirlayout_latency.avgcount=0,req_setdirlayout_latency.avgtime=0,req_setdirlayout_latency.sum=0,req_setfilelock_latency.avgcount=0,req_setfilelock_latency.avgtime=0,req_setfilelock_latency.sum=0,req_setlayout_latency.avgcount=0,req_setlayout_latency.avgtime=0,req_setlayout_latency.sum=0,req_setxattr_latency.avgcount=0,req_setxattr_latency.avgtime=0,req_setxattr_latency.sum=0,req_symlink_latency.avgcount=0,req_symlink_latency.avgtime=0,req_symlink_latency.sum=0,req_unlink_latency.avgcount=0,req_unlink_latency.avgtime=0,req_unlink_latency.sum=0 1587117476000000000
ceph,collection=AsyncMessenger::Worker-2,host=stefanmds1,id=stefanmds1,type=mds msgr_active_connections=84,msgr_created_connections=68511,msgr_recv_bytes=238078,msgr_recv_messages=2655,msgr_running_fast_dispatch_time=0.004247777,msgr_running_recv_time=25.369012545,msgr_running_send_time=3.743427461,msgr_running_total_time=130.277111559,msgr_send_bytes=172767043,msgr_send_messages=18172 1587117476000000000
ceph,collection=mds_log,host=stefanmds1,id=stefanmds1,type=mds ev=0,evadd=0,evex=0,evexd=0,evexg=0,evtrm=0,expos=4194304,jlat.avgcount=0,jlat.avgtime=0,jlat.sum=0,rdpos=4194304,replayed=1,seg=1,segadd=0,segex=0,segexd=0,segexg=0,segtrm=0,wrpos=0 1587117476000000000
ceph,collection=AsyncMessenger::Worker-0,host=stefanmds1,id=stefanmds1,type=mds msgr_active_connections=595,msgr_created_connections=943825,msgr_recv_bytes=78618003,msgr_recv_messages=914080,msgr_running_fast_dispatch_time=0.001544386,msgr_running_recv_time=459.627068807,msgr_running_send_time=469.337032316,msgr_running_total_time=2744.084305898,msgr_send_bytes=61684163658,msgr_send_messages=1858008 1587117476000000000
ceph,collection=throttle-msgr_dispatch_throttler-mds,host=stefanmds1,id=stefanmds1,type=mds get=1216458,get_or_fail_fail=0,get_or_fail_success=1216458,get_started=0,get_sum=51976882,max=104857600,put=1216458,put_sum=51976882,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117476000000000
ceph,collection=AsyncMessenger::Worker-1,host=stefanmds1,id=stefanmds1,type=mds msgr_active_connections=226,msgr_created_connections=42679,msgr_recv_bytes=63140151,msgr_recv_messages=299727,msgr_running_fast_dispatch_time=26.316138629,msgr_running_recv_time=36.969916165,msgr_running_send_time=70.457421128,msgr_running_total_time=226.230019936,msgr_send_bytes=193154464,msgr_send_messages=310481 1587117476000000000
ceph,collection=mds,host=stefanmds1,id=stefanmds1,type=mds caps=0,dir_commit=0,dir_fetch=12,dir_merge=0,dir_split=0,exported=0,exported_inodes=0,forward=0,imported=0,imported_inodes=0,inode_max=2147483647,inodes=10,inodes_bottom=3,inodes_expired=0,inodes_pin_tail=0,inodes_pinned=10,inodes_top=7,inodes_with_caps=0,load_cent=0,openino_backtrace_fetch=0,openino_dir_fetch=0,openino_peer_discover=0,q=0,reply=0,reply_latency.avgcount=0,reply_latency.avgtime=0,reply_latency.sum=0,request=0,subtrees=2,traverse=0,traverse_dir_fetch=0,traverse_discover=0,traverse_forward=0,traverse_hit=0,traverse_lock=0,traverse_remote_ino=0 1587117476000000000
ceph,collection=purge_queue,host=stefanmds1,id=stefanmds1,type=mds pq_executed=0,pq_executing=0,pq_executing_ops=0 1587117476000000000
ceph,collection=throttle-write_buf_throttle,host=stefanmds1,id=stefanmds1,type=mds get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=3758096384,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117476000000000
ceph,collection=throttle-write_buf_throttle-0x5624e9377f40,host=stefanmds1,id=stefanmds1,type=mds get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=3758096384,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117476000000000
ceph,collection=mds_cache,host=stefanmds1,id=stefanmds1,type=mds ireq_enqueue_scrub=0,ireq_exportdir=0,ireq_flush=0,ireq_fragmentdir=0,ireq_fragstats=0,ireq_inodestats=0,num_recovering_enqueued=0,num_recovering_prioritized=0,num_recovering_processing=0,num_strays=0,num_strays_delayed=0,num_strays_enqueuing=0,recovery_completed=0,recovery_started=0,strays_created=0,strays_enqueued=0,strays_migrated=0,strays_reintegrated=0 1587117476000000000
ceph,collection=throttle-objecter_bytes,host=stefanmds1,id=stefanmds1,type=mds get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=104857600,put=16,put_sum=1016,take=33,take_sum=1016,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117476000000000
ceph,collection=throttle-objecter_ops,host=stefanmds1,id=stefanmds1,type=mds get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=1024,put=33,put_sum=33,take=33,take_sum=33,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117476000000000
ceph,collection=mds_mem,host=stefanmds1,id=stefanmds1,type=mds cap=0,cap+=0,cap-=0,dir=12,dir+=12,dir-=0,dn=10,dn+=10,dn-=0,heap=322284,ino=13,ino+=13,ino-=0,rss=76032 1587117476000000000
ceph,collection=finisher-PurgeQueue,host=stefanmds1,id=stefanmds1,type=mds complete_latency.avgcount=4,complete_latency.avgtime=0.000176985,complete_latency.sum=0.000707941,queue_len=0 1587117476000000000
ceph,collection=cct,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw total_workers=0,unhealthy_workers=0 1587117156000000000
ceph,collection=throttle-objecter_bytes,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw get=791732,get_or_fail_fail=0,get_or_fail_success=791732,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117156000000000
ceph,collection=rgw,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw cache_hit=0,cache_miss=791706,failed_req=0,get=0,get_b=0,get_initial_lat.avgcount=0,get_initial_lat.avgtime=0,get_initial_lat.sum=0,keystone_token_cache_hit=0,keystone_token_cache_miss=0,pubsub_event_lost=0,pubsub_event_triggered=0,pubsub_events=0,pubsub_push_failed=0,pubsub_push_ok=0,pubsub_push_pending=0,pubsub_store_fail=0,pubsub_store_ok=0,put=0,put_b=0,put_initial_lat.avgcount=0,put_initial_lat.avgtime=0,put_initial_lat.sum=0,qactive=0,qlen=0,req=791705 1587117156000000000
ceph,collection=throttle-msgr_dispatch_throttler-radosclient,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw get=2697988,get_or_fail_fail=0,get_or_fail_success=2697988,get_started=0,get_sum=444563051,max=104857600,put=2697988,put_sum=444563051,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117156000000000
ceph,collection=finisher-radosclient,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw complete_latency.avgcount=2,complete_latency.avgtime=0.003530161,complete_latency.sum=0.007060323,queue_len=0 1587117156000000000
ceph,collection=throttle-rgw_async_rados_ops,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw get=0,get_or_fail_fail=0,get_or_fail_success=0,get_started=0,get_sum=0,max=64,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117156000000000
ceph,collection=throttle-objecter_ops,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw get=791732,get_or_fail_fail=0,get_or_fail_success=791732,get_started=0,get_sum=791732,max=24576,put=791732,put_sum=791732,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117156000000000
ceph,collection=throttle-objecter_bytes-0x5598969981c0,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw get=1637900,get_or_fail_fail=0,get_or_fail_success=1637900,get_started=0,get_sum=0,max=104857600,put=0,put_sum=0,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117156000000000
ceph,collection=objecter,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw command_active=0,command_resend=0,command_send=0,linger_active=8,linger_ping=1905736,linger_resend=4,linger_send=13,map_epoch=203,map_full=0,map_inc=17,omap_del=0,omap_rd=0,omap_wr=0,op=2697488,op_active=0,op_laggy=0,op_pg=0,op_r=791730,op_reply=2697476,op_resend=1,op_rmw=0,op_send=2697490,op_send_bytes=362,op_w=1905758,osd_laggy=5,osd_session_close=59558,osd_session_open=59566,osd_sessions=8,osdop_append=0,osdop_call=1,osdop_clonerange=0,osdop_cmpxattr=0,osdop_create=8,osdop_delete=0,osdop_getxattr=0,osdop_mapext=0,osdop_notify=0,osdop_other=791714,osdop_pgls=0,osdop_pgls_filter=0,osdop_read=16,osdop_resetxattrs=0,osdop_rmxattr=0,osdop_setxattr=0,osdop_sparse_read=0,osdop_src_cmpxattr=0,osdop_stat=791706,osdop_truncate=0,osdop_watch=1905750,osdop_write=0,osdop_writefull=0,osdop_writesame=0,osdop_zero=0,poolop_active=0,poolop_resend=0,poolop_send=0,poolstat_active=0,poolstat_resend=0,poolstat_send=0,statfs_active=0,statfs_resend=0,statfs_send=0 1587117156000000000
ceph,collection=AsyncMessenger::Worker-2,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw msgr_active_connections=11,msgr_created_connections=59839,msgr_recv_bytes=342697143,msgr_recv_messages=1441603,msgr_running_fast_dispatch_time=161.807937536,msgr_running_recv_time=118.174064257,msgr_running_send_time=207.679154333,msgr_running_total_time=698.527662129,msgr_send_bytes=530785909,msgr_send_messages=1679950 1587117156000000000
ceph,collection=mempool,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw bloom_filter_bytes=0,bloom_filter_items=0,bluefs_bytes=0,bluefs_items=0,bluestore_alloc_bytes=0,bluestore_alloc_items=0,bluestore_cache_data_bytes=0,bluestore_cache_data_items=0,bluestore_cache_onode_bytes=0,bluestore_cache_onode_items=0,bluestore_cache_other_bytes=0,bluestore_cache_other_items=0,bluestore_fsck_bytes=0,bluestore_fsck_items=0,bluestore_txc_bytes=0,bluestore_txc_items=0,bluestore_writing_bytes=0,bluestore_writing_deferred_bytes=0,bluestore_writing_deferred_items=0,bluestore_writing_items=0,buffer_anon_bytes=225471,buffer_anon_items=163,buffer_meta_bytes=0,buffer_meta_items=0,mds_co_bytes=0,mds_co_items=0,osd_bytes=0,osd_items=0,osd_mapbl_bytes=0,osd_mapbl_items=0,osd_pglog_bytes=0,osd_pglog_items=0,osdmap_bytes=33904,osdmap_items=278,osdmap_mapping_bytes=0,osdmap_mapping_items=0,pgmap_bytes=0,pgmap_items=0,unittest_1_bytes=0,unittest_1_items=0,unittest_2_bytes=0,unittest_2_items=0 1587117156000000000
ceph,collection=throttle-msgr_dispatch_throttler-radosclient-0x559896998120,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw get=1652935,get_or_fail_fail=0,get_or_fail_success=1652935,get_started=0,get_sum=276333029,max=104857600,put=1652935,put_sum=276333029,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117156000000000
ceph,collection=AsyncMessenger::Worker-1,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw msgr_active_connections=17,msgr_created_connections=84859,msgr_recv_bytes=211170759,msgr_recv_messages=922646,msgr_running_fast_dispatch_time=31.487443762,msgr_running_recv_time=83.190789333,msgr_running_send_time=174.670510496,msgr_running_total_time=484.22086275,msgr_send_bytes=1322113179,msgr_send_messages=1636839 1587117156000000000
ceph,collection=finisher-radosclient-0x559896998080,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw complete_latency.avgcount=0,complete_latency.avgtime=0,complete_latency.sum=0,queue_len=0 1587117156000000000
ceph,collection=throttle-objecter_ops-0x559896997b80,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw get=1637900,get_or_fail_fail=0,get_or_fail_success=1637900,get_started=0,get_sum=1637900,max=24576,put=1637900,put_sum=1637900,take=0,take_sum=0,val=0,wait.avgcount=0,wait.avgtime=0,wait.sum=0 1587117156000000000
ceph,collection=AsyncMessenger::Worker-0,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw msgr_active_connections=18,msgr_created_connections=74757,msgr_recv_bytes=489001094,msgr_recv_messages=1986686,msgr_running_fast_dispatch_time=168.60950961,msgr_running_recv_time=142.903031533,msgr_running_send_time=267.911165712,msgr_running_total_time=824.885614951,msgr_send_bytes=707973504,msgr_send_messages=2463727 1587117156000000000
ceph,collection=objecter-0x559896997720,host=stefanrgw1,id=rgw.stefanrgw1.4219.94113851143184,type=rgw command_active=0,command_resend=0,command_send=0,linger_active=0,linger_ping=0,linger_resend=0,linger_send=0,map_epoch=203,map_full=0,map_inc=8,omap_del=0,omap_rd=0,omap_wr=0,op=1637998,op_active=0,op_laggy=0,op_pg=0,op_r=1062803,op_reply=1637998,op_resend=15,op_rmw=0,op_send=1638013,op_send_bytes=63321099,op_w=575195,osd_laggy=0,osd_session_close=125555,osd_session_open=125563,osd_sessions=8,osdop_append=0,osdop_call=1637886,osdop_clonerange=0,osdop_cmpxattr=0,osdop_create=0,osdop_delete=0,osdop_getxattr=0,osdop_mapext=0,osdop_notify=0,osdop_other=112,osdop_pgls=0,osdop_pgls_filter=0,osdop_read=0,osdop_resetxattrs=0,osdop_rmxattr=0,osdop_setxattr=0,osdop_sparse_read=0,osdop_src_cmpxattr=0,osdop_stat=0,osdop_truncate=0,osdop_watch=0,osdop_write=0,osdop_writefull=0,osdop_writesame=0,osdop_zero=0,poolop_active=0,poolop_resend=0,poolop_send=0,poolstat_active=0,poolstat_resend=0,poolstat_send=0,statfs_active=0,statfs_resend=0,statfs_send=0 1587117156000000000
```

*Cluster Stats*

* ceph\_pgmap\_state has the following tags:
  * state (state for which the value applies e.g. active+clean, active+remapped+backfill)
* ceph\_pool\_usage has the following tags:
  * id
  * name
* ceph\_pool\_stats has the following tags:
  * id
  * name

### Example Output:

*Admin Socket Stats*

<pre>
telegraf --config /etc/telegraf/telegraf.conf --config-directory /etc/telegraf/telegraf.d --input-filter ceph --test
* Plugin: ceph, Collection 1
> ceph,collection=paxos, id=node-2,role=openstack,type=mon accept_timeout=0,begin=14931264,begin_bytes.avgcount=14931264,begin_bytes.sum=180309683362,begin_keys.avgcount=0,begin_keys.sum=0,begin_latency.avgcount=14931264,begin_latency.sum=9293.29589,collect=1,collect_bytes.avgcount=1,collect_bytes.sum=24,collect_keys.avgcount=1,collect_keys.sum=1,collect_latency.avgcount=1,collect_latency.sum=0.00028,collect_timeout=0,collect_uncommitted=0,commit=14931264,commit_bytes.avgcount=0,commit_bytes.sum=0,commit_keys.avgcount=0,commit_keys.sum=0,commit_latency.avgcount=0,commit_latency.sum=0,lease_ack_timeout=0,lease_timeout=0,new_pn=0,new_pn_latency.avgcount=0,new_pn_latency.sum=0,refresh=14931264,refresh_latency.avgcount=14931264,refresh_latency.sum=8706.98498,restart=4,share_state=0,share_state_bytes.avgcount=0,share_state_bytes.sum=0,share_state_keys.avgcount=0,share_state_keys.sum=0,start_leader=0,start_peon=1,store_state=14931264,store_state_bytes.avgcount=14931264,store_state_bytes.sum=353119959211,store_state_keys.avgcount=14931264,store_state_keys.sum=289807523,store_state_latency.avgcount=14931264,store_state_latency.sum=10952.835724 1462821234814535148
> ceph,collection=throttle-mon_client_bytes,id=node-2,type=mon get=1413017,get_or_fail_fail=0,get_or_fail_success=0,get_sum=71211705,max=104857600,put=1413013,put_sum=71211459,take=0,take_sum=0,val=246,wait.avgcount=0,wait.sum=0 1462821234814737219
> ceph,collection=throttle-mon_daemon_bytes,id=node-2,type=mon get=4058121,get_or_fail_fail=0,get_or_fail_success=0,get_sum=6027348117,max=419430400,put=4058121,put_sum=6027348117,take=0,take_sum=0,val=0,wait.avgcount=0,wait.sum=0 1462821234814815661
> ceph,collection=throttle-msgr_dispatch_throttler-mon,id=node-2,type=mon get=54276277,get_or_fail_fail=0,get_or_fail_success=0,get_sum=370232877040,max=104857600,put=54276277,put_sum=370232877040,take=0,take_sum=0,val=0,wait.avgcount=0,wait.sum=0 1462821234814872064
</pre>

*Cluster Stats*

<pre>
> ceph_osdmap,host=ceph-mon-0 epoch=170772,full=false,nearfull=false,num_in_osds=340,num_osds=340,num_remapped_pgs=0,num_up_osds=340 1468841037000000000
> ceph_pgmap,host=ceph-mon-0 bytes_avail=634895531270144,bytes_total=812117151809536,bytes_used=177221620539392,data_bytes=56979991615058,num_pgs=22952,op_per_sec=15869,read_bytes_sec=43956026,version=39387592,write_bytes_sec=165344818 1468841037000000000
> ceph_pgmap_state,host=ceph-mon-0,state=active+clean count=22952 1468928660000000000
> ceph_pgmap_state,host=ceph-mon-0,state=active+degraded count=16 1468928660000000000
> ceph_usage,host=ceph-mon-0 total_avail_bytes=634895514791936,total_bytes=812117151809536,total_used_bytes=177221637017600 1468841037000000000
> ceph_pool_usage,host=ceph-mon-0,id=150,name=cinder.volumes bytes_used=12648553794802,kb_used=12352103316,max_avail=154342562489244,objects=3026295 1468841037000000000
> ceph_pool_usage,host=ceph-mon-0,id=182,name=cinder.volumes.flash bytes_used=8541308223964,kb_used=8341121313,max_avail=39388593563936,objects=2075066 1468841037000000000
> ceph_pool_stats,host=ceph-mon-0,id=150,name=cinder.volumes op_per_sec=1706,read_bytes_sec=28671674,write_bytes_sec=29994541 1468841037000000000
> ceph_pool_stats,host=ceph-mon-0,id=182,name=cinder.volumes.flash op_per_sec=9748,read_bytes_sec=9605524,write_bytes_sec=45593310 1468841037000000000
</pre>
