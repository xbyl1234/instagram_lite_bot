package http

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net"
	"net/url"
	"strconv"
)

// HttpProxy is an HTTP/HTTPS proxy.
type HttpProxy struct {
	TlsConfig *tls.Config
	Host      string
	HaveAuth  bool
	Username  string
	Password  string
	dialer    proxy.Dialer
}

func NewHTTPProxy(u *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	s := new(HttpProxy)
	s.Host = u.Host
	s.dialer = forward
	if u.User != nil {
		s.HaveAuth = true
		s.Username = u.User.Username()
		s.Password, _ = u.User.Password()
	}
	if u.Scheme == "https" {
		serverName := u.Query().Get("sni")
		if serverName == "" {
			serverName = u.Hostname()
		}
		skipVerify, _ := strconv.ParseBool(u.Query().Get("allowInsecure"))
		s.TlsConfig = &tls.Config{
			NextProtos:         []string{"h2", "http/1.1"},
			ServerName:         serverName,
			InsecureSkipVerify: skipVerify,
		}
	}
	return s, nil
}

func (s *HttpProxy) Dial(network, addr string) (net.Conn, error) {
	// Dial and create the https client connection.
	c, err := s.dialer.Dial("tcp", s.Host)
	if err != nil {
		return nil, err
	}
	if s.TlsConfig != nil {
		c = tls.Client(c, s.TlsConfig)
	}
	return NewConn(c, s, addr), nil
}
