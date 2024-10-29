package trojanc

import (
	"encoding/binary"
	"fmt"
	"github.com/mzz2017/softwind/protocol"
	"github.com/mzz2017/softwind/protocol/vmess"
	"net"
)

func CompleteFromInstructionData(m *vmess.Metadata, instructionData []byte) (err error) {
	m.Type = vmess.ParseMetadataType(instructionData[3])
	switch m.Type {
	case protocol.MetadataTypeIPv4:
		m.Hostname = net.IP(instructionData[4:8]).String()
	case protocol.MetadataTypeIPv6:
		m.Hostname = net.IP(instructionData[4:20]).String()
	case protocol.MetadataTypeDomain:
		m.Hostname = string(instructionData[5 : 5+instructionData[4]])
	case protocol.MetadataTypeMsg:
		m.Cmd = protocol.MetadataCmd(instructionData[4])
	default:
		return fmt.Errorf("CompleteFromInstructionData: %w: invalid type: %v", vmess.ErrInvalidMetadata, instructionData[3])
	}
	m.Port = binary.BigEndian.Uint16(instructionData[1:])
	m.Network = vmess.ParseNetwork(instructionData[0])
	return nil
}
