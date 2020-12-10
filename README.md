# domaininfo

This project makes use of the https://team-cymru.com/community-services/ip-asn-mapping/ service.

## Use

A command line tool for looking up details about domains and IP addresses.

```sh
domaininfo git:main ❯ make build
go build -o ./bin/domaininfo ./cmd/domaininfo
domaininfo git:main ❯ ./bin/domaininfo www.cnn.com
Canonical Name (0): turner-tls.map.fastly.net.
---
IPv4 Address: 151.101.193.67
ASN Info: [54113 | 151.101.192.0/22 | US | arin | 2016-02-01]
IPv4 Address: 151.101.1.67
ASN Info: [54113 | 151.101.0.0/22 | US | arin | 2016-02-01]
IPv4 Address: 151.101.65.67
ASN Info: [54113 | 151.101.64.0/22 | US | arin | 2016-02-01]
IPv4 Address: 151.101.129.67
ASN Info: [54113 | 151.101.128.0/22 | US | arin | 2016-02-01]
---
IPv6 Address: 2a04:4e42::323
ASN Info: [54113 | 2a04:4e42::/36 | EU | ripencc | 2013-07-18]
IPv6 Address: 2a04:4e42:400::323
ASN Info: [54113 | 2a04:4e42::/36 | EU | ripencc | 2013-07-18]
IPv6 Address: 2a04:4e42:600::323
ASN Info: [54113 | 2a04:4e42::/36 | EU | ripencc | 2013-07-18]
IPv6 Address: 2a04:4e42:200::323
ASN Info: [54113 | 2a04:4e42::/36 | EU | ripencc | 2013-07-18]
---
ASN Description: [54113 | US | arin | 2011-10-04 | FASTLY, US]
---
No CAA record on www.cnn.com
No CAA record on turner-tls.map.fastly.net.
No CAA record on cnn.com
```

The command line output provides:

 * The canonical name which is the final name after following zero or more CNAME records
 * IPv4 and IPv6 addresses from DNS lookup
 * Autonomous system number (ASN) info for an IP address
 * A description of all autonomous system numbers found for the IP addresses
 * CAA record lookup according to https://docs.digicert.com/manage-certificates/dns-caa-resource-record-check/

 # Further Reading

  * https://en.wikipedia.org/wiki/Autonomous_system_(Internet)
  * https://team-cymru.com/community-services/
  * https://github.com/miekg/dns
