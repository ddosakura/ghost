package cs

import (
	"net"
)

// CS Builder
type CS interface {
	S() (S, error) // Listen
	C() (C, error) // Client
}

// S wrapper Server
type S interface {
	Accept() (P, error) // Open Conn
	Close()             // Close Listen
}

// C wrapper Client
type C interface {
	Open() (P, error) // Open Conn
	Close()
}

// P wrapper pipe
type P interface {
	// TODO: R/W
	Read([]byte)  // Read from Conn
	Write([]byte) // Write to Conn
	Close()       // Close Conn

	Conn() net.Conn // Get Conn
}

// MustCS handle error
func MustCS(m CS, e error) CS {
	if e != nil {
		panic(e)
	}
	return m
}

// MustS handle error
func MustS(m S, e error) S {
	if e != nil {
		panic(e)
	}
	return m
}

// MustC handle error
func MustC(m C, e error) C {
	if e != nil {
		panic(e)
	}
	return m
}

// MustP handle error
func MustP(m P, e error) P {
	if e != nil {
		panic(e)
	}
	return m
}
