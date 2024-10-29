package net

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"CentralizedControl/ins_lite/config"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/three/czlib"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

type SocketIo interface {
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
	Close() error
}

type MsgWrap struct {
	Parent  *MsgWrap
	Code    uint64
	RecvIdx int
	SendIdx int
	Reader  *io.Reader
	Body    *reflect.Value
}

type SocketEvent interface {
	OnSocketConnected(isReConnect bool)
	OnSocketRecv(msgEvent *MsgWrap)
	OnSocketLost(isSend bool, err any)
}

const (
	ClientStatusSocketNotStart  = -1
	ClientStatusSocketClosed    = 0b00000000
	ClientStatusSocketLost      = 0b00000010
	ClientStatusSocketConnected = 0b00000100
)

type InsLiteSocket struct {
	proxy           proxys.Proxy
	sessionId       int64
	senderQueue     *ChanQueue[proto.MessageInterface]
	socket          SocketIo
	socketReader    *io.Reader
	lzmaReader      *common.LZMA2InputStream
	zipWriter       *czlib.ZipWriter
	senderIdx       int32
	recverIdx       int32
	callback        SocketEvent
	lastError       string
	ClientStatus    int
	closeWrite      chan byte
	isReadRunning   bool
	isWriteRunning  bool
	waitWriteThread sync.WaitGroup
	waitReadThread  sync.WaitGroup
	lock            sync.Mutex
}

var (
	FirstData       int32 = 0x1869F
	defaultDictSize       = 16384
)

func CreateInsLiteSocketClient(callback SocketEvent, p proxys.Proxy) (*InsLiteSocket, error) {
	r := common.NewLZMA2InputStream(defaultDictSize)
	w, _ := czlib.NewWriterLevel(9)
	client := &InsLiteSocket{
		proxy:        p,
		senderQueue:  CreateChanQueue[proto.MessageInterface](),
		lzmaReader:   r,
		zipWriter:    w,
		callback:     callback,
		ClientStatus: ClientStatusSocketNotStart,
	}
	return client, nil
}

func (this *InsLiteSocket) SetSession(SessionId int64, dictSize int32) {
	this.sessionId = SessionId
	if int32(defaultDictSize) != dictSize {
		this.lzmaReader = common.NewLZMA2InputStream(int(dictSize))
	}
}

func (this *InsLiteSocket) connect(transientToken int32) error {
	//s, err := CreateSocket(config.InsHost, "443")
	//s, err := CreateGoTls(config.InsHost, "443", this.proxy)
	var err error
	var s net.Conn
	for i := 0; i < 3; i++ {
		s, err = CreateJa3Socket(config.InsHost, "443", this.proxy)
		if err != nil {
			if s != nil {
				s.Close()
			}
			log.Error("connect error: %v", err)
			continue
		}
		err = nil
		break
	}
	if err != nil {
		return err
	}
	log.Info("remote ip: %s", s.RemoteAddr().String())
	wrapSocket := WrapSocketLoger(s)
	//wrapSocket := s
	defer func() {
		if err != nil {
			s.Close()
		}
	}()
	this.socket = wrapSocket
	this.socketReader = io.CreateReader(wrapSocket)

	if transientToken == 0 {
		transientToken = FirstData
	}
	buff := io.CreateWriter(4)
	buff.WriteInt(transientToken)
	_, err = this.socket.Write(buff.GetBytes())
	if err != nil {
		return err
	}
	this.closeWrite = make(chan byte)
	go this.reader()
	go this.writer()
	this.ClientStatus = ClientStatusSocketConnected
	return nil
}

func (this *InsLiteSocket) Start(transientToken int32) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	err := this.connect(transientToken)
	if err != nil {
		return err
	}
	this.callback.OnSocketConnected(false)
	return nil
}

func (this *InsLiteSocket) ReConnect(transientToken int32) error {
	if this.ClientStatus == ClientStatusSocketConnected {
		log.Warn("socket had connected!")
		return nil
	}
	err := this.connect(transientToken)
	if err != nil {
		return err
	}
	this.callback.OnSocketConnected(true)
	return nil
}

func (this *InsLiteSocket) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.ClientStatus = ClientStatusSocketClosed
	this.senderQueue.Close()
	this.socket.Close()
	//this.waitReadThread.Wait()
	//this.waitWriteThread.Wait()
}

func (this *InsLiteSocket) SendRaw(b []byte) error {
	_, err := this.socket.Write(b)
	return err
}

