package protocol

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
)

type Metadata struct {
	Type     MetadataType
	Hostname string
	Port     uint16
	// Cmd is valid only if Type is MetadataTypeMsg.
	Cmd      MetadataCmd
	Cipher   string
	IsClient bool
}

type MetadataCmd uint8

const (
	MetadataCmdPing MetadataCmd = iota
	MetadataCmdSyncPassages
	MetadataCmdResponse
)

type MetadataType int

const (
	MetadataTypeIPv4 MetadataType = iota
	MetadataTypeIPv6
	MetadataTypeDomain
	MetadataTypeMsg
	MetadataTypeInvalid
)

func ParseMetadata(tgt string) (mdata Metadata, err error) {
	host, strPort, err := net.SplitHostPort(tgt)
	if err != nil {
		return mdata, fmt.Errorf("SplitHostPort: %w", err)
	}
	port, err := strconv.Atoi(strPort)
	if err != nil {
		return mdata, fmt.Errorf("failed to parse port: %w", err)
	}
	tgtIP, err := netip.ParseAddr(host)
	var typ MetadataType
	if err != nil {
		typ = MetadataTypeDomain
	} else if tgtIP.Is4() {
		typ = MetadataTypeIPv4
	} else {
		typ = MetadataTypeIPv6
	}
	return Metadata{
		Type:     typ,
		Hostname: host,
		Port:     uint16(port),
	}, nil
}
