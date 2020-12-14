package types

// ASNInfo contains ASN info for an IP address
type ASNInfo struct {
	ASN              string `json:"asn"`
	AddressBlock     string `json:"addressBlock"`
	Country          string `json:"country"`
	InternetRegistry string `json:"internetRegistry"`
	Date             string `json:"date"`
}

// ASNDescription contains ASN info for a specific ASN
type ASNDescription struct {
	ASN              string `json:"asn"`
	Country          string `json:"country"`
	InternetRegistry string `json:"internetRegistry"`
	Date             string `json:"date"`
	Org              string `json:"org"`
}

// CAAInfo contains certificate authority lists per domain
type CAAInfo struct {
	Domain string   `json:"domain"`
	CAs    []string `json:"cas"`
}

// DomainInfo contains all domain information
type DomainInfo struct {
	Domain                string               `json:"domain"`
	CanonicalNamesTargets []string             `json:"canonicalNamesTargets"`
	IPv4AddressInfo       map[string][]ASNInfo `json:"ipv4AddressInfo"`
	IPv6AddressInfo       map[string][]ASNInfo `json:"ipv6AddressInfo"`
	ASNDescriptions       []ASNDescription     `json:"asnDescriptions"`
	CAAInfos              []CAAInfo            `json:"caaInfos"`
}
