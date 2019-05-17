package gateproxy

import (
	"testing"
	"time"
)

func TestHTTP(t *testing.T) {
	c := InitHTTP(&HTTPConfig{
		Addr: ":80",
	})
	c.Start()
	time.Sleep(time.Second * 10)
	c.Stop()
}