func (this *InsLiteSocket) checkSend() {
	if this.ClientStatus == ClientStatusSocketNotStart {
		panic(fmt.Sprintf("socket not start"))
	}
	if this.ClientStatus == ClientStatusSocketClosed {
		panic(fmt.Sprintf("socket had close: %s", this.lastError))
	}
	if this.ClientStatus == ClientStatusSocketLost {
		log.Warn("socket had lost: %s", this.lastError)
	}
}

func (this *InsLiteSocket) SendMsg(msg proto.MessageInterface) {
	this.checkSend()
	this.senderQueue.Push(msg)
}

func (this *InsLiteSocket) SendMsgFront(msg proto.MessageInterface) {
	this.checkSend()
	this.senderQueue.PushFront(msg)
}

func (this *InsLiteSocket) parseRecvBody(body []byte) {
	atomic.AddInt32(&this.recverIdx, 1)
	reader := io.CreateReaderBuffer(body)
	msgCode := reader.ReadVarUInt32()
	msgIdx := reader.ReadInt()
	recvIdx := int(uint32(msgIdx) & uint32(0x0000ffff))
	sendIdx := int((uint32(msgIdx) & uint32(0xffff0000)) >> 16)

	msgEvent := &MsgWrap{
		Code:    uint64(msgCode),
		RecvIdx: recvIdx,
		SendIdx: sendIdx,
		Reader:  reader,
	}

	this.callback.OnSocketRecv(msgEvent)
}

type ReadStreamHeader struct {
	StreamIdx  int
	MsgDataLen uint32
	IsEncode   int
}

func (this *InsLiteSocket) reader() {
	if config.EnablePanic {
		defer func() {
			this.isReadRunning = false
			r := recover()
			if this.ClientStatus == ClientStatusSocketClosed {
				return
			}
			log.Error("reader recover error: %v", r)
			this.onSocketLost(false, r)
			this.waitReadThread.Done()
		}()
	}
	this.isReadRunning = true
	this.waitReadThread.Add(1)
	for {
		var header ReadStreamHeader
		oneByte := this.socketReader.ReadByte()
		if (oneByte & 0x80) == 0 {
			header.IsEncode = 0
			header.MsgDataLen = uint32(oneByte&0x7f&0xff)<<8 | uint32(this.socketReader.ReadByte())
			header.StreamIdx = -1
		} else {
			switch oneByte & 0x7f {
			case 0:
				header.IsEncode = 0
				header.MsgDataLen = this.socketReader.ReadVarUInt32()
				header.StreamIdx = -1
			case 1:
				header.IsEncode = 1
				header.MsgDataLen = this.socketReader.ReadVarUInt32()
				header.StreamIdx = -1
			case 3:
				header.IsEncode = 1
				header.MsgDataLen = this.socketReader.ReadVarUInt32()
				idx := this.socketReader.ReadByte()
				header.StreamIdx = int((idx&1|2)<<(idx>>1) + 11)
			default:
				panic(fmt.Sprintf("unknow streem type: %d", oneByte&0x7f))
				//case 2:
				//case 0x7f:
			}
		}
		//log.Info("stream id: %d", header.StreamIdx)
		body := this.socketReader.ReadBytes(header.MsgDataLen)
		if header.IsEncode == 1 {
			this.lzmaReader.Write(body)
			all, err := this.lzmaReader.ReadAll()
			if err != nil {
				log.Error("recv lzma error: %v", err)
				panic(err)
			}
			body = all
		}
		this.parseRecvBody(body)
	}
}

func (this *InsLiteSocket) makeSendZipData(msg proto.MessageInterface) []byte {
	msgBody := msg.WriteTo()
	wrapBody := io.CreateWriter(len(msgBody) + 0x50)
	wrapBody.WriteVarUInt32(uint32(msg.GetCode(true)))

	v := reflect.ValueOf(msg).Elem()
	if strings.Contains(v.Type().Name(), "MessageC") {
		magic := v.FieldByName("Magic").Interface().(int32)
		clientId := v.FieldByName("ClientId").Interface().(int64)
		wrapBody.WriteInt(magic)
		wrapBody.WriteLong(clientId)
	}

	wrapBody.WriteLong(this.sessionId)

	idx := msg.GetSenderIdx() | (msg.GetRecverIdx() << 16) //& int32(uint32(0xFFFF0000))
	wrapBody.WriteInt(idx)

	wrapBody.WriteBytes(msgBody)

	//zip
	this.zipWriter.Write(wrapBody.GetBytes())
	this.zipWriter.Flush()
	zipData, err := this.zipWriter.Read()
	if err != nil {
		panic(err)
	}
	zipData = zipData[:len(zipData)-4]
	wrapZipData := io.CreateWriter(len(zipData) + 0x50)
	if len(zipData) > 0x7fff {
		log.Warn("len(zipData) > 0x7fff")
	} else {
		var retry = false
		retry = msg.GetCode(true) == proto.MsgCodeAppInitMsg
		var l = int16(len(zipData))
		if retry {
			wrapZipData.WriteByte(0x80)
			wrapZipData.WriteByte(0)
		} else {
			tmp := uint16(0x8000)
			l |= int16(tmp)
		}
		wrapZipData.WriteShort(l)
	}
	wrapZipData.WriteBytes(zipData)
	return wrapZipData.GetBytes()
}

