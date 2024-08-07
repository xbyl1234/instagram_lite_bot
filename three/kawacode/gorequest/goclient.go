package gorequest

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	http2 "github.com/bogdanfinn/fhttp/http2"
	tls "github.com/bogdanfinn/utls"
	goproxy "github.com/kawacode/goproxy"
	gostruct "github.com/kawacode/gostruct"
	gotools "github.com/kawacode/gotools"
)

// Create a client based on the protocol version
func CreateClient(bot *gostruct.BotData) (*http.Client, error) {
	var client *http.Client
	if bot.HttpRequest.Request.Protocol == "1" {
		var err error
		client, err = CreateHttp1Client(bot)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		client, err = CreateHttp2Client(bot)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

// It creates an HTTP/1.1 client with the ability to use a proxy, disable redirects, and set a timeout, that uses a custom TLS dialer that uses a custom JA3 fingerprint
func CreateHttp1Client(bot *gostruct.BotData) (*http.Client, error) {
	http1transport := http.Transport{
		DisableKeepAlives: bot.HttpRequest.Request.DisableKeepAlive,
		ForceAttemptHTTP2: bot.HttpRequest.Request.ForceAttemptHTTP2,
		PseudoHeaderOrder: bot.HttpRequest.Request.ClientProfile.PseudoHeaderOrder,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			tls.EnableWeakCiphers()
			var conn net.Conn
			if len(bot.HttpRequest.Request.Proxy) > 0 {
				dialer, err := goproxy.CreateProxyDialer(bot.HttpRequest.Request.Proxy)
				if err != nil {
					return nil, err
				}
				con, err := dialer.Dial(network, addr)
				if err != nil {
					return nil, err
				}
				conn = con
			} else {
				var err error
				conn, err = net.Dial(network, addr)
				if err != nil {
					return nil, err
				}
			}
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			config := &tls.Config{ServerName: host, InsecureSkipVerify: bot.HttpRequest.Request.InsecureSkipVerify}
			var uconn *tls.UConn
			if strings.Contains(bot.HttpRequest.Request.ClientProfile.Client.Str(), "-") {
				uconn = tls.UClient(conn, config, bot.HttpRequest.Request.ClientProfile.Client, bot.HttpRequest.Request.RandomTLSExtensionOrder, bot.HttpRequest.Request.ForceHttp1)
				if strings.Contains(bot.HttpRequest.Request.ClientProfile.Client.Str(), "CustomInternal") {
					if bot.HttpRequest.Request.ClientProfile.ClientSpec == "" {
						return nil, errors.New("missing clientspec/Ja3")
					}
					if bot.HttpRequest.Request.Protocol != "2" && bot.HttpRequest.Request.Protocol != "1" {
						bot.HttpRequest.Request.Protocol = "1"
					}
					var tlsspec *tls.ClientHelloSpec
					if strings.Contains(bot.HttpRequest.Request.ClientProfile.ClientSpec, "-") {
						var err error
						tlsspec, err = gotools.ParseJA3(bot.HttpRequest.Request.ClientProfile.ClientSpec, bot.HttpRequest.Request.Protocol)
						if err != nil {
							return nil, err
						}
					} else {
						var err error
						var data []byte
						data, err = hex.DecodeString(bot.HttpRequest.Request.ClientProfile.ClientSpec)
						if err != nil {
							return nil, err
						}
						fingerprinter := &tls.Fingerprinter{KeepPSK: true, AllowBluntMimicry: true, AlwaysAddPadding: true}
						tlsspec, err = fingerprinter.FingerprintClientHello(data)
						if err != nil {
							return nil, err
						}
					}
					if err := uconn.ApplyPreset(tlsspec); err != nil {
						return nil, err
					}
				}
			} else {
				uconn = tls.UClient(conn, config, bot.HttpRequest.Request.ClientProfile.Client, bot.HttpRequest.Request.RandomTLSExtensionOrder, bot.HttpRequest.Request.ForceHttp1)
			}
			if err := uconn.Handshake(); err != nil {
				return nil, err
			}
			return uconn, nil
		},
	}
	timeout := gotools.IsInt(bot.HttpRequest.Request.Timeout)
	var client *http.Client
	if timeout {
		timeoutsec, _ := strconv.ParseInt(bot.HttpRequest.Request.Timeout, 0, 64)
		client = &http.Client{
			Transport: &http1transport,
			Timeout:   time.Duration(time.Duration(timeoutsec) * time.Second),
		}
	} else {
		client = &http.Client{
			Transport: &http1transport,
			Timeout:   time.Duration(time.Duration(30) * time.Second),
		}
	}
	if gotools.IsInt(bot.HttpRequest.Request.MaxRedirects) {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			maxredirects, _ := strconv.ParseInt(bot.HttpRequest.Request.MaxRedirects, 0, 16)
			if len(via) >= int(maxredirects) {
				return http.ErrUseLastResponse
			}
			return nil
		}
	} else if bot.HttpRequest.Request.MaxRedirects == "false" {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			maxredirects := 10
			if len(via) >= maxredirects {
				return http.ErrUseLastResponse
			}
			return nil
		}
	}
	return client, nil
}

