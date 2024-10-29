package tools

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/net"
	"CentralizedControl/ins_lite/proto/io"
	"encoding/hex"
	"fmt"
	"os"
)

func ParseLiteRecvWarpPkg(data []byte, debug bool) *ParseResult {
	var parse *ParseResult
	if !debug {
		defer func() {
			if r := recover(); r != nil {
				code := uint32(0)
				if parse != nil {
					code = parse.MsgCode
				}
				log.Error("ParseLiteRecvWarpPkg code:%d error: %v, body: %s", code, r, hex.EncodeToString(data))
			}
		}()
	}
	parse = ParseLiteRecv(data)
	return parse
}

func ParseRecvStream(data string, callback ParseCallback) string {
	reader := io.CreateReaderBuffer(DecodeHexData(data))
	lzmaRead := common.NewLZMA2InputStream(65535)
	//	lzmaRead := lzma.NewReader2(16384)
	parse := ""
	for !reader.EOF() {
		var header net.ReadStreamHeader
		oneByte := reader.ReadByte()
		if (oneByte & 0x80) == 0 {
			header.IsEncode = 0
			header.MsgDataLen = uint32(oneByte&0x7f&0xff)<<8 | uint32(reader.ReadByte())
			header.StreamIdx = -1
		} else {
			switch oneByte & 0x7f {
			case 0:
				header.IsEncode = 0
				header.MsgDataLen = reader.ReadVarUInt32()
				header.StreamIdx = -1
			case 1:
				header.IsEncode = 1
				header.MsgDataLen = reader.ReadVarUInt32()
				header.StreamIdx = -1
			case 3:
				header.IsEncode = 1
				header.MsgDataLen = reader.ReadVarUInt32()
				idx := reader.ReadByte()
				header.StreamIdx = int((idx&1|2)<<(idx>>1) + 11)
			default:
				panic(fmt.Sprintf("unknow streem type: %d", oneByte&0x7f))
				//case 2:
				//case 0x7f:
			}
		}
		body := reader.ReadBytes(header.MsgDataLen)
		if header.IsEncode == 1 {
			lzmaRead.Write(body)
			all, err := lzmaRead.ReadAll()
			if err != nil {
				panic(err)
			}
			body = all
		}
		parseObj := ParseLiteRecvWarpPkg(body, true)
		if callback != nil {
			callback(parseObj)
		}
		parse += parseObj.toString() + "\n"
	}
	return parse
}

func ParseRecvFile(path string) {
	data, _ := os.ReadFile(path)
	parse := ParseRecvStream(string(data), func(parse *ParseResult) {

	})
	os.WriteFile(path+"_parse", []byte(parse), 0777)
}
