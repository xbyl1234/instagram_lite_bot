package protocol

import (
	"bytes"

	"github.com/v2rayA/shadowsocksR/tools"
)

func init() {
	register("auth_aes128_sha1", NewAuthAES128SHA1)
}

func NewAuthAES128SHA1() IProtocol {
	a := &authAES128{
		salt:       "auth_aes128_sha1",
		hmac:       tools.HmacSHA1,
		hashDigest: tools.SHA1Sum,
		packID:     1,
		recvInfo: recvInfo{
			recvID: 1,
			buffer: bytes.NewBuffer(nil),
		},
	}
	return a
}
