package ip

import (
	"fmt"
	"net"
	"strings"
)

// isZeros checks if the IP is all zeroes
func isZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}
	return true
}

// UncompressedIPv6 returns a string representation of an IPv6 in an uncompressed format
func UncompressedIPv6(ip net.IP) (string, error) {
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
