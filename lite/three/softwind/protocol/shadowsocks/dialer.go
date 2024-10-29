package shadowsocks

import (
	"github.com/mzz2017/softwind/protocol"
	"golang.org/x/net/proxy"
	"net"
)

func init() {
	protocol.Register("shadowsocks", NewDialer)
}

type Dialer struct {
	proxyAddress string
	nextDialer   proxy.Dialer
	metadata     protocol.Metadata
	key          []byte
}

func NewDialer(nextDialer proxy.Dialer, header protocol.Header) (proxy.Dialer, error) {
	//log.Trace("shadowsocks.NewDialer: metadata: %v, password: %v", metadata, password)
	return &Dialer{
		proxyAddress: header.ProxyAddress,
		nextDialer:   nextDialer,
		metadata: protocol.Metadata{
			Cipher:   header.Cipher,
			IsClient: header.IsClient,
		},
		key: EVPBytesToKey(header.Password, CiphersConf[header.Cipher].KeyLen),
	}, nil
}

func (d *Dialer) Dial(network string, addr string) (c net.Conn, err error) {
	mdata, err := protocol.ParseMetadata(addr)
	if err != nil {
		return nil, err
	}
	mdata.Cipher = d.metadata.Cipher
	mdata.IsClient = d.metadata.IsClient
	switch network {
	case "tcp":
		conn, err := d.nextDialer.Dial(network, d.proxyAddress)
		if err != nil {
			return nil, err
		}
		return NewTCPConn(conn, mdata, d.key, nil)
	case "udp":
		conn, err := d.nextDialer.Dial(network, d.proxyAddress)
		if err != nil {
			return nil, err
		}
		return NewUDPConn(conn.(net.PacketConn), d.proxyAddress, mdata, d.key, nil)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}
