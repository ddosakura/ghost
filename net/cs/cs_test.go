package cs

import (
	"testing"
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

func TestTCPandSMUX(t *testing.T) {
	cs := MustCS(TCP("127.0.0.1:8100"))
	// TODO: be short
	s := MustS(SMUX(MustP(MustS(cs.S()).Accept())).S())
	c := MustC(SMUX(MustP(MustC(cs.C()).Open())).C())

	testCS(t, s, c)
}

func testCS(t *testing.T, s S, c C) {
	// defer hanlderr(t)

	go func() {
		p := MustP(s.Accept())

		var buf []byte
		p.Read(buf)
		p.Write(buf)

		p.Close()
		s.Close()
	}()

	p := MustP(c.Open())

	buf := []byte("Hello World!")
	p.Write(buf)
	p.Read(buf)

	p.Close()
	c.Close()
}
