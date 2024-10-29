package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type GetClientStoredServerAccessibleKey struct {
	StringKey string
	IntKey    int32
}

type GetPkgInfo struct {
	Pkg types.ListValue[string, int32]
}

type Config89 struct {
	Config types.MapValue[string, string, int32]
}

type UpdateApp struct {
	Update byte
}

type Unknow19 struct {
	Unknow1 int32
}

type TransactionID struct {
	TransactionID int32
	Status        byte
}

type IcoLink struct {
	Unknow1 byte
	Unknow2 int64
	Unknow3 string
}

type RecvUserIcoItem struct {
	Host string
	Addr types.ListValue[types.ListValue[byte, types.VarUInt32], types.VarUInt32]
}

type RecvUserIco struct {
	IcoLink types.ListValue[IcoLink, int16]
	Unknow  types.ListValue[RecvUserIcoItem, types.VarUInt32]
}

type RecvBloks struct {
	Bloks types.ListValue[byte, types.VarUInt32]
}

type GetSystemPropertiesMsg struct {
	Properties types.ListValue[string, int16]
}

type RecvImageHeader struct {
	ImageId        uint64
	ParallelChunks uint16
	PartNumber     uint16 //offset=partNumber*14000
	TotalSize      uint32
	DataLength     uint32
	Options        byte
	Flag           int32 `ins:"Options&4 != 0"`
}

type RecvImage struct {
	RecvImageHeader
	ImageData []byte
	Url       string
}

func (this *RecvImage) Write(to io.BufferWriter) {

}

func (this *RecvImage) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.RecvImageHeader)
	if this.Options&1 == 0 {
		this.ImageData = from.ReadRemain()
	} else {
		this.Url = string(from.ReadBytes(this.DataLength))
	}
}

type Notification struct {
	Json string
}
