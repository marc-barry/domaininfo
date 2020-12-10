# domaininfo

This project makes use of the https://team-cymru.com/community-services/ip-asn-mapping/ service.

## Use

A command line tool for looking up details about domains and IP addresses.

```sh
domaininfo git:main ❯ make build
go build -o ./bin/domaininfo ./cmd/domaininfo
domaininfo git:main ❯ ./bin/domaininfo www.cnn.com
{
  "domain": "www.cnn.com",
  "canonicalNamesInfo": {
    "targets": [
      "turner-tls.map.fastly.net."
    ]
  },
  "addresses": {
    "ipv4AddressInfo": {
      "151.101.1.67": [
        {
          "asn": "54113",
          "addressBlock": "151.101.0.0/22",
          "country": "US",
          "internetRegistry": "arin",
          "date": "2016-02-01"
        }
      ],
      "151.101.129.67": [
        {
          "asn": "54113",
          "addressBlock": "151.101.128.0/22",
          "country": "US",
          "internetRegistry": "arin",
          "date": "2016-02-01"
        }
      ],
      "151.101.193.67": [
        {
          "asn": "54113",
          "addressBlock": "151.101.192.0/22",
          "country": "US",
          "internetRegistry": "arin",
          "date": "2016-02-01"
        }
      ],
      "151.101.65.67": [
        {
          "asn": "54113",
          "addressBlock": "151.101.64.0/22",
          "country": "US",
          "internetRegistry": "arin",
          "date": "2016-02-01"
        }
      ]
    },
    "ipv6AddressInfo": {
      "2a04:4e42:200::323": [
        {
          "asn": "54113",
          "addressBlock": "2a04:4e42::/36",
          "country": "EU",
          "internetRegistry": "ripencc",
          "date": "2013-07-18"
        }
      ],
      "2a04:4e42:400::323": [
        {
          "asn": "54113",
          "addressBlock": "2a04:4e42::/36",
          "country": "EU",
          "internetRegistry": "ripencc",
          "date": "2013-07-18"
        }
      ],
      "2a04:4e42:600::323": [
        {
          "asn": "54113",
          "addressBlock": "2a04:4e42::/36",
          "country": "EU",
          "internetRegistry": "ripencc",
          "date": "2013-07-18"
        }
      ],
      "2a04:4e42::323": [
        {
          "asn": "54113",
          "addressBlock": "2a04:4e42::/36",
          "country": "EU",
          "internetRegistry": "ripencc",
          "date": "2013-07-18"
        }
      ]
    }
  },
  "asns": {
    "descriptions": {
      "54113": {
        "country": "US",
        "internetRegistry": "arin",
        "date": "2011-10-04",
        "org": "FASTLY, US"
      }
    }
  },
  "caas": {
    "caaInfos": [
      {
        "domain": "www.cnn.com",
        "cas": []
      },
      {
        "domain": "turner-tls.map.fastly.net.",
        "cas": []
      },
      {
        "domain": "cnn.com",
        "cas": []
      }
    ]
  }
}
domaininfo git:main ❯ ./bin/domaininfo www.google.com
{
  "domain": "www.google.com",
  "canonicalNamesInfo": {
    "targets": []
  },
  "addresses": {
    "ipv4AddressInfo": {
      "172.217.15.100": [
        {
          "asn": "15169",
          "addressBlock": "172.217.15.0/24",
          "country": "US",
          "internetRegistry": "arin",
          "date": "2012-04-16"
        }
      ]
    },
    "ipv6AddressInfo": {
      "2607:f8b0:4004:801::2004": [
        {
          "asn": "15169",
          "addressBlock": "2607:f8b0:4004::/48",
          "country": "US",
          "internetRegistry": "arin",
          "date": "2009-03-12"
        }
      ]
    }
  },
  "asns": {
    "descriptions": {
      "15169": {
        "country": "US",
        "internetRegistry": "arin",
        "date": "2000-03-30",
        "org": "GOOGLE, US"
      }
    }
  },
  "caas": {
    "caaInfos": [
      {
        "domain": "www.google.com",
        "cas": []
      },
      {
        "domain": "google.com",
        "cas": [
          "google.com. 86052 IN CAA 0 issue \"pki.goog\""
        ]
      }
    ]
  }
}
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
