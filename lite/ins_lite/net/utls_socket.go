package net

import (
	tls "github.com/bogdanfinn/utls"
	"net"
	"time"
)

type Socket struct {
	Host      string
	Port      string
	Ip        string
	client    *tls.UConn
	tlsConfig *tls.Config
}

func CreateSocket(host string, port string) (*Socket, error) {
	addr, err := net.ResolveIPAddr("ip", host)
	socket := &Socket{
		Host: host,
		Port: port,
		Ip:   addr.String(),
		tlsConfig: &tls.Config{
			ServerName: host,
		},
	}

	dialConn, err := net.DialTimeout("tcp", socket.Ip+":"+socket.Port, 20*time.Second)
	if err != nil {
		return nil, err
	}

	socket.client = tls.UClient(dialConn, socket.tlsConfig, tls.HelloCustom, false, false)
	// do not use this particular spec in production
	// make sure to generate a separate copy of ClientHelloSpec for every connection
	spec := tls.ClientHelloSpec{
		TLSVersMax: tls.VersionTLS13,
		TLSVersMin: tls.VersionTLS10,
		CipherSuites: []uint16{
			tls.GREASE_PLACEHOLDER,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_AES_128_GCM_SHA256, // tls 1.3
			tls.FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Extensions: []tls.TLSExtension{
			&tls.SNIExtension{},
			&tls.SupportedCurvesExtension{Curves: []tls.CurveID{tls.X25519, tls.CurveP256}},
			&tls.SupportedPointsExtension{SupportedPoints: []byte{0}}, // uncompressed
			&tls.SessionTicketExtension{},
			&tls.ALPNExtension{AlpnProtocols: []string{"myFancyProtocol", "http/1.1"}},
			&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
				tls.ECDSAWithP256AndSHA256,
				tls.ECDSAWithP384AndSHA384,
				tls.ECDSAWithP521AndSHA512,
				tls.PSSWithSHA256,
				tls.PSSWithSHA384,
				tls.PSSWithSHA512,
				tls.PKCS1WithSHA256,
				tls.PKCS1WithSHA384,
				tls.PKCS1WithSHA512,
				tls.ECDSAWithSHA1,
				tls.PKCS1WithSHA1}},
			&tls.KeyShareExtension{[]tls.KeyShare{
				{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
				{Group: tls.X25519},
			}},
			&tls.PSKKeyExchangeModesExtension{[]uint8{1}}, // pskModeDHE
			&tls.SupportedVersionsExtension{[]uint16{
				tls.VersionTLS13,
				tls.VersionTLS12,
				tls.VersionTLS11,
				tls.VersionTLS10}},
		},
		GetSessionID: nil,
	}
	err = socket.client.ApplyPreset(&spec)

	if err != nil {
		socket.client.Close()
		return nil, err
	}
	return socket, nil
}

func (this *Socket) Read(b []byte) (int, error) {
	return this.client.Read(b)
}

func (this *Socket) Write(b []byte) (int, error) {
	return this.client.Write(b)
}
func (this *Socket) Close() error {
	return this.client.Close()
}
