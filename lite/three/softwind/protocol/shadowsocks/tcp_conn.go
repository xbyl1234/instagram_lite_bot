package shadowsocks

import (
	"crypto/cipher"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	disk_bloom "github.com/mzz2017/disk-bloom"
	"github.com/mzz2017/softwind/common"
	"github.com/mzz2017/softwind/pool"
	"github.com/mzz2017/softwind/protocol"
	"golang.org/x/crypto/hkdf"
	"hash"
	"hash/fnv"
	"io"
	"math"
	"net"
	"sync"
	"time"
)

const (
	TCPChunkMaxLen = (1 << (16 - 2)) - 1
)

var (
	ErrFailInitCipher = fmt.Errorf("fail to initiate cipher")
)

type TCPConn struct {
	net.Conn
	metadata   protocol.Metadata
	cipherConf CipherConf
	masterKey  []byte

	cipherRead  cipher.AEAD
	cipherWrite cipher.AEAD
	onceRead    sync.Once
	onceWrite   bool
	writeMutex  sync.Mutex
	nonceRead   []byte
	nonceWrite  []byte

	readMutex   sync.Mutex
	leftToRead  []byte
	indexToRead int

	bloom *disk_bloom.FilterGroup
	sg    SaltGenerator
}

type Key struct {
	CipherConf CipherConf
	MasterKey  []byte
}

func EncryptedPayloadLen(plainTextLen int, tagLen int) int {
	n := plainTextLen / TCPChunkMaxLen
	if plainTextLen%TCPChunkMaxLen > 0 {
		n++
	}
	return plainTextLen + n*(2+tagLen+tagLen)
}

func NewTCPConn(conn net.Conn, metadata protocol.Metadata, masterKey []byte, bloom *disk_bloom.FilterGroup) (crw *TCPConn, err error) {
	conf := CiphersConf[metadata.Cipher]
	if conf.NewCipher == nil {
		return nil, fmt.Errorf("invalid CipherConf")
	}
	sg, err := GetSaltGenerator(masterKey, conf.SaltLen)
	if err != nil {
		return nil, err
	}
	// DO NOT use pool here because Close() cannot interrupt the reading or writing, which will modify the value of the pool buffer.
	key := make([]byte, len(masterKey))
	copy(key, masterKey)
	c := TCPConn{
		Conn:       conn,
		metadata:   metadata,
		cipherConf: conf,
		masterKey:  key,
		nonceRead:  make([]byte, conf.NonceLen),
		nonceWrite: make([]byte, conf.NonceLen),
		bloom:      bloom,
		sg:         sg,
	}
	if metadata.IsClient {
		time.AfterFunc(100*time.Millisecond, func() {
			// avoid the situation where the server sends messages first
			c.writeMutex.Lock()
			defer c.writeMutex.Unlock()
			if !c.onceWrite {
				buf, offset, toWrite, err := c.initWriteFromPool(nil)
				if err != nil {
					return
				}
				defer pool.Put(buf)
				defer pool.Put(toWrite)
				buf = c.seal(buf[offset:], toWrite)
				if _, err = c.Conn.Write(buf); err != nil {
					return
				}
				c.onceWrite = true
			}
		})
	}
	return &c, nil
}

func (c *TCPConn) Close() error {
	return c.Conn.Close()
}

