# x509 Certificate Input Plugin

This plugin provides information about X509 certificate accessible via local
file, tcp, udp, https or smtp protocol.

When using a UDP address as a certificate source, the server must support
[DTLS](https://en.wikipedia.org/wiki/Datagram_Transport_Layer_Security).

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md

## Configuration

```toml @sample.conf
# Reads metrics from a SSL certificate
[[inputs.x509_cert]]
  ## List certificate sources, support wildcard expands for files
  ## Prefix your entry with 'file://' if you intend to use relative paths
  sources = ["tcp://example.org:443", "https://influxdata.com:443",
            "smtp://mail.localhost:25", "udp://127.0.0.1:4433",
            "/etc/ssl/certs/ssl-cert-snakeoil.pem",
            "/etc/mycerts/*.mydomain.org.pem", "file:///path/to/*.pem"]

  ## Timeout for SSL connection
  # timeout = "5s"

  ## Pass a different name into the TLS request (Server Name Indication).
  ## This is synonymous with tls_server_name, and only one of the two
  ## options may be specified at one time.
  ##   example: server_name = "myhost.example.org"
  # server_name = "myhost.example.org"

  ## Only output the leaf certificates and omit the root ones.
  # exclude_root_certs = false

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  # tls_server_name = "myhost.example.org"

  ## Set the proxy URL
  # use_proxy = true
  # proxy_url = "http://localhost:8888"
```

## Metrics

  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
```


### Metrics

- x509_cert
  - tags:
    - source - source of the certificate
    - organization
    - organizational_unit
    - country
    - province
    - locality
    - verification
    - serial_number
    - signature_algorithm
    - public_key_algorithm
    - issuer_common_name
    - issuer_serial_number
    - san
  - fields:
    - verification_code (int)
    - verification_error (string)
  - fields:
    - expiry (int, seconds)
    - age (int, seconds)
    - startdate (int, seconds)
    - enddate (int, seconds)

## Example Output

```shell
x509_cert,common_name=ubuntu,source=/etc/ssl/certs/ssl-cert-snakeoil.pem,verification=valid age=7693222i,enddate=1871249033i,expiry=307666777i,startdate=1555889033i,verification_code=0i 1563582256000000000
x509_cert,common_name=www.example.org,country=US,locality=Los\ Angeles,organization=Internet\ Corporation\ for\ Assigned\ Names\ and\ Numbers,organizational_unit=Technology,province=California,source=https://example.org:443,verification=invalid age=20219055i,enddate=1606910400i,expiry=43328144i,startdate=1543363200i,verification_code=1i,verification_error="x509: certificate signed by unknown authority" 1563582256000000000
x509_cert,common_name=DigiCert\ SHA2\ Secure\ Server\ CA,country=US,organization=DigiCert\ Inc,source=https://example.org:443,verification=valid age=200838255i,enddate=1678276800i,expiry=114694544i,startdate=1362744000i,verification_code=0i 1563582256000000000
x509_cert,common_name=DigiCert\ Global\ Root\ CA,country=US,organization=DigiCert\ Inc,organizational_unit=www.digicert.com,source=https://example.org:443,verification=valid age=400465455i,enddate=1952035200i,expiry=388452944i,startdate=1163116800i,verification_code=0i 1563582256000000000

### Example output

```
x509_cert,host=myhost,source=https://example.org age=1753627i,expiry=5503972i,startdate=1516092060i,enddate=1523349660i 1517845687000000000
x509_cert,host=myhost,source=/etc/ssl/certs/ssl-cert-snakeoil.pem age=7522207i,expiry=308002732i,startdate=1510323480i,enddate=1825848420i 1517845687000000000
```
