package vmess

import (
	"encoding/binary"
	"errors"
	"net"
	"net/netip"

	"github.com/mzz2017/softwind/protocol"
)

const SeqPacketMagicAddress = "sp.packet-addr.v2fly.arpa"

func (m *Metadata) IsPacketAddr() bool {
	return m.Network == "udp" && m.Type == protocol.MetadataTypeDomain && m.Hostname == SeqPacketMagicAddress
}

func ExtractPacketAddr(src []byte) (protocol.MetadataType, net.UDPAddr, error) {
	addrType := ParsePacketAddrType(src[0])

	if addrType == protocol.MetadataTypeInvalid {
		return addrType, net.UDPAddr{}, errors.New("invalid packet addr type")
	}

	if len(src) < PacketAddrLength(addrType) {
		return addrType, net.UDPAddr{}, errors.New("invalid packet addr")
	}

	addr := net.UDPAddr{}
	switch addrType {
	case protocol.MetadataTypeIPv4:
		addr.IP = net.IP(src[1:5])
		addr.Port = int(binary.BigEndian.Uint16(src[5:7]))
	case protocol.MetadataTypeIPv6:
		addr.IP = net.IP(src[1:17])
		addr.Port = int(binary.BigEndian.Uint16(src[17:19]))
	}
	return addrType, addr, nil
}

func PutPacketAddr(src []byte, addr *net.UDPAddr) error {
	nip, ok := netip.AddrFromSlice(addr.IP)
	if !ok {
		return errors.New("invalid IP")
	}

	if nip.Is4() {
		src[0] = 1
		copy(src[1:5], nip.AsSlice())
		binary.BigEndian.PutUint16(src[5:7], uint16(addr.Port))
	} else {
		src[0] = 2
		copy(src[1:17], nip.AsSlice())
		binary.BigEndian.PutUint16(src[17:19], uint16(addr.Port))
	}

	return nil
}

func ParsePacketAddrType(t byte) protocol.MetadataType {
	switch t {
	case 1:
		return protocol.MetadataTypeIPv4
	case 2:
		return protocol.MetadataTypeIPv6
	default:
		return protocol.MetadataTypeInvalid
	}
}

func UDPAddrToPacketAddrLength(addr *net.UDPAddr) int {
	nip, ok := netip.AddrFromSlice(addr.IP)
	if !ok {
		return 0
	}

	if nip.Is4() {
		return 1 + 4 + 2
	} else {
		return 1 + 16 + 2
	}
}

func PacketAddrLength(typ protocol.MetadataType) int {
	switch typ {
	case protocol.MetadataTypeIPv4:
		return 1 + 4 + 2
	case protocol.MetadataTypeIPv6:
		return 1 + 16 + 2
	default:
		return 0
	}
}
