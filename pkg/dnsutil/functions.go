package dnsutil

import (
	"net"
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
			return targets, err
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
		return ips, err
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
		return ips, err
	}
	for _, r := range res {
		ips = append(ips, r.AAAA)
	}
	return ips, nil
}
