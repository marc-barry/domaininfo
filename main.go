package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// IPv4ORIGINLOOKUPDNSSERVER conatains the domain of the IPv4 DNS server
const IPv4ORIGINLOOKUPDNSSERVER = "origin.asn.cymru.com"

// IPv6ORIGINLOOKUPDNSSERVER contains the domain of the IPv6 DNS server
const IPv6ORIGINLOOKUPDNSSERVER = "origin6.asn.cymru.com"

// ASNLOOKUPDNSSERVER contains the domain of the IPv4 DNS server
const ASNLOOKUPDNSSERVER = "asn.cymru.com"

// IPv4LOOKUPTEMPLATE is the template for looking up IPv4 addresses
const IPv4LOOKUPTEMPLATE = "%s." + IPv4ORIGINLOOKUPDNSSERVER

// IPv6LOOKUPTEMPLATE is the template for looking up IPv6 addresses
const IPv6LOOKUPTEMPLATE = "%s." + IPv6ORIGINLOOKUPDNSSERVER

const ASNLOOKUPTEMPLATE = "AS%s." + ASNLOOKUPDNSSERVER

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

	cname, err := net.LookupCNAME(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Canonical Name: %s\n", cname)
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
}
