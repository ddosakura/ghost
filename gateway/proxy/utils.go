package gateproxy

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)

// Hop-by-hop headers
// These headers are meaningful only for a single transport-level connection
// and must not be retransmitted by proxies or cached. Such headers are:
// Connection, Keep-Alive, Proxy-Authenticate, Proxy-Authorization, TE,
// Trailer, Transfer-Encoding and Upgrade.
// Note that only hop-by-hop headers may be set using the Connection general header.
//
// [RFC7231] 4.3.6 关于响应头
// A server MUST NOT send any Transfer-Encoding or Content-Length header
// fields in a 2xx (Successful) response to CONNECT.  A client MUST
// ignore any Content-Length or Transfer-Encoding header fields received
// in a successful response to CONNECT.
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			// 8 个
			switch k {
			case "Connection", "Upgrade", "Transfer-Encoding", "Trailer",
				"Keep-Alive",
				"TE", "Proxy-Authorization",
				"Proxy-Authenticate":
			default:
				dst.Add(k, v)
			}
		}
	}
}

func parseBaseCredential(basicCredential string) (user string, pass string, err error) {
	auths := strings.SplitN(basicCredential, " ", 2)
	if len(auths) != 2 {
		return "", "", errAuthMethodUnsupported
	}
	authMethod := auths[0]
	authB64 := auths[1]
	switch authMethod {
	case "Basic":
		authstr, err := base64.StdEncoding.DecodeString(authB64)
		if err != nil {
			return "", "", err
		}
		//fmt.Println(string(authstr))
		userPwd := strings.SplitN(string(authstr), ":", 2)
		if len(userPwd) != 2 {
			return "", "", errAuthNull
		}
		user = userPwd[0]
		pass = userPwd[1]
	default:
		return "", "", errAuthMethodUnsupported
	}
	return
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	//io.Copy(destination, source)
	io.CopyBuffer(destination, source, make([]byte, 1024))
}
