package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

// IPv4ORIGINLOOKUPDNSSERVER conatains the domain of the IPv4 DNS lookup server
const IPv4ORIGINLOOKUPDNSSERVER = "origin.asn.cymru.com"

// IPv6ORIGINLOOKUPDNSSERVER contains the domain of the IPv6 DNS lookup server
const IPv6ORIGINLOOKUPDNSSERVER = "origin6.asn.cymru.com"

// ASNLOOKUPDNSSERVER contains the domain of the ASN lookup server
const ASNLOOKUPDNSSERVER = "asn.cymru.com"

// IPv4LOOKUPTEMPLATE is the template for looking up IPv4 addresses
const IPv4LOOKUPTEMPLATE = "%s." + IPv4ORIGINLOOKUPDNSSERVER

// IPv6LOOKUPTEMPLATE is the template for looking up IPv6 addresses
const IPv6LOOKUPTEMPLATE = "%s." + IPv6ORIGINLOOKUPDNSSERVER

// ASNLOOKUPTEMPLATE is the template for looking up ASN descriptions
const ASNLOOKUPTEMPLATE = "AS%s." + ASNLOOKUPDNSSERVER

var dnsClient = &dns.Client{}

// LookupCAA looks up CAA records for a domain
func LookupCAA(c *dns.Client, name string) ([]*dns.CAA, error) {
	var rrs []*dns.CAA

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCAA)

	rsp, _, err := c.Exchange(msg, "1.1.1.1:53")
	if err != nil {
		return nil, err
	}

	if rsp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("lookup code %s", dns.RcodeToString[rsp.Rcode])
	}

	for _, rr := range rsp.Answer {
		if a, ok := rr.(*dns.CAA); ok {
			rrs = append(rrs, a)
		}
	}

	return rrs, nil
}

// LookupCNAME looks up CNAME records for a domain
func LookupCNAME(c *dns.Client, name string) ([]*dns.CNAME, error) {
	var rrs []*dns.CNAME

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCNAME)

	rsp, _, err := c.Exchange(msg, "1.1.1.1:53")
	if err != nil {
		return nil, err
	}

	if rsp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("lookup code %s", dns.RcodeToString[rsp.Rcode])
	}

	for _, rr := range rsp.Answer {
		if a, ok := rr.(*dns.CNAME); ok {
			rrs = append(rrs, a)
		}
	}

	return rrs, nil
}

func isZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}
	return true
}

