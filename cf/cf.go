package cf

import (
	"net"
	"net/netip"

	"go4.org/netipx"
)

var CFIPv4s = []string{
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
}

var CFIPv6s = []string{
	"2400:cb00::/32",
	"2606:4700::/32",
	"2803:f800::/32",
	"2405:b500::/32",
	"2405:8100::/32",
	"2a06:98c0::/29",
	"2c0f:f248::/32",
}

func IsCfIP(ip net.IP) bool {
	nip, ok := netipx.FromStdIP(ip)
	if !ok {
		return false
	}
	return CFIPSet.Contains(nip)
}

// CFIPSet contains all Cloudflare IPv4 and IPv6 prefixes.
// It is initialized at program startup and can be used to test
// whether an IP address belongs to Cloudflare.
var CFIPSet *netipx.IPSet

func init() {
	builder := netipx.IPSetBuilder{}

	for _, cidr := range CFIPv4s {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			panic(err)
		}
		builder.AddPrefix(prefix)
	}

	for _, cidr := range CFIPv6s {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			panic(err)
		}
		builder.AddPrefix(prefix)
	}

	set, err := builder.IPSet()
	if err != nil {
		panic(err)
	}

	CFIPSet = set
}
