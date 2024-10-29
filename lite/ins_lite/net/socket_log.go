package net

import (
	"CentralizedControl/common/log"
	"encoding/hex"
	"os"
)

type SocketLoger struct {
	socket SocketIo
}

var readFile *os.File

func init() {
	readFile, _ = os.OpenFile("./log/debug_stream.txt", os.O_WRONLY|os.O_CREATE, 777)
}

func WrapSocketLoger(s SocketIo) *SocketLoger {
	return &SocketLoger{
		socket: s,
	}
}

func (this *SocketLoger) Read(b []byte) (int, error) {
	read, err := this.socket.Read(b)
	if err != nil {
		log.Error("read error: %v", err)
	} else if read != len(b) {
		log.Debug("read size error: %d -> %d", len(b), read)
	} else {
		log.Debug("read: %s", hex.EncodeToString(b))
	}
	if err == nil {
		readFile.Write([]byte("read: " + hex.EncodeToString(b[:read]) + "\n"))
		readFile.Sync()
	}
	return read, err
}

func (this *SocketLoger) Write(b []byte) (int, error) {
	read, err := this.socket.Write(b)
	if err != nil {
		log.Error("send error: %v", err)
	} else {
		log.Debug("send: %s", hex.EncodeToString(b))
	}
	if err == nil {
		readFile.Write([]byte("write: " + hex.EncodeToString(b[:read]) + "\n"))
		readFile.Sync()
	}
	return read, err
}

func (this *SocketLoger) Close() error {
	return this.socket.Close()
}
