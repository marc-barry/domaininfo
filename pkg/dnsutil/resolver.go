package dnsutil

import (
	"fmt"

	"github.com/miekg/dns"
)

// Resolver represents a DNS resolver that can be used to lookup DNS records
type Resolver struct {
	c *dns.Client
}

// NewResolver constructs a new DNS resolver with an underlying DNS client
func NewResolver() *Resolver {
	r := new(Resolver)
	r.c = &dns.Client{}
	return r
}

// LookupCAA looks up CAA records for a domain
func (r *Resolver) LookupCAA(name string) ([]*dns.CAA, error) {
	var rrs []*dns.CAA

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCAA)

	rsp, _, err := r.c.Exchange(msg, "1.1.1.1:53")
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
func (r *Resolver) LookupCNAME(name string) ([]*dns.CNAME, error) {
	var rrs []*dns.CNAME

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCNAME)

	rsp, _, err := r.c.Exchange(msg, "1.1.1.1:53")
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
