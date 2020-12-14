package types

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
