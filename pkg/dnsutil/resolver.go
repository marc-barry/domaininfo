package dnsutil

import (
	"fmt"

	"github.com/miekg/dns"
)

// Resolver represents a DNS resolver that can be used to lookup DNS records
type Resolver struct {
	c       *dns.Client
	address string
}

// NewResolver constructs a new DNS resolver with an underlying DNS client
func NewResolver(address string) *Resolver {
	r := new(Resolver)
	r.c = &dns.Client{}
	r.address = address
	return r
}

// LookupA looks up A records for a domain
func (r *Resolver) LookupA(name string) ([]*dns.A, error) {
	var rrs []*dns.A

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeA)

	rsp, _, err := r.c.Exchange(msg, r.address)
	if err != nil {
		return nil, err
	}

	if rsp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("lookup code %s", dns.RcodeToString[rsp.Rcode])
	}

	for _, rr := range rsp.Answer {
		if a, ok := rr.(*dns.A); ok {
			rrs = append(rrs, a)
		}
	}

	return rrs, nil
}

// LookupAAAA looks up AAAA records for a domain
func (r *Resolver) LookupAAAA(name string) ([]*dns.AAAA, error) {
	var rrs []*dns.AAAA

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeAAAA)

	rsp, _, err := r.c.Exchange(msg, r.address)
	if err != nil {
		return nil, err
	}

	if rsp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("lookup code %s", dns.RcodeToString[rsp.Rcode])
	}

	for _, rr := range rsp.Answer {
		if a, ok := rr.(*dns.AAAA); ok {
			rrs = append(rrs, a)
		}
	}

	return rrs, nil
}

// LookupCAA looks up CAA records for a domain
func (r *Resolver) LookupCAA(name string) ([]*dns.CAA, error) {
	var rrs []*dns.CAA

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeCAA)

	rsp, _, err := r.c.Exchange(msg, r.address)
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

	rsp, _, err := r.c.Exchange(msg, r.address)
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

// LookupTXT looks up TXT records for a domain
func (r *Resolver) LookupTXT(name string) ([]*dns.TXT, error) {
	var rrs []*dns.TXT

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeTXT)

	rsp, _, err := r.c.Exchange(msg, r.address)
	if err != nil {
		return nil, err
	}

	if rsp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("lookup code %s", dns.RcodeToString[rsp.Rcode])
	}

	for _, rr := range rsp.Answer {
		if a, ok := rr.(*dns.TXT); ok {
			rrs = append(rrs, a)
		}
	}

	return rrs, nil
}
