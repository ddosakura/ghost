package cs

import (
	"github.com/ccding/go-stun/stun"
)

// P2Ptest discover nat type
func P2Ptest() (stun.NATType, *stun.Host, error) {
	c := stun.NewClient()
	//c.SetServerAddr("stun4.l.google.com:19302")
	c.SetServerAddr("stun.pjsip.org:3478") // default port -> 3478
	//c.SetVerbose(true)
	//c.SetVVerbose(true)

	return c.Discover()
}
