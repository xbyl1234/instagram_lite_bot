package tools

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/three/czlib"
	"encoding/hex"
	"os"
)

func ParseLiteSendWarpPkg(data []byte, debug bool) *ParseResult {
	var parse *ParseResult
	if !debug {
		defer func() {
			if r := recover(); r != nil {
				code := uint32(0)
				if parse != nil {
					code = parse.MsgCode
				}
				log.Error("ParseLiteSendWarpPkg code:%d error: %v, body: %s", code, r, hex.EncodeToString(data))
			}
		}()
	}
	parse = ParseLiteSend(data)
	return parse
}

func ParseSendStream(data []byte, callback ParseCallback) string {
	data = data[4:]
	dezip, _ := czlib.NewReader()
	idx := 0
	reader := io.CreateReaderBuffer(data)
	var parse string
	for !reader.EOF() {
		readShort := int16(reader.ReadShort())
		l := readShort & (int16(^uint16(0x8000)))
		if l == 0 && idx == 0 {
			l = reader.ReadShort()
		}
		pkg := reader.ReadBytes(uint32(l))
		var decode []byte
		tmp := 0x8000
		if readShort&int16(uint16(tmp)) == 0 {
			decode = pkg
		} else {
			dezip.SetInput(pkg)
			dezip.SetInput([]byte{0, 0, 0xFF, 0xFF})
			inflate, err := dezip.Inflate()
			if err != nil {
				panic(err)
			}
			decode = inflate
		}
		parseObj := ParseLiteSendWarpPkg(decode, true)
		callback(parseObj)
		parse += parseObj.toString() + "\n"
		idx++
	}
	return parse
}

func ParseSendFile(path string) {
	data, _ := os.ReadFile(path)
	parse := ParseSendStream(DecodeHexData(string(data)), func(parse *ParseResult) {

	})
	os.WriteFile(path+"_parse", []byte(parse), 0777)
}
