package protocol

import (
	"github.com/v2rayA/shadowsocksR/ssr"
)

func init() {
	register("origin", NewOrigin)
}

type origin struct {
	ssr.ServerInfo
}

func NewOrigin() IProtocol {
	a := &origin{}
	return a
}

func (o *origin) SetServerInfo(s *ssr.ServerInfo) {
	o.ServerInfo = *s
}

func (o *origin) GetServerInfo() (s *ssr.ServerInfo) {
	return &o.ServerInfo
}

func (o *origin) PreEncrypt(data []byte) (encryptedData []byte, err error) {
	return data, nil
}

func (o *origin) PostDecrypt(data []byte) ([]byte, int, error) {
	return data, len(data), nil
}

func (o *origin) SetData(data interface{}) {

}

func (o *origin) GetData() interface{} {
	return nil
}

func (o *origin) GetOverhead() int {
	return 0
}
