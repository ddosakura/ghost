package cs

import (
	"github.com/ccding/go-stun/stun"
)

// P2Ptest discover nat type
func P2Ptest(addr string) (stun.NATType, *stun.Host, error) {
	if addr == "" {
		// default port -> 3478
		addr = "stun.pjsip.org:3478"
	}
	c := stun.NewClient()
	c.SetServerAddr(addr)
	return c.Discover()
}
