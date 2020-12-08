# domaininfo

This project makes use of the https://team-cymru.com/community-services/ip-asn-mapping/ service.

## Use

A command line tool for looking up details about domains and IP addresses.

```sh
domaininfo git:main ❯ ./domaininfo google.com
Canonical Name: google.com.
---
IPv4 Address: 172.217.5.238
ASN Info: [15169 | 172.217.5.0/24 | US | arin | 2012-04-16]
---
IPv6 Address: 2607:f8b0:4004:804::200e
ASN Info: [15169 | 2607:f8b0:4004::/48 | US | arin | 2009-03-12]
---
ASN Description: [15169 | US | arin | 2000-03-30 | GOOGLE, US]
```

The command line output provides:

 * The canonical name which is the final name after following zero or more CNAME records
 * IPv4 and IPv6 addresses from DNS lookup
 * Autonomous system number (ASN) info for an IP address
 * A description of all autonomous system numbers found for the IP addresses

 # Further Reading

  * https://en.wikipedia.org/wiki/Autonomous_system_(Internet)
  * https://team-cymru.com/community-services/