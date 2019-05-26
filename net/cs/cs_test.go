package cs

import (
	"crypto/sha1"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func hanlderr(t *testing.T) {
	e := recover()
	if e != nil {
		t.Fatal(e)
	}
}

func TestTCP(t *testing.T) {
	cs := MustCS(TCP("127.0.0.1:8100"))
	s := MustS(cs.S())
	c := MustC(cs.C())

	testCS(t, s, c)
}

func TestKCP(t *testing.T) {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	cs := MustCS(KCP("127.0.0.1:8300", key))
	s := MustS(cs.S())
	c := MustC(cs.C())

	testCS(t, s, c)
}

func TestTCPwithSMUX(t *testing.T) {
	cs := MustCS(TCP("127.0.0.1:8200"))

	ps := MustP(MustS(cs.S()).Accept())
	s := MustS(SMUX(ps).S())
	pc := MustP(MustC(cs.C()).Open())
	c := MustC(SMUX(pc).C())

	testCS(t, s, c)

	ps.Close()
	pc.Close()
}

func TestKCPwithSMUX(t *testing.T) {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	cs := MustCS(KCP("127.0.0.1:8200", key))

	ps := MustP(MustS(cs.S()).Accept())
	s := MustS(SMUX(ps).S())
	pc := MustP(MustC(cs.C()).Open())
	c := MustC(SMUX(pc).C())

	testCS(t, s, c)

	ps.Close()
	pc.Close()
}

func testCS(t *testing.T, s S, c C) {
	// defer hanlderr(t)

	go func() {
		p := MustP(s.Accept())

		bufA := make([]byte, 64)
		bufB := []byte("yes!")
		p.Read(bufA)
		fmt.Println("Server Receive:", string(bufA))
		fmt.Println("Server Send:", string(bufB))
		p.Write(bufB)

		p.Close()
		s.Close()
	}()

	time.Sleep(time.Second)

	p := MustP(c.Open())

	bufA := make([]byte, 64)
	bufB := []byte("Hello World!")
	fmt.Println("Client Send:", string(bufB))
	p.Write(bufB)
	p.Read(bufA)
	fmt.Println("Client Receive:", string(bufA))

	p.Close()
	c.Close()
}
