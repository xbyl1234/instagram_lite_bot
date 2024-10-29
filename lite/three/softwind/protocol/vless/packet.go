package trojanc

import (
	"encoding/binary"
	"fmt"
	"github.com/mzz2017/softwind/pool"
	"github.com/mzz2017/softwind/protocol"
	"io"
	"net"
	"strconv"
)

func (c *Conn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	// FIXME: a compromise on Symmetric NAT
	if c.cachedRAddrIP == nil {
		c.cachedRAddrIP, err = net.ResolveUDPAddr("udp", net.JoinHostPort(c.metadata.Hostname, strconv.Itoa(int(c.metadata.Port))))
		if err != nil {
			return 0, nil, err
		}
	}
	addr = c.cachedRAddrIP

	bLen := pool.Get(2)
	defer pool.Put(bLen)
	if _, err = io.ReadFull(c, bLen); err != nil {
		return 0, nil, err
	}
	length := int(binary.BigEndian.Uint16(bLen))
	if len(p) < length {
		return 0, nil, fmt.Errorf("buf size is not enough")
	}
	n, err = io.ReadFull(c, p[:length])
	return n, addr, err
}

func (c *Conn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	bLen := pool.Get(2)
	defer pool.Put(bLen)
	binary.BigEndian.PutUint16(bLen, uint16(len(p)))
	if _, err = c.Write(bLen); err != nil {
		return 0, err
	}
	return c.Write(p)
}

func (c *Conn) LocalAddr() net.Addr {
	switch c.metadata.Network {
	case "udp":
		return protocol.TCPAddrToUDPAddr(c.Conn.LocalAddr().(*net.TCPAddr))
	default:
		return c.Conn.LocalAddr()
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	switch c.metadata.Network {
	case "udp":
		return protocol.TCPAddrToUDPAddr(c.Conn.RemoteAddr().(*net.TCPAddr))
	default:
		return c.Conn.RemoteAddr()
	}
}
