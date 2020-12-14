package dnsutil

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/marc-barry/domaininfo/pkg/ip"
	"github.com/marc-barry/domaininfo/pkg/types"
)

// CNAMEChain produces a list containing the chain of CNAMES starting from the domain
func CNAMEChain(resolver *Resolver, domain string) ([]string, error) {
	targetsQueue := []string{domain}
	domainToLookup := ""
	targets := []string{}

	for len(targetsQueue) != 0 {
		domainToLookup, targetsQueue = targetsQueue[0], targetsQueue[1:]
		lc, err := resolver.LookupCNAME(domainToLookup)
		if err != nil {
			return nil, err
		}

		for _, cname := range lc {
			targets = append(targets, cname.Target)
			targetsQueue = append(targetsQueue, cname.Target)
		}
	}
	return targets, nil
}

// IPv4List returns a list of IPv4 addresses via A record lookups
func IPv4List(resolver *Resolver, domain string) ([]net.IP, error) {
	ips := make([]net.IP, 0)
	res, err := resolver.LookupA(domain)
	if err != nil {
		return nil, err
	}
	for _, r := range res {
		ips = append(ips, r.A)
	}
	return ips, nil
}

// IPv6List returns a list of IPv4 addresses via AAAA record lookups
func IPv6List(resolver *Resolver, domain string) ([]net.IP, error) {
	ips := make([]net.IP, 0)
	res, err := resolver.LookupAAAA(domain)
	if err != nil {
		return nil, err
	}
	for _, r := range res {
		ips = append(ips, r.AAAA)
	}
	return ips, nil
}

// NewAddressesInfo returns a an address into type based on IPv4 and IPv6 input lists
func NewAddressesInfo(resolver *Resolver, ipv4s []net.IP, ipv6s []net.IP) (map[string][]types.ASNInfo, map[string][]types.ASNInfo, map[string]string, error) {
	asns := make(map[string]string)

	ipv4Info := make(map[string][]types.ASNInfo)
	for _, ipv4 := range ipv4s {
		ipv4Info[ipv4.String()] = make([]types.ASNInfo, 0)
		if len(ipv4) != net.IPv4len {
			ipv4 = ipv4[12:]
		}
		strs := make([]string, len(ipv4))
		for i := 0; i < len(ipv4); i++ {
			strs[len(ipv4)-1-i] = strconv.Itoa(int(ipv4[i]))
		}
		if res, err := resolver.LookupTXT(fmt.Sprintf(types.IPv4LOOKUPTEMPLATE, strings.Join(strs, "."))); err == nil {
			for _, r := range res {
				for _, txt := range r.Txt {
					sp := strings.Split(txt, " | ")
					if len(sp) == 5 {
						ipv4Info[ipv4.String()] = append(ipv4Info[ipv4.String()], types.ASNInfo{ASN: sp[0], AddressBlock: sp[1], Country: sp[2], InternetRegistry: sp[3], Date: sp[4]})
						asns[sp[0]] = sp[0]
					}
				}
			}
		}
	}

	ipv6Info := make(map[string][]types.ASNInfo)
	for _, ipv6 := range ipv6s {
		ipv6Info[ipv6.String()] = make([]types.ASNInfo, 0)
		ipv6decom, err := ip.UncompressedIPv6(ipv6)
		if err != nil {
			continue
		}
		ipv6decomstrip := strings.ReplaceAll(ipv6decom, ":", "")
		strs := make([]string, len(ipv6decomstrip))
		for i := 0; i < len(ipv6decomstrip); i++ {
			strs[len(ipv6decomstrip)-1-i] = string(ipv6decomstrip[i])
		}
		if res, err := resolver.LookupTXT(fmt.Sprintf(types.IPv6LOOKUPTEMPLATE, strings.Join(strs, "."))); err == nil {
			for _, r := range res {
				for _, txt := range r.Txt {
					sp := strings.Split(txt, " | ")
					if len(sp) == 5 {
						ipv6Info[ipv6.String()] = append(ipv6Info[ipv6.String()], types.ASNInfo{ASN: sp[0], AddressBlock: sp[1], Country: sp[2], InternetRegistry: sp[3], Date: sp[4]})
						asns[sp[0]] = sp[0]
					}
				}
			}
		}
	}

	return ipv4Info, ipv6Info, asns, nil
}
