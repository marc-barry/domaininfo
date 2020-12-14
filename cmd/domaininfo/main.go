package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/marc-barry/domaininfo/pkg/dnsutil"
	"github.com/marc-barry/domaininfo/pkg/types"
)

var resolver = dnsutil.NewResolver()

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Requires at least one command line argument")
	}

	domain := os.Args[1]

	targets, err := dnsutil.CNAMEChain(resolver, domain)
	if err != nil {
		log.Fatal(err)
	}

	ipv4s, err := dnsutil.IPv4List(resolver, domain)
	if err != nil {
		log.Fatal(err)
	}
	ipv6s, err := dnsutil.IPv6List(resolver, domain)
	if err != nil {
		log.Fatal(err)
	}

	ipv4Info, ipv6Info, asns, err := dnsutil.NewAddressesInfo(resolver, ipv4s, ipv6s)
	if err != nil {
		log.Fatal(err)
	}

	asnDescriptions := make(map[string]types.ASNDescription)

	for k := range asns {
		if res, err := resolver.LookupTXT(fmt.Sprintf(types.ASNLOOKUPTEMPLATE, k)); err == nil {
			for _, r := range res {
				for _, txt := range r.Txt {
					sp := strings.Split(txt, " | ")
					if len(sp) == 5 && sp[0] == k {
						asnDescriptions[k] = types.ASNDescription{Country: sp[1], InternetRegistry: sp[2], Date: sp[3], Org: sp[4]}
					}
				}
			}
		}
	}

	caaInfos := make([]types.CAAInfo, 0)
	found := false

	res, err := resolver.LookupCAA(domain)
	if err != nil {
		log.Printf("Error looking up CAA records: %s\n", err)
	}
	info := types.CAAInfo{Domain: domain, CAs: make([]string, 0)}
	for _, r := range res {
		found = true
		info.CAs = append(info.CAs, r.Value)
	}
	caaInfos = append(caaInfos, info)

	if !found {
		for i, domain := range targets {
			if i == 7 {
				break
			}
			res, err := resolver.LookupCAA(domain)
			if err != nil {
				log.Printf("Error looking up CAA records: %s\n", err)
				continue
			}
			info := types.CAAInfo{Domain: domain, CAs: make([]string, 0)}
			for _, r := range res {
				info.CAs = append(info.CAs, r.Value)
			}
			caaInfos = append(caaInfos, info)
			if len(res) != 0 {
				found = true
			}
			break
		}
		if !found {
			i := strings.IndexAny(domain, ".")
			if i > 0 {
				parent := domain[i+1:]
				res, err := resolver.LookupCAA(parent)
				if err != nil {
					log.Printf("Error looking up CAA records: %s\n", err)
				}
				info := types.CAAInfo{Domain: parent, CAs: make([]string, 0)}
				for _, r := range res {
					info.CAs = append(info.CAs, r.Value)
				}
				caaInfos = append(caaInfos, info)
			}
		}
	}

	dinfo := types.DomainInfo{
		Domain:                domain,
		CanonicalNamesTargets: targets,
		IPv4AddressInfo:       ipv4Info,
		IPv6AddressInfo:       ipv6Info,
		ASNDescriptions:       asnDescriptions,
		CAAInfos:              caaInfos,
	}

	if b, err := json.MarshalIndent(dinfo, "", "  "); err == nil {
		fmt.Println(string(b))
	}
}
