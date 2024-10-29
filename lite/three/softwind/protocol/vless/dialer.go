package trojanc

import (
	"github.com/mzz2017/softwind/protocol"
	"github.com/mzz2017/softwind/protocol/vmess"
	"golang.org/x/net/proxy"
	"net"
)

func init() {
	protocol.Register("vless", NewDialer)
}

type Dialer struct {
	proxyAddress string
	nextDialer   proxy.Dialer
	metadata     protocol.Metadata
	key          []byte
}

func NewDialer(nextDialer proxy.Dialer, header protocol.Header) (proxy.Dialer, error) {
	metadata := protocol.Metadata{
		IsClient: header.IsClient,
	}
	//log.Trace("vless.NewDialer: metadata: %v, password: %v", metadata, password)
	id, err := Password2Key(header.Password)
	if err != nil {
		return nil, err
	}
	return &Dialer{
		proxyAddress: header.ProxyAddress,
		nextDialer:   nextDialer,
		metadata:     metadata,
		key:          id,
	}, nil
}

func (d *Dialer) Dial(network string, addr string) (c net.Conn, err error) {
	switch network {
	case "tcp", "udp":
		mdata, err := protocol.ParseMetadata(addr)
		if err != nil {
			return nil, err
		}
		mdata.IsClient = d.metadata.IsClient

		conn, err := d.nextDialer.Dial("tcp", d.proxyAddress)
		if err != nil {
			return nil, err
		}

		return NewConn(conn, vmess.Metadata{
			Metadata: mdata,
			Network:  network,
		}, d.key)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}
