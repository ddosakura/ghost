package cs

import (
	"fmt"
	"net"
	"testing"

	"github.com/kr/pretty"
)

func TestP2P(t *testing.T) {
	nat, host, err := P2Ptest()

	pretty.Println(nat, host, err)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("NAT Type:", nat)
	if host != nil {
		fmt.Println("External IP Family:", host.Family())
		fmt.Println("External IP:", host.IP())
		fmt.Println("External Port:", host.Port())
	}

	fmt.Println(net.ParseIP("0.0.0.0"))
}