func decompressIPv6(ip net.IP) (string, error) {
	if len(ip) == net.IPv4len {
		return "", fmt.Errorf("Not an IPv6 address")
	}
	if len(ip) != net.IPv6len {
		return "", fmt.Errorf("Not an IPv6 address. Got length %d", len(ip))
	}
	if len(ip) == net.IPv6len &&
		isZeros(ip[0:10]) &&
		ip[10] == 0xff &&
		ip[11] == 0xff {
		return "", fmt.Errorf("Not an IPv6 address")
	}
	return strings.ToLower(fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15])), nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Requires at least one command line argument")
	}

	domains := []string{os.Args[1]}
	domain := ""
	cnames := []string{}

	for len(domains) != 0 {
		domain, domains = domains[0], domains[1:]
		lc, err := LookupCNAME(dnsClient, domain)
		if err != nil {
			log.Fatal(err)
		}

		for _, cname := range lc {
			cnames = append(cnames, cname.Target)
			domains = append(domains, cname.Target)
		}
	}
	for i, cname := range cnames {
		fmt.Printf("Canonical Name (%d): %s\n", i, cname)
	}

	fmt.Println("---")

	ips, err := net.LookupIP(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ipv4s := make([]net.IP, 0)
	ipv6s := make([]net.IP, 0)
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4s = append(ipv4s, ip)
			continue
		}
		ipv6s = append(ipv6s, ip)
	}

	asnSet := map[string]string{}

	for _, ipv4 := range ipv4s {
		fmt.Printf("IPv4 Address: %s\n", ipv4.String())
		if len(ipv4) != net.IPv4len {
			ipv4 = ipv4[12:]
		}
		strs := make([]string, len(ipv4))
		for i := 0; i < len(ipv4); i++ {
			strs[len(ipv4)-1-i] = strconv.Itoa(int(ipv4[i]))
		}
		if res, err := net.LookupTXT(fmt.Sprintf(IPv4LOOKUPTEMPLATE, strings.Join(strs, "."))); err == nil {
			fmt.Printf("ASN Info: %s\n", res)
			for _, r := range res {
				sp := strings.Split(r, " | ")
				if len(sp) > 0 {
					asnSet[sp[0]] = sp[0]
				}
			}
		}
	}

	fmt.Println("---")

	for _, ipv6 := range ipv6s {
		fmt.Printf("IPv6 Address: %s\n", ipv6.String())
		ipv6decom, err := decompressIPv6(ipv6)
		if err != nil {
			fmt.Printf("Error decompressing IPv6 address: %s\n", err)
			continue
		}
		ipv6decomstrip := strings.ReplaceAll(ipv6decom, ":", "")
		strs := make([]string, len(ipv6decomstrip))
		for i := 0; i < len(ipv6decomstrip); i++ {
			strs[len(ipv6decomstrip)-1-i] = string(ipv6decomstrip[i])
		}
		if res, err := net.LookupTXT(fmt.Sprintf(IPv6LOOKUPTEMPLATE, strings.Join(strs, "."))); err == nil {
			fmt.Printf("ASN Info: %s\n", res)
			for _, r := range res {
				sp := strings.Split(r, " | ")
				if len(sp) > 0 {
					asnSet[sp[0]] = sp[0]
				}
			}
		}
	}

	fmt.Println("---")

	for k := range asnSet {
		if res, err := net.LookupTXT(fmt.Sprintf(ASNLOOKUPTEMPLATE, k)); err == nil {
			fmt.Printf("ASN Description: %s\n", res)
		}
	}

	fmt.Println("---")

	caas, err := LookupCAA(dnsClient, os.Args[1])
	if err != nil {
		fmt.Printf("Error looking up CAA records: %s\n", err)
	}
	if len(caas) == 0 {
		fmt.Printf("No CAA record on %s\n", os.Args[1])
	}
	for _, r := range caas {
		fmt.Printf("CAA Record: %s\n", r.String())
	}

	curDomains := []string{}
	done := true
	if len(caas) == 0 {
		curDomains = append(curDomains, os.Args[1])
		done = false
	}
	i := 0

	for (!done || i <= 8) && len(curDomains) != 0 {
		i++
		cnames := []*dns.CNAME{}
		for _, domain := range curDomains {
			lc, err := LookupCNAME(dnsClient, domain)
			if err != nil {
				fmt.Printf("Error looking up CNAME records: %s\n", err)
				continue
			}
			cnames = append(cnames, lc...)
		}
		curDomains = []string{}
		if len(cnames) != 0 {
			for _, r := range cnames {
				caas, err := LookupCAA(dnsClient, r.Target)
				if err != nil {
					fmt.Printf("Error looking up CAA records: %s\n", err)
					continue
				}
				if len(caas) == 0 {
					fmt.Printf("No CAA record on %s\n", r.Target)
					continue
				}
				done = true
				for _, r := range caas {
					fmt.Printf("CAA Record: %s\n", r.String())
				}
			}
		}
	}
	if !done {
		i := strings.IndexAny(os.Args[1], ".")
		if i > 0 {
			parent := os.Args[1][i+1:]
			caas, err := LookupCAA(dnsClient, parent)
			if err != nil {
				fmt.Printf("Error looking up CAA records: %s\n", err)
			}
			if len(caas) == 0 {
				fmt.Printf("No CAA record on %s\n", parent)
			}
			for _, r := range caas {
				fmt.Printf("CAA Record: %s\n", r.String())
			}
		}
	}
}
