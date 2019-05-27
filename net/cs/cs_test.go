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
	cs := MustCS(KCP("127.0.0.1:8200", key))
	s := MustS(cs.S())
	c := MustC(cs.C())

	testCS(t, s, c)
}

func TestTCPwithSMUX(t *testing.T) {
	cs := MustCS(TCP("127.0.0.1:8300"))

	testMUX(cs, func(ps, pc P) {
		s := MustS(SMUX(ps).S())
		c := MustC(SMUX(pc).C())

		testCS(t, s, c)
	})
}

func TestKCPwithSMUX(t *testing.T) {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	cs := MustCS(KCP("127.0.0.1:8400", key))

	testMUX(cs, func(ps, pc P) {
		s := MustS(SMUX(ps).S())
		c := MustC(SMUX(pc).C())

		testCS(t, s, c)
	})
}

func testMUX(cs CS, fn func(s P, c P)) {
	s := MustS(cs.S())
	var ps, pc P
	go func() {
		ps = MustP(s.Accept())
		fmt.Println("服务器通信检测中...")
		buf := make([]byte, 64)
		ps.Write([]byte("S"))
		ps.Read(buf)
		fmt.Println("服务器 LOG", string(buf))
	}()
	time.Sleep(time.Second)
	go func() {
		pc = MustP(MustC(cs.C()).Open())
		fmt.Println("客户端通信检测中...")
		buf := make([]byte, 64)
		pc.Write([]byte("C"))
		pc.Read(buf)
		fmt.Println("客户端 LOG", string(buf))
	}()

	// 确保通信连接的建立
	// retry := 0
	for ps == nil || pc == nil {
		fmt.Println("state: ps=", ps, "; pc=", pc)
		time.Sleep(time.Second)
	}

	fmt.Println("连接复用开始")
	fn(ps, pc)

	ps.Close()
	pc.Close()
}

func testCS(t *testing.T, s S, c C) {
	// defer hanlderr(t)

	time.Sleep(time.Second)

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
