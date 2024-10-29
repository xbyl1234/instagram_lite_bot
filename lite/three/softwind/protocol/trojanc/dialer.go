package trojanc

import (
	"github.com/mzz2017/softwind/protocol"
	"golang.org/x/net/proxy"
	"net"
)

func init() {
	protocol.Register("trojanc", NewDialer)
}

type Dialer struct {
	proxyAddress string
	nextDialer   proxy.Dialer
	metadata     protocol.Metadata
	password     string
}

func NewDialer(nextDialer proxy.Dialer, header protocol.Header) (proxy.Dialer, error) {
	metadata := protocol.Metadata{
		IsClient: header.IsClient,
	}
	//log.Trace("trojanc.NewDialer: metadata: %v, password: %v", metadata, password)
	return &Dialer{
		proxyAddress: header.ProxyAddress,
		nextDialer:   nextDialer,
		metadata:     metadata,
		password:     header.Password,
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

		return NewConn(conn, Metadata{
			Metadata: mdata,
			Network:  network,
		}, d.password)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}
