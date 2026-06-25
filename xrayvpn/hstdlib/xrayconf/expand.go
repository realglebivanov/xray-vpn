package xrayconf

import "github.com/realglebivanov/hstd/hstdlib"

func InvertRules(rules []RouteRule) []RouteRule {
	out := make([]RouteRule, 0, len(rules)+1)
	for _, r := range rules {
		switch r.OutboundTag {
		case hstdlib.DirectTag:
			r.OutboundTag = hstdlib.BlockTag
		case hstdlib.ProxyTag:
			r.OutboundTag = hstdlib.DirectTag
		default:
			continue
		}
		out = append(out, r)
	}
	out = append(out, RouteRule{
		Type:        "field",
		OutboundTag: hstdlib.DirectTag,
		Network:     "tcp,udp",
	})
	return out
}

const (
	TokenRUCIDR    = "cidr:ru"
	TokenSteamCIDR = "cidr:steam"
)

func ExpandRules(rules []RouteRule, cidrsByToken map[string][]string) []RouteRule {
	out := make([]RouteRule, 0, len(rules))
	for _, r := range rules {
		if len(r.IP) > 0 {
			r.IP = expandIPs(r.IP, cidrsByToken)
			if len(r.IP) == 0 {
				continue
			}
		}
		out = append(out, r)
	}
	return out
}

func expandIPs(ips []string, cidrsByToken map[string][]string) []string {
	out := make([]string, 0, len(ips))
	for _, ip := range ips {
		if cidrs, ok := cidrsByToken[ip]; ok {
			out = append(out, cidrs...)
		} else {
			out = append(out, ip)
		}
	}
	return out
}
