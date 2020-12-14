package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	ipv4Info, ipv6Info, asns, err := dnsutil.AddressesInfos(resolver, ipv4s, ipv6s)
	if err != nil {
		log.Fatal(err)
	}

	if b, err := json.MarshalIndent(
		types.DomainInfo{
			Domain:                domain,
			CanonicalNamesTargets: targets,
			IPv4AddressInfo:       ipv4Info,
			IPv6AddressInfo:       ipv6Info,
			ASNDescriptions:       dnsutil.ASNDescriptions(resolver, asns),
			CAAInfos:              dnsutil.CAAInfos(resolver, domain, targets),
		}, "", "  "); err == nil {
		fmt.Println(string(b))
	}
}
