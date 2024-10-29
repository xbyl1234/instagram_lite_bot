package trojanc

import (
	"encoding/binary"
	"github.com/mzz2017/softwind/pool"
	"github.com/mzz2017/softwind/protocol"
	"io"
	"net"
	"strconv"
)

func (c *Conn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	buf := pool.Get(1)
	defer pool.Put(buf)
	if _, err = io.ReadFull(c, buf); err != nil {
		return 0, nil, err
	}
	m := Metadata{Metadata: protocol.Metadata{Type: ParseMetadataType(buf[0])}}
	buf = pool.Get(m.Len())
	buf[0] = MetadataTypeToByte(m.Type)
	defer pool.Put(buf)
	if _, err = io.ReadFull(c, buf[1:]); err != nil {
		return 0, nil, err
	}
	m.Unpack(buf)

	// TODO: should we check the address type?
	if addr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(m.Hostname, strconv.Itoa(int(m.Port)))); err != nil {
		return 0, nil, err
	}
	if _, err = io.ReadFull(c, buf[:2]); err != nil {
		return 0, nil, err
	}
	length := binary.BigEndian.Uint16(buf)
	buf = pool.Get(2 + int(length))
	defer pool.Put(buf)
	if _, err = io.ReadFull(c, buf); err != nil {
		return 0, nil, err
	}
	copy(p, buf[2:])
	return int(length), addr, nil
}

func (c *Conn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	_metadata, err := protocol.ParseMetadata(addr.String())
	if err != nil {
		return 0, err
	}
	metadata := Metadata{
		Metadata: _metadata,
		Network:  "udp",
	}
	buf := pool.Get(metadata.Len() + 4 + len(p))
	defer pool.Put(buf)
	SealUDP(metadata, buf, p)
	_, err = c.Write(buf)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func SealUDP(metadata Metadata, dst []byte, data []byte) []byte {
	n := metadata.Len()
	// copy first to allow overlap
	copy(dst[n+4:], data)
	metadata.PackTo(dst)
	binary.BigEndian.PutUint16(dst[n:], uint16(len(data)))
	copy(dst[n+2:], CRLF)
	return dst[:n+4+len(data)]
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
