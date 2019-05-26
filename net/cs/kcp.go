package cs

import (
	"net"

	kcp "github.com/xtaci/kcp-go"
)

type tKCP struct {
	*kcpC
}

func (b *tKCP) S() (S, error) {
	l, err := kcp.ListenWithOptions(b.addr, b.block, 10, 3)
	if err != nil {
		return nil, err
	}
	return &kcpS{l}, nil
}

func (b *tKCP) C() (C, error) {
	return b, nil
}

// KCP Protocol
//   Need key, like: key := pbkdf2.Key([]byte("pass"), []byte("salt"), 1024, 32, sha1.New)
func KCP(addr string, key []byte) (CS, error) {
	block, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		return nil, err
	}
	return &tKCP{&kcpC{addr, block}}, nil
}

type kcpS struct {
	l *kcp.Listener
}

func (s *kcpS) Accept() (P, error) {
	session, err := s.l.AcceptKCP()
	if err != nil {
		return nil, err
	}
	return &kcpP{session}, nil
}
func (s *kcpS) Close() {
	s.l.Close()
}

type kcpC struct {
	addr  string
	block kcp.BlockCrypt
}

func (c *kcpC) Open() (P, error) {
	session, err := kcp.DialWithOptions(c.addr, c.block, 10, 3)
	if err != nil {
		return nil, err
	}
	return &kcpP{session}, nil
}

func (c *kcpC) Close() {}

type kcpP struct {
	session *kcp.UDPSession
}

func (p *kcpP) Read(b []byte) (int, error) {
	return p.session.Read(b)
}
func (p *kcpP) Write(b []byte) (int, error) {
	return p.session.Write(b)
}
func (p *kcpP) Close() {
	p.session.Close()
}
func (p *kcpP) Conn() net.Conn {
	return p.session
}
