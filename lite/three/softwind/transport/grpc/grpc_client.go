package grpc

import (
	"context"
	"fmt"
	"github.com/mzz2017/softwind/pkg/cert"
	proto "github.com/mzz2017/softwind/pkg/gun_proto"
	"github.com/mzz2017/softwind/pool"
	"golang.org/x/net/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// https://github.com/v2fly/v2ray-core/blob/v5.0.6/transport/internet/grpc/dial.go
type clientConnMeta struct {
	cc         *grpc.ClientConn
	addrTagger *addrTagger
}

var (
	globalCCMap    map[string]*clientConnMeta
	globalCCAccess sync.Mutex
)

type ccCanceller func()

type ClientConn struct {
	tun       proto.GunService_TunClient
	closer    context.CancelFunc
	muReading sync.Mutex // muReading protects reading
	muWriting sync.Mutex // muWriting protects writing
	muRecv    sync.Mutex // muReading protects recv
	muSend    sync.Mutex // muWriting protects send
	buf       []byte
	offset    int

	deadlineMu    sync.Mutex
	readDeadline  *time.Timer
	writeDeadline *time.Timer
	readClosed    chan struct{}
	writeClosed   chan struct{}
	closed        chan struct{}

	addrTagger *addrTagger
}

func NewClientConn(tun proto.GunService_TunClient, addrTagger *addrTagger, closer context.CancelFunc) *ClientConn {
	return &ClientConn{
		tun:         tun,
		closer:      closer,
		readClosed:  make(chan struct{}),
		writeClosed: make(chan struct{}),
		closed:      make(chan struct{}),
		addrTagger:  addrTagger,
	}
}

type RecvResp struct {
	hunk *proto.Hunk
	err  error
}

func (c *ClientConn) Read(p []byte) (n int, err error) {
	select {
	case <-c.readClosed:
		return 0, os.ErrDeadlineExceeded
	case <-c.closed:
		return 0, io.EOF
	default:
	}

	c.muReading.Lock()
	defer c.muReading.Unlock()
	if c.buf != nil {
		n = copy(p, c.buf[c.offset:])
		c.offset += n
		if c.offset == len(c.buf) {
			pool.Put(c.buf)
			c.buf = nil
		}
		return n, nil
	}
	// set 1 to avoid channel leak
	readDone := make(chan RecvResp, 1)
	// pass channel to the function to avoid closure leak
	go func(readDone chan RecvResp) {
		// FIXME: not really abort the send so there is some problems when recover
		c.muRecv.Lock()
		defer c.muRecv.Unlock()
		recv, e := c.tun.Recv()
		readDone <- RecvResp{
			hunk: recv,
			err:  e,
		}
	}(readDone)
	select {
	case <-c.readClosed:
		return 0, os.ErrDeadlineExceeded
	case <-c.closed:
		return 0, io.EOF
	case recvResp := <-readDone:
		err = recvResp.err
		if err != nil {
			if code := status.Code(err); code == codes.Unavailable || status.Code(err) == codes.OutOfRange {
				err = io.EOF
			}
			return 0, err
		}
		n = copy(p, recvResp.hunk.Data)
		c.buf = pool.Get(len(recvResp.hunk.Data) - n)
		copy(c.buf, recvResp.hunk.Data[n:])
		c.offset = 0
		return n, nil
	}
}

func (c *ClientConn) Write(p []byte) (n int, err error) {
	select {
	case <-c.writeClosed:
		return 0, os.ErrDeadlineExceeded
	case <-c.closed:
		return 0, io.EOF
	default:
	}

	c.muWriting.Lock()
	defer c.muWriting.Unlock()
	// set 1 to avoid channel leak
	sendDone := make(chan error, 1)
	// pass channel to the function to avoid closure leak
	go func(sendDone chan error) {
		// FIXME: not really abort the send so there is some problems when recover
		c.muSend.Lock()
		defer c.muSend.Unlock()
		e := c.tun.Send(&proto.Hunk{Data: p})
		sendDone <- e
	}(sendDone)
	select {
	case <-c.writeClosed:
		return 0, os.ErrDeadlineExceeded
	case <-c.closed:
		return 0, io.EOF
	case err = <-sendDone:
		if code := status.Code(err); code == codes.Unavailable || status.Code(err) == codes.OutOfRange {
			err = io.EOF
		}
		return len(p), err
	}
}

func (c *ClientConn) Close() error {
	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
	c.closer()
	return nil
}
func (c *ClientConn) CloseWrite() error {
	return c.tun.CloseSend()
}

// FIXME: LocalAddr is not RELIABLE.
func (c *ClientConn) LocalAddr() net.Addr {
	return c.addrTagger.ConnTagInfo.LocalAddr
}
func (c *ClientConn) RemoteAddr() net.Addr {
	return c.addrTagger.ConnTagInfo.RemoteAddr
}

func (c *ClientConn) SetDeadline(t time.Time) error {
	c.deadlineMu.Lock()
	defer c.deadlineMu.Unlock()
	if now := time.Now(); t.After(now) {
		// refresh the deadline if the deadline has been exceeded
		select {
		case <-c.readClosed:
			c.readClosed = make(chan struct{})
		default:
		}
		select {
		case <-c.writeClosed:
			c.writeClosed = make(chan struct{})
		default:
		}
		// reset the deadline timer to make the c.readClosed and c.writeClosed with the new pointer (if it is)
		if c.readDeadline != nil {
			c.readDeadline.Stop()
		}
		c.readDeadline = time.AfterFunc(t.Sub(now), func() {
			c.deadlineMu.Lock()
			defer c.deadlineMu.Unlock()
			select {
			case <-c.readClosed:
			default:
				close(c.readClosed)
			}
		})
		if c.writeDeadline != nil {
			c.writeDeadline.Stop()
		}
		c.writeDeadline = time.AfterFunc(t.Sub(now), func() {
			c.deadlineMu.Lock()
			defer c.deadlineMu.Unlock()
			select {
			case <-c.writeClosed:
			default:
				close(c.writeClosed)
			}
		})
	} else {
		select {
		case <-c.readClosed:
		default:
			close(c.readClosed)
		}
		select {
		case <-c.writeClosed:
		default:
			close(c.writeClosed)
		}
	}
	return nil
}

func (c *ClientConn) SetReadDeadline(t time.Time) error {
	c.deadlineMu.Lock()
	defer c.deadlineMu.Unlock()
	if now := time.Now(); t.After(now) {
		// refresh the deadline if the deadline has been exceeded
		select {
		case <-c.readClosed:
			c.readClosed = make(chan struct{})
		default:
		}
		// reset the deadline timer to make the c.readClosed and c.writeClosed with the new pointer (if it is)
		if c.readDeadline != nil {
			c.readDeadline.Stop()
		}
		c.readDeadline = time.AfterFunc(t.Sub(now), func() {
			c.deadlineMu.Lock()
			defer c.deadlineMu.Unlock()
			select {
			case <-c.readClosed:
			default:
				close(c.readClosed)
			}
		})
	} else {
		select {
		case <-c.readClosed:
		default:
			close(c.readClosed)
		}
	}
	return nil
}

func (c *ClientConn) SetWriteDeadline(t time.Time) error {
	c.deadlineMu.Lock()
	defer c.deadlineMu.Unlock()
	if now := time.Now(); t.After(now) {
		// refresh the deadline if the deadline has been exceeded
		select {
		case <-c.writeClosed:
			c.writeClosed = make(chan struct{})
		default:
		}
		if c.writeDeadline != nil {
			c.writeDeadline.Stop()
		}
		c.writeDeadline = time.AfterFunc(t.Sub(now), func() {
			c.deadlineMu.Lock()
			defer c.deadlineMu.Unlock()
			select {
			case <-c.writeClosed:
			default:
				close(c.writeClosed)
			}
		})
	} else {
		select {
		case <-c.writeClosed:
		default:
			close(c.writeClosed)
		}
	}
	return nil
}

type Dialer struct {
	NextDialer    proxy.ContextDialer
	ServiceName   string
	ServerName    string
	AllowInsecure bool
}

func (d *Dialer) Dial(network string, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

func (d *Dialer) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	meta, cancel, err := getGrpcClientConn(ctx, d.NextDialer, d.ServerName, address, d.AllowInsecure)
	if err != nil {
		cancel()
		return nil, err
	}
	client := proto.NewGunServiceClient(meta.cc)

	clientX := client.(proto.GunServiceClientX)
	serviceName := d.ServiceName
	if serviceName == "" {
		serviceName = "GunService"
	}
	// ctx is the lifetime of the tun
	ctxStream, streamCloser := context.WithCancel(context.Background())
	tun, err := clientX.TunCustomName(ctxStream, serviceName)
	if err != nil {
		streamCloser()
		return nil, err
	}
	return NewClientConn(tun, meta.addrTagger, streamCloser), nil
}

func getGrpcClientConn(ctx context.Context, dialer proxy.ContextDialer, serverName string, address string, allowInsecure bool) (*clientConnMeta, ccCanceller, error) {
	// allowInsecure?
	var certOption grpc.DialOption
	if allowInsecure {
		certOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		roots, err := cert.GetSystemCertPool()
		if err != nil {
			return nil, func() {}, fmt.Errorf("failed to get system certificate pool")
		}
		certOption = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(roots, serverName))
	}

	globalCCAccess.Lock()
	if globalCCMap == nil {
		globalCCMap = make(map[string]*clientConnMeta)
	}
	globalCCAccess.Unlock()

	canceller := func() {
		globalCCAccess.Lock()
		defer globalCCAccess.Unlock()
		globalCCMap[address].cc.Close()
		delete(globalCCMap, address)
	}

	// TODO Should support chain proxy to the same destination
	globalCCAccess.Lock()
	if meta, found := globalCCMap[address]; found && meta.cc.GetState() != connectivity.Shutdown {
		globalCCAccess.Unlock()
		return meta, canceller, nil
	}
	globalCCAccess.Unlock()
	meta := &clientConnMeta{
		cc:         nil,
		addrTagger: &addrTagger{},
	}
	var err error
	meta.cc, err = grpc.DialContext(ctx, address,
		certOption,
		grpc.WithContextDialer(func(ctxGrpc context.Context, s string) (net.Conn, error) {
			return dialer.DialContext(ctxGrpc, "tcp", s)
		}), grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  500 * time.Millisecond,
				Multiplier: 1.5,
				Jitter:     0.2,
				MaxDelay:   19 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}), grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}), grpc.WithStatsHandler(meta.addrTagger),
	)
	if err != nil {
		return nil, canceller, err
	}
	globalCCAccess.Lock()
	globalCCMap[address] = meta
	globalCCAccess.Unlock()
	return meta, canceller, err
}
