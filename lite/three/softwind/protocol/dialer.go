package protocol

import (
	"fmt"
	"golang.org/x/net/proxy"
	"strconv"
)

type Creator func(nextDialer proxy.Dialer, header Header) (proxy.Dialer, error)

var Mapper = make(map[string]Creator)

func Register(name string, c Creator) {
	Mapper[name] = c
}

func NewDialer(name string, nextDialer proxy.Dialer, header Header) (proxy.Dialer, error) {
	creator, ok := Mapper[name]
	if !ok {
		return nil, fmt.Errorf("no conn creator registered for %v", strconv.Quote(name))
	}
	return creator(nextDialer, header)
}
