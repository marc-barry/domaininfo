package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/marc-barry/domaininfo/pkg/dnsutil"
	"github.com/marc-barry/domaininfo/pkg/ip"
	"github.com/marc-barry/domaininfo/pkg/types"
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

var resolver = dnsutil.NewResolver()

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Requires at least one command line argument")
	}

	domain := os.Args[1]

	targetsQueue := []string{domain}
	domainToLookup := ""
	targets := []string{}

	for len(targetsQueue) != 0 {
		domainToLookup, targetsQueue = targetsQueue[0], targetsQueue[1:]
		lc, err := resolver.LookupCNAME(domainToLookup)
		if err != nil {
			log.Fatal(err)
		}

		for _, cname := range lc {
			targets = append(targets, cname.Target)
			targetsQueue = append(targetsQueue, cname.Target)
		}
	}

	cnameInfo := types.CanonicalNamesInfo{Targets: targets}

	ipv4s := make([]net.IP, 0)
	ipv6s := make([]net.IP, 0)

	res4, err := resolver.LookupA(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range res4 {
		ipv4s = append(ipv4s, r.A)
	}

	res6, err := resolver.LookupAAAA(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range res6 {
		ipv6s = append(ipv6s, r.AAAA)
	}

	addresses := types.Addresses{IPv4AddressInfo: make(map[string][]types.ASNInfo), IPv6AddressInfo: make(map[string][]types.ASNInfo)}
	asns := types.ASNs{Descriptions: make(map[string]*types.ASNDescription)}

	for _, ipv4 := range ipv4s {
		addresses.IPv4AddressInfo[ipv4.String()] = make([]types.ASNInfo, 0)
		info := addresses.IPv4AddressInfo[ipv4.String()]
		if len(ipv4) != net.IPv4len {
			ipv4 = ipv4[12:]
		}
		strs := make([]string, len(ipv4))
		for i := 0; i < len(ipv4); i++ {
			strs[len(ipv4)-1-i] = strconv.Itoa(int(ipv4[i]))
		}
		if res, err := resolver.LookupTXT(fmt.Sprintf(IPv4LOOKUPTEMPLATE, strings.Join(strs, "."))); err == nil {
			for _, r := range res {
				for _, txt := range r.Txt {
					sp := strings.Split(txt, " | ")
					if len(sp) == 5 {
						asns.Descriptions[sp[0]] = &types.ASNDescription{}
						addresses.IPv4AddressInfo[ipv4.String()] = append(info, types.ASNInfo{ASN: sp[0], AddressBlock: sp[1], Country: sp[2], InternetRegistry: sp[3], Date: sp[4]})
					}
				}
			}
		}
	}

	for _, ipv6 := range ipv6s {
		addresses.IPv6AddressInfo[ipv6.String()] = make([]types.ASNInfo, 0)
		info := addresses.IPv6AddressInfo[ipv6.String()]
		ipv6decom, err := ip.UncompressedIPv6(ipv6)
		if err != nil {
			log.Printf("Error decompressing IPv6 address: %s\n", err)
			continue
		}
		ipv6decomstrip := strings.ReplaceAll(ipv6decom, ":", "")
		strs := make([]string, len(ipv6decomstrip))
		for i := 0; i < len(ipv6decomstrip); i++ {
			strs[len(ipv6decomstrip)-1-i] = string(ipv6decomstrip[i])
		}
		if res, err := resolver.LookupTXT(fmt.Sprintf(IPv6LOOKUPTEMPLATE, strings.Join(strs, "."))); err == nil {
			for _, r := range res {
				for _, txt := range r.Txt {
					sp := strings.Split(txt, " | ")
					if len(sp) == 5 {
						asns.Descriptions[sp[0]] = &types.ASNDescription{}
						addresses.IPv6AddressInfo[ipv6.String()] = append(info, types.ASNInfo{ASN: sp[0], AddressBlock: sp[1], Country: sp[2], InternetRegistry: sp[3], Date: sp[4]})
					}
				}
			}
		}
	}

	for k, v := range asns.Descriptions {
		if res, err := resolver.LookupTXT(fmt.Sprintf(ASNLOOKUPTEMPLATE, k)); err == nil {
			for _, r := range res {
				for _, txt := range r.Txt {
					sp := strings.Split(txt, " | ")
					if len(sp) == 5 && sp[0] == k {
						v.Country = sp[1]
						v.InternetRegistry = sp[2]
						v.Date = sp[3]
						v.Org = sp[4]
					}
				}
			}
		}
	}

	caas := types.CAAs{CAAInfos: make([]types.CAAInfo, 0)}
	found := false

	res, err := resolver.LookupCAA(domain)
	if err != nil {
		log.Printf("Error looking up CAA records: %s\n", err)
	}
	info := types.CAAInfo{Domain: domain, CAs: make([]string, 0)}
	for _, r := range res {
		found = true
		info.CAs = append(info.CAs, strings.Replace(r.String(), "\t", " ", -1))
	}
	caas.CAAInfos = append(caas.CAAInfos, info)

	if !found {
		for i, domain := range cnameInfo.Targets {
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
				info.CAs = append(info.CAs, strings.Replace(r.String(), "\t", " ", -1))
			}
			caas.CAAInfos = append(caas.CAAInfos, info)
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
					info.CAs = append(info.CAs, strings.Replace(r.String(), "\t", " ", -1))
				}
				caas.CAAInfos = append(caas.CAAInfos, info)
			}
		}
	}

	dinfo := types.DomainInfo{
		Domain:             domain,
		CanonicalNamesInfo: cnameInfo,
		Addresses:          addresses,
		ASNs:               asns,
		CAAs:               caas,
	}

	if b, err := json.MarshalIndent(dinfo, "", "  "); err == nil {
		fmt.Println(string(b))
	}
}