// It creates an HTTP2 client with the ability to use a proxy, disable redirects, and set a timeout, that uses a custom TLS dialer that uses a custom JA3 fingerprint
func CreateHttp2Client(bot *gostruct.BotData) (*http.Client, error) {
	http2transport := http2.Transport{
		StrictMaxConcurrentStreams: false,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			tls.EnableWeakCiphers()
			var conn net.Conn
			if len(bot.HttpRequest.Request.Proxy) > 0 {
				dialer, err := goproxy.CreateProxyDialer(bot.HttpRequest.Request.Proxy)
				if err != nil {
					return nil, err
				}
				con, err := dialer.Dial(network, addr)
				if err != nil {
					return nil, err
				}
				conn = con
			} else {
				var err error
				conn, err = net.Dial(network, addr)
				if err != nil {
					return nil, err
				}
			}
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			config := &tls.Config{ServerName: host, InsecureSkipVerify: bot.HttpRequest.Request.InsecureSkipVerify}
			var uconn *tls.UConn
			if strings.Contains(bot.HttpRequest.Request.ClientProfile.Client.Str(), "-") {
				uconn = tls.UClient(conn, config, bot.HttpRequest.Request.ClientProfile.Client, bot.HttpRequest.Request.RandomTLSExtensionOrder, bot.HttpRequest.Request.ForceHttp1)
				if strings.Contains(bot.HttpRequest.Request.ClientProfile.Client.Str(), "CustomInternal") {
					if bot.HttpRequest.Request.ClientProfile.ClientSpec == "" {
						return nil, errors.New("missing clientspec/Ja3")
					}
					if bot.HttpRequest.Request.Protocol != "2" && bot.HttpRequest.Request.Protocol != "1" {
						bot.HttpRequest.Request.Protocol = "2"
					}
					var tlsspec *tls.ClientHelloSpec
					if strings.Contains(bot.HttpRequest.Request.ClientProfile.ClientSpec, "-") {
						var err error
						tlsspec, err = gotools.ParseJA3(bot.HttpRequest.Request.ClientProfile.ClientSpec, bot.HttpRequest.Request.Protocol)
						if err != nil {
							return nil, err
						}
					} else {
						var err error
						var data []byte
						data, err = hex.DecodeString(bot.HttpRequest.Request.ClientProfile.ClientSpec)
						if err != nil {
							return nil, err
						}
						fingerprinter := &tls.Fingerprinter{KeepPSK: true, AllowBluntMimicry: true, AlwaysAddPadding: true}
						tlsspec, err = fingerprinter.FingerprintClientHello(data)
						if err != nil {
							return nil, err
						}
					}
					if err := uconn.ApplyPreset(tlsspec); err != nil {
						return nil, err
					}
				}
			} else {
				uconn = tls.UClient(conn, config, bot.HttpRequest.Request.ClientProfile.Client, bot.HttpRequest.Request.RandomTLSExtensionOrder, bot.HttpRequest.Request.ForceHttp1)
			}
			if err := uconn.Handshake(); err != nil {
				return nil, err
			}
			return uconn, nil
		},
		PseudoHeaderOrder: bot.HttpRequest.Request.ClientProfile.PseudoHeaderOrder,
	}
	http2transport.SettingsOrder = bot.HttpRequest.Request.ClientProfile.SettingsOrder
	http2transport.Settings = bot.HttpRequest.Request.ClientProfile.Settings
	http2transport.Priorities = bot.HttpRequest.Request.ClientProfile.Priorities
	http2transport.HeaderPriority = bot.HttpRequest.Request.ClientProfile.HeaderPriority
	http2transport.ConnectionFlow = uint32(bot.HttpRequest.Request.ClientProfile.ConnectionFlow)
	timeout := gotools.IsInt(bot.HttpRequest.Request.Timeout)
	var client *http.Client
	if timeout {
		timeoutsec, _ := strconv.ParseInt(bot.HttpRequest.Request.Timeout, 0, 64)
		client = &http.Client{
			Transport: &http2transport,
			Timeout:   time.Duration(time.Duration(timeoutsec) * time.Second),
		}
	} else {
		client = &http.Client{
			Transport: &http2transport,
			Timeout:   time.Duration(time.Duration(30) * time.Second),
		}
	}
	if gotools.IsInt(bot.HttpRequest.Request.MaxRedirects) {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			maxredirects, _ := strconv.ParseInt(bot.HttpRequest.Request.MaxRedirects, 0, 16)
			if len(via) >= int(maxredirects) {
				return http.ErrUseLastResponse
			}
			return nil
		}
	} else if bot.HttpRequest.Request.MaxRedirects == "false" {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			maxredirects := 10
			if len(via) >= maxredirects {
				return http.ErrUseLastResponse
			}
			return nil
		}
	}
	return client, nil
}
