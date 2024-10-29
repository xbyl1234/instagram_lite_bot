// protocol spec:
// https://trojan-gfw.github.io/trojan/protocol

package trojanc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mzz2017/softwind/pool"
	"io"
	"net"
	"sync"
	"time"
)

var (
	CRLF        = []byte{13, 10}
	FailAuthErr = fmt.Errorf("incorrect password")
)

type Conn struct {
	net.Conn
	metadata Metadata
	pass     [56]byte

	writeMutex sync.Mutex
	onceWrite  bool
	onceRead   sync.Once
}

func NewConn(conn net.Conn, metadata Metadata, password string) (c *Conn, err error) {
	hash := sha256.New224()
	hash.Write([]byte(password))
	c = &Conn{
		Conn:     conn,
		metadata: metadata,
		pass:     [56]byte{},
	}
	hex.Encode(c.pass[:], hash.Sum(nil))
	if metadata.IsClient {
		time.AfterFunc(100*time.Millisecond, func() {
			// avoid the situation where the server sends messages first
			c.writeMutex.Lock()
			defer c.writeMutex.Unlock()
			if !c.onceWrite {
				buf := c.reqHeaderFromPool(nil)
				defer pool.Put(buf)
				if _, err = c.Conn.Write(buf); err != nil {
					return
				}
				c.onceWrite = true
			}
		})
	}
	return c, nil
}

func (c *Conn) reqHeaderFromPool(payload []byte) (buf []byte) {
	reqLen := c.metadata.Len()
	buf = pool.Get(56 + 2 + 1 + reqLen + 2 + len(payload))
	copy(buf, c.pass[:])
	copy(buf[56:], CRLF)
	buf[58] = NetworkToByte(c.metadata.Network)
	c.metadata.PackTo(buf[59:])
	copy(buf[59+reqLen:], CRLF)
	copy(buf[61+reqLen:], payload)

	return buf
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	if !c.onceWrite {
		if c.metadata.IsClient {
			buf := c.reqHeaderFromPool(b)
			defer pool.Put(buf)
			if _, err = c.Conn.Write(buf); err != nil {
				return 0, fmt.Errorf("write header: %w", err)
			}
			c.onceWrite = true
			return len(b), nil
		}
	}
	return c.Conn.Write(b)
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.onceRead.Do(func() {
		if !c.metadata.IsClient {
			if err = c.ReadReqHeader(); err != nil {
				return
			}
		}
	})
	return c.Conn.Read(b)
}

func (c *Conn) ReadReqHeader() (err error) {
	buf := pool.Get(56)
	defer pool.Put(buf)
	if _, err = io.ReadFull(c.Conn, buf); err != nil {
		return err
	}
	if !bytes.Equal(c.pass[:], buf[:56]) {
		return FailAuthErr
	}
	if _, err = io.ReadFull(c.Conn, buf[:2]); err != nil {
		return err
	}
	c.metadata.Network = ParseNetwork(buf[0])
	c.metadata.Type = ParseMetadataType(buf[1])
	n := c.metadata.Len()
	if n > cap(buf) {
		buf = pool.Get(n)
		defer pool.Put(buf)
	} else {
		buf = buf[:n]
	}
	buf[0] = MetadataTypeToByte(c.metadata.Type)
	if _, err = io.ReadFull(c.Conn, buf[1:]); err != nil {
		return err
	}
	c.metadata.Unpack(buf)
	return nil
}
