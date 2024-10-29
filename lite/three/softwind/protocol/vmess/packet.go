package vmess

import (
	"net"
	"strconv"

	"github.com/mzz2017/softwind/pool"
	"github.com/mzz2017/softwind/protocol"
)

func (c *Conn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	buf := pool.Get(MaxUDPSize)
	defer pool.Put(buf)
	n, err = c.Read(buf)

	if c.metadata.IsPacketAddr() {
		addrTyp, address, err := ExtractPacketAddr(buf)
		addrLen := PacketAddrLength(addrTyp)
		copy(p, buf[addrLen:n])
		return n - addrLen, &address, err
	} else {
		if c.cachedRAddrIP == nil {
			c.cachedRAddrIP, err = net.ResolveUDPAddr("udp", net.JoinHostPort(c.metadata.Hostname, strconv.Itoa(int(c.metadata.Port))))
			if err != nil {
				return 0, nil, err
			}
		}
		addr = c.cachedRAddrIP
		copy(p, buf[:n])
		return n, addr, err
	}
}

func (c *Conn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	if c.metadata.IsPacketAddr() {
		address := addr.(*net.UDPAddr)
		packetAddrLen := UDPAddrToPacketAddrLength(address)
		buf := pool.Get(packetAddrLen + len(p))
		defer pool.Put(buf)

		err := PutPacketAddr(buf, address)
		if err != nil {
			return 0, err
		}
		copy(buf[packetAddrLen:], p)
		return c.Write(buf)
	}

	return c.Write(p)
}

func (c *Conn) LocalAddr() net.Addr {
	switch c.metadata.Network {
	case "udp":
		if c.Conn.LocalAddr() != nil {
			return protocol.TCPAddrToUDPAddr(c.Conn.LocalAddr().(*net.TCPAddr))
		} else {
			return nil
		}
	default:
		return c.Conn.LocalAddr()
	}
}

func (c *Conn) RemoteAddr() net.Addr {
	switch c.metadata.Network {
	case "udp":
		if c.Conn.RemoteAddr() != nil {
			return protocol.TCPAddrToUDPAddr(c.Conn.RemoteAddr().(*net.TCPAddr))
		} else {
			return nil
		}
	default:
		return c.Conn.RemoteAddr()
	}
}
