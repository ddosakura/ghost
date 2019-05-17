package gateproxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ddosakura/ghost"
)

// HTTPConfig for Gateway
type HTTPConfig struct {
	Addr string // 监听地址

	HTTPS     *TLS       // 启用 HTTPS
	BasicAuth *BasicAuth // 启用 Basic 认证

	// 连接配置
	Proxy       func(*http.Request) (*url.URL, error) // can use http.ProxyFromEnvironment
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

// TLS for Proxy
type TLS struct {
	Crt string
	Key string
}

// BasicAuth for Proxy
type BasicAuth struct {
	User string
	Pass string
}

// InitHTTP Gateway
func InitHTTP(c *HTTPConfig) *Controller {
	if c.DialContext == nil {
		c.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	}
	transport := &http.Transport{
		Proxy:                 c.Proxy,
		DialContext:           c.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	s := &http.Server{
		Addr: c.Addr,
		Handler: handleHTTP{
			config:    c,
			transport: transport,
		},
	}
	s.RegisterOnShutdown(func() {
		ghost.Info("http server (", c.Addr, ")", "Shutdown")
	})

	controller := newController()
	controller.run = func() {
		var e error
		if c.HTTPS != nil {
			ghost.Info("https server (", c.Addr, ")")
			e = s.ListenAndServeTLS(c.HTTPS.Crt, c.HTTPS.Key)
		} else {
			ghost.Info("http server (", c.Addr, ")")
			e = s.ListenAndServe()
		}
		if e != nil {
			panic(e)
		}
	}
	controller.shutdown = func() {
		// TODO: check Hijacker
		s.Shutdown(nil)

	}
	return controller
}

type handleHTTP struct {
	config    *HTTPConfig
	transport *http.Transport
}

// https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
func (h handleHTTP) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if h.config.BasicAuth != nil {
		//u, p, ok := r.BasicAuth()
		//if !ok ||
		//	u != h.config.BasicAuth.User ||
		//	p != h.config.BasicAuth.Pass {
		//	http.Error(rw,
		//		"Username/Password Error",
		//		http.StatusProxyAuthRequired)
		//	return
		//}

		// TODO: auth

		// https://www.jb51.net/article/89070.htm
		auth := r.Header["Proxy-Authorization"]
		if auth == nil || len(auth) < 1 {
			rw.Header().Set("Proxy-Authenticate", `Basic realm="*"`)
			http.Error(rw,
				"Username/Password Error",
				http.StatusProxyAuthRequired)
			return
		}
		u, p, err := parseBaseCredential(auth[0])
		if err != nil {
			http.Error(rw,
				err.Error(),
				http.StatusProxyAuthRequired)
		}
		if u != h.config.BasicAuth.User || p != h.config.BasicAuth.Pass {
			http.Error(rw,
				"Username/Password Error",
				http.StatusProxyAuthRequired)
			return
		}
	}

	if r.Method != http.MethodConnect {
		h.xForwarded(rw, r)
		return
	}
	h.connect(rw, r)
}

// 转发 X-Forwarded-For
func (h handleHTTP) xForwarded(rw http.ResponseWriter, r *http.Request) {
	res, err := h.transport.RoundTrip(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer res.Body.Close()
	copyHeader(rw.Header(), res.Header)
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
}

// 隧道 CONNECT
func (h handleHTTP) connect(rw http.ResponseWriter, r *http.Request) {
	// TODO: HTTP/2
	if r.ProtoMajor == 2 {
		http.Error(rw,
			"HTTP/2 not supported",
			// TODO: check this code
			http.StatusInternalServerError)
	}

	// [RFC7231] 4.3.6
	// the request-target
	// consists of only the host name and port number of the tunnel
	// destination, separated by a colon.
	//
	// --- So, need't to auto add :80/:443
	//
	host := r.Host
	//if !regHasPort.MatchString(host) {
	//	switch r.URL.Scheme {
	//	case "https":
	//		host += ":443"
	//	case "http":
	//		host += ":80"
	//	default:
	//		http.Error(rw,
	//			"Need http/https",
	//			http.StatusBadRequest)
	//		return
	//	}
	//}
	//fmt.Println(host)

	conn, err := h.config.DialContext(context.Background(), "tcp", host)
	//conn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		http.Error(rw,
			err.Error(),
			http.StatusServiceUnavailable)
		return
	}

	// [RFC7231] 4.3.6 关于返回码
	// CONNECT is intended only for use in requests to a proxy.  An origin
	// server that receives a CONNECT request for itself MAY respond with a
	// 2xx (Successful) status code to indicate that a connection is
	// established.
	rw.WriteHeader(http.StatusOK)
	// [x] OK? Connection established? Connection Established?
	//rw.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	//fmt.Println(host, "conn ok!")

	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw,
			"Hijacking not supported",
			http.StatusInternalServerError)
		return
	}
	client, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(rw,
			err.Error(),
			http.StatusServiceUnavailable)
		return
	}

	go transfer(conn, client)
	go transfer(client, conn)
}
