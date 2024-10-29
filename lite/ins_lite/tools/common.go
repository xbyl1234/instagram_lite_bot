package tools

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type ParseCallback = func(parse *ParseResult)

func ParseLiteStructByHexStr(isSend bool, code uint64, data string) interface{} {
	decode := DecodeHexData(data)
	p := ParseLiteStruct(isSend, code, decode)
	return p
}

func ParseLiteStruct(isSend bool, code uint64, data []byte) interface{} {
	log.Info("start parse: %d, data: %d", code, len(data))
	reader := io.CreateReaderBuffer(data)
	return parseLiteStruct(isSend, code, reader)
}

func parseLiteStruct(isSend bool, code uint64, reader *io.Reader) interface{} {
	v := proto.CreateMsgByCode(isSend, code)
	if v != nil {
		types.ReadMsg(reader, v.Interface())
		if !reader.EOF() {
			log.Error("msg code: %d, remain :%s", code, hex.EncodeToString(reader.PeekRemain()))
		}
		return v.Interface()
	} else {
		return nil
	}
}

func DecodeHexData(data string) []byte {
	s := data
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	decodeString, _ := hex.DecodeString(s)
	return decodeString
}

type ParseResult struct {
	Reader        *io.Reader
	MsgCode       uint32
	SendRecvIdx   int32
	SendIdx       int
	RecvIdx       int
	Remain        []byte
	Body          interface{}
	Data          []byte
	Magic         int32
	ClientId      int64
	ResponseInt64 int64
}

func (this *ParseResult) toString() string {
	marshal, _ := json.Marshal(this.Body)
	return fmt.Sprintf("msgCode: %v, sendIdx: 0x%x,recvIdx: 0x%x, magic: %v, clientId: %v, responseInt64: %v, json: \n%s",
		this.MsgCode,
		this.SendIdx,
		this.RecvIdx,
		this.Magic,
		this.ClientId,
		this.ResponseInt64,
		marshal)
}

func ParseLiteRecvStr(data string) *ParseResult {
	return ParseLiteRecv(DecodeHexData(data))
}

func ParseLiteSendStr(data string) *ParseResult {
	return ParseLiteSend(DecodeHexData(data))
}

func ParseLiteRecv(data []byte) *ParseResult {
	reader := io.CreateReaderBuffer(data)
	msgCode := reader.ReadVarUInt32()
	sendRecvIdx := reader.ReadInt()
	sendIdx := int(uint32(sendRecvIdx) & uint32(0x0000ffff))
	recvIdx := int((uint32(sendRecvIdx) & uint32(0xffff0000)) >> 16)
	noHeaderData := reader.PeekRemain()
	parseBody := parseLiteStruct(false, uint64(msgCode), reader)
	return &ParseResult{
		Data:        noHeaderData,
		Reader:      reader,
		MsgCode:     msgCode,
		SendRecvIdx: sendRecvIdx,
		SendIdx:     sendIdx,
		RecvIdx:     recvIdx,
		Remain:      reader.PeekRemain(),
		Body:        parseBody,
	}
}

func ParseLiteSend(data []byte) *ParseResult {
	reader := io.CreateReaderBuffer(data)
	msgCode := reader.ReadVarUInt32()
	//if msgCode != 83 && msgCode != 2 {
	//	return ""
	//}
	var magic int32
	var clientId int64
	if msgCode == 1 {
		magic = reader.ReadInt()
		clientId = reader.ReadLong()
	}
	responseInt64 := reader.ReadLong()
	sendRecvIdx := reader.ReadInt()
	sendIdx := int(uint32(sendRecvIdx) & uint32(0x0000ffff))
	recvIdx := int((uint32(sendRecvIdx) & uint32(0xffff0000)) >> 16)
	remain := reader.PeekRemain()
	parseBody := ParseLiteStruct(true, uint64(msgCode), remain)
	return &ParseResult{
		Reader:        reader,
		MsgCode:       msgCode,
		SendRecvIdx:   sendRecvIdx,
		SendIdx:       sendIdx,
		RecvIdx:       recvIdx,
		Remain:        remain,
		Body:          parseBody,
		Magic:         magic,
		ClientId:      clientId,
		ResponseInt64: responseInt64,
	}
}