func (c *TCPConn) Read(b []byte) (n int, err error) {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()
	c.onceRead.Do(func() {
		var salt = pool.Get(c.cipherConf.SaltLen)
		defer pool.Put(salt)
		n, err = io.ReadFull(c.Conn, salt)
		if err != nil {
			return
		}
		if c.bloom != nil {
			if c.bloom.ExistOrAdd(salt) {
				err = protocol.ErrReplayAttack
				return
			}
		}
		//log.Warn("salt: %v", hex.EncodeToString(salt))
		subKey := pool.Get(c.cipherConf.KeyLen)
		defer pool.Put(subKey)
		kdf := hkdf.New(
			sha1.New,
			c.masterKey,
			salt,
			ReusedInfo,
		)
		_, err = io.ReadFull(kdf, subKey)
		if err != nil {
			return
		}
		c.cipherRead, err = c.cipherConf.NewCipher(subKey)
	})
	if c.cipherRead == nil {
		if errors.Is(err, protocol.ErrReplayAttack) {
			return 0, fmt.Errorf("%v: %w", ErrFailInitCipher, err)
		}
		return 0, fmt.Errorf("%v: %w", ErrFailInitCipher, err)
	}
	if c.indexToRead < len(c.leftToRead) {
		n = copy(b, c.leftToRead[c.indexToRead:])
		c.indexToRead += n
		if c.indexToRead >= len(c.leftToRead) {
			// Put the buf back
			pool.Put(c.leftToRead)
		}
		return n, nil
	}
	// Chunk
	chunk, err := c.readChunkFromPool()
	if err != nil {
		return 0, err
	}
	n = copy(b, chunk)
	if n < len(chunk) {
		// Wait for the next read
		c.leftToRead = chunk
		c.indexToRead = n
	} else {
		// Full reading. Put the buf back
		pool.Put(chunk)
	}
	return n, nil
}

func (c *TCPConn) readChunkFromPool() ([]byte, error) {
	bufLen := pool.Get(2 + c.cipherConf.TagLen)
	defer pool.Put(bufLen)
	//log.Warn("len(bufLen): %v, c.nonceRead: %v", len(bufLen), c.nonceRead)
	if _, err := io.ReadFull(c.Conn, bufLen); err != nil {
		return nil, err
	}
	bLenPayload, err := c.cipherRead.Open(bufLen[:0], c.nonceRead, bufLen, nil)
	if err != nil {
		//log.Warn("read length of payload: %v: %v", protocol.ErrFailAuth, err)
		return nil, protocol.ErrFailAuth
	}
	common.BytesIncLittleEndian(c.nonceRead)
	lenPayload := binary.BigEndian.Uint16(bLenPayload)
	bufPayload := pool.Get(int(lenPayload) + c.cipherConf.TagLen) // delay putting back
	if _, err = io.ReadFull(c.Conn, bufPayload); err != nil {
		return nil, err
	}
	payload, err := c.cipherRead.Open(bufPayload[:0], c.nonceRead, bufPayload, nil)
	if err != nil {
		//log.Warn("read payload: %v: %v", protocol.ErrFailAuth, err)
		return nil, protocol.ErrFailAuth
	}
	common.BytesIncLittleEndian(c.nonceRead)
	return payload, nil
}

func (c *TCPConn) initWriteFromPool(b []byte) (buf []byte, offset int, toWrite []byte, err error) {
	var mdata = Metadata{
		Metadata: c.metadata,
	}
	var prefix, suffix []byte
	if c.metadata.Type == protocol.MetadataTypeMsg {
		mdata.LenMsgBody = uint32(len(b))
		suffix = pool.Get(CalcPaddingLen(c.masterKey, b, c.metadata.IsClient))
		defer pool.Put(suffix)
	}
	if c.metadata.IsClient || c.metadata.Type == protocol.MetadataTypeMsg {
		prefix, err = mdata.BytesFromPool()
		if err != nil {
			return nil, 0, nil, err
		}
		defer pool.Put(prefix)
	}
	toWrite = pool.Get(len(prefix) + len(b) + len(suffix))
	copy(toWrite, prefix)
	copy(toWrite[len(prefix):], b)
	copy(toWrite[len(prefix)+len(b):], suffix)

	buf = pool.Get(c.cipherConf.SaltLen + EncryptedPayloadLen(len(toWrite), c.cipherConf.TagLen))
	salt := c.sg.Get()
	copy(buf, salt)
	pool.Put(salt)
	subKey := pool.Get(c.cipherConf.KeyLen)
	defer pool.Put(subKey)
	kdf := hkdf.New(
		sha1.New,
		c.masterKey,
		buf[:c.cipherConf.SaltLen],
		ReusedInfo,
	)
	_, err = io.ReadFull(kdf, subKey)
	if err != nil {
		pool.Put(buf)
		pool.Put(toWrite)
		return nil, 0, nil, err
	}
	c.cipherWrite, err = c.cipherConf.NewCipher(subKey)
	if err != nil {
		pool.Put(buf)
		pool.Put(toWrite)
		return nil, 0, nil, err
	}
	offset += c.cipherConf.SaltLen
	if c.bloom != nil {
		c.bloom.ExistOrAdd(buf[:c.cipherConf.SaltLen])
	}
	//log.Trace("salt(%p): %v", &b, hex.EncodeToString(buf[:c.cipherConf.SaltLen]))
	return buf, offset, toWrite, nil
}

