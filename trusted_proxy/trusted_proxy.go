package trusted_proxy

import (
	"net"
	"strings"
)

// -----------------------------------------------------------------------------

type TrustedProxy struct {
	proxiesMap  map[string]struct{}
	proxyRanges []*net.IPNet
}

// -----------------------------------------------------------------------------

func NewTrustedProxy(proxies []string) *TrustedProxy {
	tp := TrustedProxy{
		proxiesMap:  make(map[string]struct{}),
		proxyRanges: make([]*net.IPNet, 0),
	}

	// Build trusted proxy list
	for _, ipAddress := range proxies {
		if strings.Contains(ipAddress, "/") {
			_, ipNet, err := net.ParseCIDR(ipAddress)
			if err == nil {
				tp.proxyRanges = append(tp.proxyRanges, ipNet)
			}
		} else {
			tp.proxiesMap[ipAddress] = struct{}{}
		}
	}

	// Done
	return &tp
}

func (tp *TrustedProxy) IsIpTrusted(ip net.IP) bool {
	if _, ok := tp.proxiesMap[ip.String()]; ok {
		return true
	}
	for _, ipNet := range tp.proxyRanges {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
