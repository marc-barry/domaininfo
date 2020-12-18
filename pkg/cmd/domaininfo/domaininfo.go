package domaininfo

import (
	"encoding/json"
	"fmt"

	"github.com/marc-barry/domaininfo/pkg/dnsutil"
	"github.com/marc-barry/domaininfo/pkg/types"
)

// RunCommand runs the domaininf command
func RunCommand(domain string) error {
	resolver := dnsutil.NewResolver("1.1.1.1:53")

	targets, err := dnsutil.CNAMEChain(resolver, domain)
	if err != nil {
		return err
	}

	ipv4s, err := dnsutil.IPv4List(resolver, domain)
	if err != nil {
		return err
	}
	ipv6s, err := dnsutil.IPv6List(resolver, domain)
	if err != nil {
		return err
	}

	ipv4Info, ipv6Info, asns, err := dnsutil.AddressesInfos(resolver, ipv4s, ipv6s)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(
		types.DomainInfo{
			Domain:                domain,
			CanonicalNamesTargets: targets,
			IPv4AddressInfo:       ipv4Info,
			IPv6AddressInfo:       ipv6Info,
			ASNDescriptions:       dnsutil.ASNDescriptions(resolver, asns),
			CAAInfos:              dnsutil.CAAInfos(resolver, domain, targets),
		}, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