func (c *TCPConn) Write(b []byte) (n int, err error) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	var buf []byte
	var toPack []byte
	var offset int
	if !c.onceWrite {
		c.onceWrite = true
		buf, offset, toPack, err = c.initWriteFromPool(b)
		if err != nil {
			return 0, err
		}
		defer pool.Put(toPack)
	}
	if buf == nil {
		buf = pool.Get(EncryptedPayloadLen(len(b), c.cipherConf.TagLen))
		toPack = b
	}
	defer pool.Put(buf)
	if c.cipherWrite == nil {
		return 0, fmt.Errorf("%v: %w", ErrFailInitCipher, err)
	}
	c.seal(buf[offset:], toPack)
	//log.Trace("to write(%p): %v", &b, hex.EncodeToString(buf[:c.cipherConf.SaltLen]))
	_, err = c.Conn.Write(buf)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *TCPConn) seal(buf []byte, b []byte) []byte {
	offset := 0
	for i := 0; i < len(b); i += TCPChunkMaxLen {
		// write chunk
		var l = common.Min(TCPChunkMaxLen, len(b)-i)
		binary.BigEndian.PutUint16(buf[offset:], uint16(l))
		_ = c.cipherWrite.Seal(buf[offset:offset], c.nonceWrite, buf[offset:offset+2], nil)
		offset += 2 + c.cipherConf.TagLen
		common.BytesIncLittleEndian(c.nonceWrite)

		_ = c.cipherWrite.Seal(buf[offset:offset], c.nonceWrite, b[i:i+l], nil)
		offset += l + c.cipherConf.TagLen
		common.BytesIncLittleEndian(c.nonceWrite)
	}
	return buf[:offset]
}

func (c *TCPConn) ReadMetadata() (metadata Metadata, err error) {
	var firstTwoBytes = pool.Get(2)
	defer pool.Put(firstTwoBytes)
	if _, err = io.ReadFull(c, firstTwoBytes); err != nil {
		return Metadata{}, err
	}
	n, err := BytesSizeForMetadata(firstTwoBytes)
	if err != nil {
		return Metadata{}, err
	}
	var bytesMetadata = pool.Get(n)
	defer pool.Put(bytesMetadata)
	copy(bytesMetadata, firstTwoBytes)
	_, err = io.ReadFull(c, bytesMetadata[2:])
	if err != nil {
		return Metadata{}, err
	}
	mdata, err := NewMetadata(bytesMetadata)
	if err != nil {
		return Metadata{}, err
	}
	metadata = *mdata
	// complete metadata
	if !c.metadata.IsClient {
		c.metadata.Type = metadata.Metadata.Type
		c.metadata.Hostname = metadata.Metadata.Hostname
		c.metadata.Port = metadata.Metadata.Port
		if metadata.Type == protocol.MetadataTypeMsg {
			c.metadata.Cmd = protocol.MetadataCmdResponse
		} else {
			c.metadata.Cmd = metadata.Metadata.Cmd
		}
	}
	return metadata, nil
}

func CalcPaddingLen(masterKey []byte, bodyWithoutAddr []byte, req bool) (length int) {
	maxPadding := common.Max(int(10*float64(len(bodyWithoutAddr))/(1+math.Log(float64(len(bodyWithoutAddr)))))-len(bodyWithoutAddr), 0)
	if maxPadding == 0 {
		return 0
	}
	var h hash.Hash32
	if req {
		h = fnv.New32a()
	} else {
		h = fnv.New32()
	}
	h.Write(masterKey)
	h.Write(bodyWithoutAddr)
	return int(h.Sum32()) % maxPadding
}
