package cs

import (
	"net"
)

type tcp struct {
	*tcpC
}

func (b *tcp) S() (S, error) {
	l, e := net.ListenTCP("tcp", b.addr)
	if e != nil {
		return nil, e
	}
	return &tcpS{l}, nil
}

func (b *tcp) C() (C, error) {
	return b, nil
}

// TCP Protocol
func TCP(address string) (CS, error) {
	addr, e := net.ResolveTCPAddr("tcp", address)
	if e != nil {
		return nil, e
	}
	return &tcp{&tcpC{addr}}, nil
}

type tcpS struct {
	l net.Listener
}

func (s *tcpS) Accept() (P, error) {
	conn, e := s.l.Accept()
	if e != nil {
		return nil, e
	}
	return &tcpP{conn}, nil
}

func (s *tcpS) Close() {
	s.l.Close()
}

type tcpC struct {
	addr *net.TCPAddr
}

func (c *tcpC) Open() (P, error) {
	conn, e := net.DialTCP("tcp", nil, c.addr)
	if e != nil {
		return nil, e
	}
	return &tcpP{conn}, nil
}

func (c *tcpC) Close() {}

type tcpP struct {
	conn net.Conn
}

func (p *tcpP) Read(buf []byte) (int, error) {
	return p.conn.Read(buf)
}

func (p *tcpP) Write(buf []byte) (int, error) {
	return p.conn.Write(buf)
}

func (p *tcpP) Close() {
	p.conn.Close()
}

func (p *tcpP) Conn() net.Conn {
	return p.conn
}
