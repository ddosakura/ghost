package cs

import (
	"fmt"
	"net"

	"github.com/xtaci/smux"
)

type sMux struct {
	conn net.Conn
}

func (b *sMux) S() (S, error) {
	session, err := smux.Server(b.conn, nil)
	if err != nil {
		return nil, err
	}
	return &smuxS{
		b,
		session,
	}, nil
}

func (b *sMux) C() (C, error) {
	session, err := smux.Client(b.conn, nil)
	if err != nil {
		return nil, err
	}
	return &smuxC{
		b,
		session,
	}, nil
}

// SMUX wrapper https://github.com/xtaci/smux
func SMUX(p P) CS {
	return &sMux{conn: p.Conn()}
}

type smuxS struct {
	*sMux
	session *smux.Session
}

func (s *smuxS) Accept() (P, error) {
	stream, err := s.session.AcceptStream()
	if err != nil {
		return nil, err
	}
	return &smuxP{stream}, nil
}

func (s *smuxS) Close() {
	s.session.Close()
	// s.conn.Close()
}

type smuxC struct {
	*sMux
	session *smux.Session
}

func (c *smuxC) Open() (P, error) {
	stream, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}
	return &smuxP{stream}, nil
}
func (c *smuxC) Close() {
	c.session.Close()
	// c.conn.Close()
}

type smuxP struct {
	stream *smux.Stream
}

func (p *smuxP) Read(b []byte) {
	n, e := p.stream.Read(b)
	fmt.Println(n, e)
}

func (p *smuxP) Write(b []byte) {
	n, e := p.stream.Write(b)
	fmt.Println(n, e)
}

func (p *smuxP) Close() {
	p.stream.Close()
}

func (p *smuxP) Conn() net.Conn {
	return p.stream
}
