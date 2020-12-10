package types

type CanonicalNamesInfo struct {
	Targets []string `json:"targets"`
}

type ASNInfo struct {
	ASN              string `json:"asn"`
	AddressBlock     string `json:"addressBlock"`
	Country          string `json:"country"`
	InternetRegistry string `json:"internetRegistry"`
	Date             string `json:"date"`
}

type Addresses struct {
	IPv4AddressInfo map[string][]ASNInfo `json:"ipv4AddressInfo"`
	IPv6AddressInfo map[string][]ASNInfo `json:"ipv6AddressInfo"`
}

type ASNDescription struct {
	Country          string `json:"country"`
	InternetRegistry string `json:"internetRegistry"`
	Date             string `json:"date"`
	Org              string `json:"org"`
}

type ASNs struct {
	Descriptions map[string]*ASNDescription `json:"descriptions"`
}

type CAAInfo struct {
	Domain string   `json:"domain"`
	CAs    []string `json:"cas"`
}

type CAAs struct {
	CAAInfos []CAAInfo `json:"caaInfos"`
}

type DomainInfo struct {
	Domain             string             `json:"domain"`
	CanonicalNamesInfo CanonicalNamesInfo `json:"canonicalNamesInfo"`
	Addresses          Addresses          `json:"addresses"`
	ASNs               ASNs               `json:"asns"`
	CAAs               CAAs               `json:"caas"`
}
