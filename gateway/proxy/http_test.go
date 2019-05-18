package gateproxy

import (
	"fmt"
	"testing"
	"time"
)

func TestHTTP(t *testing.T) {
	c := InitHTTP(&HTTPConfig{
		Addr: ":80",
	})
	c.Sub("a", func(bool) {})
	c.Sub("a", nil)
	c.Sub("b", func(b bool) {
		fmt.Println("Alive:", b)
	})
	c.Start()
	time.Sleep(time.Second * 10)
	c.Stop()
}