func (this *InsLiteSocket) makeSendNoZipData(msg proto.MessageInterface) []byte {
	msgBody := msg.WriteTo()
	wrapBody := io.CreateWriter(len(msgBody) + 0x50)
	if len(msgBody) > 0x7fff {
		panic("len(zipData) > 0x7fff")
	} else {
		wrapBody.WriteShort(int16(len(msgBody) + io.GetVarUInt32Len(uint32(msg.GetCode(true))) + 12))
		wrapBody.WriteVarUInt32(uint32(msg.GetCode(true)))
	}
	wrapBody.WriteLong(this.sessionId)
	idx := msg.GetSenderIdx() | (msg.GetRecverIdx() << 16) // & int32(uint32(0xFFFF0000))
	wrapBody.WriteInt(idx)
	wrapBody.WriteBytes(msgBody)
	return wrapBody.GetBytes()
}

func logSend(msg proto.MessageInterface) {
	marshal, _ := json.Marshal(msg.GetBody())
	//log.Info("send msgCode: %d, senderIdx: 0x%x, recverIdx: 0x%x, body: %s",
	//	msg.GetCode(true),
	//	this.senderIdx,
	//	this.recverIdx,
	//	marshal)
	log.Info("packet send\t%s\t%s", proto.GetMessageName(true, msg.GetCode(true)), marshal)
}

func (this *InsLiteSocket) makeSendData(msg proto.MessageInterface) []byte {
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Error("makeSendData error: %v, body: %v", r, msg)
	//	}
	//}()
	if msg.GetCode(true) == proto.MsgCodeReConnect {
		msg.SetSenderIdx(0)
		msg.SetRecverIdx(atomic.LoadInt32(&this.recverIdx))
		logSend(msg)
		return this.makeSendNoZipData(msg)
	}
	msg.SetRecverIdx(atomic.LoadInt32(&this.recverIdx))
	msg.SetSenderIdx(atomic.AddInt32(&this.senderIdx, 1))
	logSend(msg)
	return this.makeSendZipData(msg)
	//if msg.GetCode(true) == lite_msg.MsgCode_ActivityResumed {
	//	return this.makeSendNoZipData(msg)
	//} else {
	//	return this.makeSendZipData(msg)
	//}
}

func (this *InsLiteSocket) writer() {
	if config.EnablePanic {
		defer func() {
			this.isWriteRunning = false
			r := recover()
			if this.ClientStatus == ClientStatusSocketClosed {
				return
			}
			log.Error("writer recover error: %v", r)
			this.onSocketLost(true, r)
			this.waitWriteThread.Done()
		}()
	}
	this.isWriteRunning = true
	this.waitWriteThread.Add(1)
	for {
		select {
		case <-this.closeWrite:
			return
		case msg := <-this.senderQueue.Get():
			count, err := this.socket.Write(this.makeSendData(msg))
			if err != nil {
				return
			}
			_ = count
		}
		//log.Debug("send msg Code: %d, count: %d", msg.GetCode(true), count)
	}
}

func (this *InsLiteSocket) onSocketLost(isSend bool, err any) {
	if !this.lock.TryLock() {
		return
	}
	defer this.lock.Unlock()
	this.ClientStatus = ClientStatusSocketLost
	this.socket.Close()    //关读
	close(this.closeWrite) //关写
	if isSend {
		this.waitReadThread.Wait() //等待读写退出
	} else {
		this.waitWriteThread.Wait() //等待读写退出
	}
	this.lastError = fmt.Sprintf("%v", err)
	this.callback.OnSocketLost(isSend, err)
}

//var logs string
//switch msgCode {
//case 233:
//case 53:
//	msg_data.readInt_AA9()
//case 40:
//	return
//case 45:
//	0CQ.read_array_len_A00(0Ig0);
//	if(v3 == 282) {
//
//	}
//case 158:
//	Received image X.0KW
//default:
//}
