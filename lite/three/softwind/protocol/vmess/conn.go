package vmess

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/mzz2017/softwind/common"
	"github.com/mzz2017/softwind/pkg/fastrand"
	"github.com/mzz2017/softwind/pool"
	"io"
	"net"
	"sync"
)

const (
	MaxChunkSize = 1 << 14
	MaxUDPSize   = 1 << 11
)

type Conn struct {
	net.Conn
	initRead      sync.Once
	initWrite     sync.Once
	metadata      Metadata
	cmdKey        []byte
	cachedRAddrIP *net.UDPAddr

	NewAEAD func(key []byte) (cipher.AEAD, error)

	writeMutex            sync.Mutex
	writeBodyCipher       cipher.AEAD
	writeNonceGenerator   BytesGenerator
	writeChunkSizeParser  ChunkSizeEncoder
	writePaddingGenerator PaddingLengthGenerator

	readBodyCipher       cipher.AEAD
	readNonceGenerator   BytesGenerator
	readChunkSizeParser  ChunkSizeDecoder
	readPaddingGenerator PaddingLengthGenerator

	requestBodyKey [16]byte
	requestBodyIV  [16]byte
	requestOptions byte

	responseBodyKey [16]byte
	responseBodyIV  [16]byte
	responseAuth    byte

	readMutex   sync.Mutex
	leftToRead  []byte
	indexToRead int
}

func NewConn(conn net.Conn, metadata Metadata, cmdKey []byte) (c *Conn, err error) {
	// DO NOT use pool here because Close() cannot interrupt the reading or writing, which will modify the value of the pool buffer.
	key := make([]byte, len(cmdKey))
	copy(key, cmdKey)
	c = &Conn{
		Conn:     conn,
		metadata: metadata,
		cmdKey:   key,
	}
	if metadata.IsClient {
		if err = c.WriteReqHeader(); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}

func (c *Conn) chunks(size int) (payloadSize int, numChunks int) {
	payloadSize = MaxChunkSize - c.writeBodyCipher.Overhead() - int(c.writeChunkSizeParser.SizeBytes()) - int(c.writePaddingGenerator.MaxPaddingLen())
	if size%payloadSize == 0 {
		return payloadSize, size / payloadSize
	}
	return payloadSize, size/payloadSize + 1
}

func GenerateChunkNonce(nonce []byte, size uint32) BytesGenerator {
	c := make([]byte, size)
	copy(c[2:], nonce[2:])
	count := uint16(0)
	return func() []byte {
		binary.BigEndian.PutUint16(c, count)
		count++
		return c[:size]
	}
}

// seal packs the b. The overhead is sizeParser.SizeBytes() + auth.Overhead() + paddingSize(no more than maxPadding).
func (c *Conn) sealFromPool(b []byte) (data []byte) {
	sizeSize := c.writeChunkSizeParser.SizeBytes()
	encryptedSize := int32(len(b) + c.writeBodyCipher.Overhead())
	paddingSize := int32(c.writePaddingGenerator.NextPaddingLen())

	data = pool.Get(int(sizeSize + encryptedSize + paddingSize))
	c.writeChunkSizeParser.Encode(uint16(encryptedSize+paddingSize), data)

	c.writeBodyCipher.Seal(data[sizeSize:sizeSize], c.writeNonceGenerator(), b, nil)
	fastrand.Read(data[len(data)-int(paddingSize):])
	//log.Warn("write: size: %v, padding: %v", encryptedSize+paddingSize, paddingSize)
	return data
}

// writeStream splits mb into multiple FIXED size (payloadSize) chunks.
// Then seal the chunks and write separately.
// If the sum size of mb less than one payloadSize, seal and write it directly.
func (c *Conn) writeStream(b []byte, preWrite []byte) (n int, err error) {
	payloadSize, numChunks := c.chunks(len(b))
	var start = 0
	if preWrite != nil {
		start++
		data := c.sealFromPool(b[n:common.Min(n+payloadSize, len(b))])
		defer pool.Put(data)
		if _, err = c.Conn.Write(bytes.Join([][]byte{preWrite, data}, nil)); err != nil {
			return 0, err
		}
		n += payloadSize
	}
	for i := start; i < numChunks; i++ {
		data := c.sealFromPool(b[n:common.Min(n+payloadSize, len(b))])
		if _, err = c.Conn.Write(data); err != nil {
			return n, err
		}
		pool.Put(data)
		n += payloadSize
	}
	if n > len(b) {
		n = len(b)
	}
	return n, nil
}

// writePacket simply seal every buffer of mb and write.
func (c *Conn) writePacket(b []byte, preWrite []byte) (n int, err error) {
	data := c.sealFromPool(b)
	defer pool.Put(data)
	if preWrite != nil {
		if _, err = c.Conn.Write(bytes.Join([][]byte{preWrite, data}, nil)); err != nil {
			return 0, err
		}
	} else {
		if _, err = c.Conn.Write(data); err != nil {
			return 0, err
		}
	}
	return len(b), nil
}

func (c *Conn) InitContext(instructionData []byte) error {
	c.responseAuth = instructionData[33]
	copy(c.requestBodyIV[:], instructionData[1:])
	copy(c.requestBodyKey[:], instructionData[17:])
	tmp := sha256.Sum256(c.requestBodyIV[:])
	copy(c.responseBodyIV[:], tmp[:16])
	tmp = sha256.Sum256(c.requestBodyKey[:])
	copy(c.responseBodyKey[:], tmp[:16])
	if c.metadata.Cipher == "" {
		ciph, err := ParseCipherFromSecurity(instructionData[35] & 0xf)
		if err != nil {
			return err
		}
		c.metadata.Cipher = string(ciph)
	}
	newAEAD, ok := NewCipherMapper[Cipher(c.metadata.Cipher)]
	if !ok {
		return fmt.Errorf("unexpected cipher: %v", c.metadata.Cipher)
	}
	c.NewAEAD = newAEAD
	c.requestOptions = instructionData[34]
	return nil
}

func (c *Conn) WriteReqHeader() (err error) {
	c.initWrite.Do(func() {
		instructionData := ReqInstructionDataFromPool(c.metadata)
		defer pool.Put(instructionData)

		if err = c.InitContext(instructionData); err != nil {
			return
		}

		var header []byte
		if header, err = EncryptReqHeaderFromPool(instructionData, c.cmdKey); err != nil {
			return
		}
		defer pool.Put(header)
		if c.writeBodyCipher, err = c.NewAEAD(c.requestBodyKey[:]); err != nil {
			return
		}

		if ContainOption(c.requestOptions, OptionChunkLengthMasking) {
			c.writeChunkSizeParser = NewShakeSizeParser(c.requestBodyIV[:])
			if ContainOption(c.requestOptions, OptionGlobalPadding) {
				c.writePaddingGenerator = c.writeChunkSizeParser.(PaddingLengthGenerator)
			}
		} else {
			c.writeChunkSizeParser = PlainChunkSizeParser{}
		}
		if c.writePaddingGenerator == nil {
			c.writePaddingGenerator = PlainPaddingGenerator{}
		}
		c.writeNonceGenerator = GenerateChunkNonce(c.requestBodyIV[:], uint32(c.writeBodyCipher.NonceSize()))
		_, err = c.Conn.Write(header)
	})
	return err
}

// Write writes data to the connection. Empty b should be written before closing the connection to indicate the terminal.
func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	var encRespHeader []byte
	c.initWrite.Do(func() {
		if !c.metadata.IsClient {
			header := RespHeaderFromPool(c.responseAuth)
			defer pool.Put(header)
			encRespHeader, err = c.EncryptRespHeaderFromPool(header)
			if err != nil {
				return
			}
			if c.writeBodyCipher, err = c.NewAEAD(c.responseBodyKey[:]); err != nil {
				return
			}
			if ContainOption(c.requestOptions, OptionChunkLengthMasking) {
				c.writeChunkSizeParser = NewShakeSizeParser(c.responseBodyIV[:])

				if ContainOption(c.requestOptions, OptionGlobalPadding) {
					c.writePaddingGenerator = c.writeChunkSizeParser.(PaddingLengthGenerator)
				}
			} else {
				c.writeChunkSizeParser = PlainChunkSizeParser{}
			}
			if c.writePaddingGenerator == nil {
				c.writePaddingGenerator = PlainPaddingGenerator{}
			}
			c.writeNonceGenerator = GenerateChunkNonce(c.responseBodyIV[:], uint32(c.writeBodyCipher.NonceSize()))
		}
	})
	if len(encRespHeader) != 0 {
		defer pool.Put(encRespHeader)
	}
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		data := c.sealFromPool(nil)
		defer pool.Put(data)
		_, err = c.Conn.Write(data)
		return 0, err
	}
	//log.Trace("vmess: write len(b)=%v", len(b))
	switch c.metadata.Network {
	case "tcp":
		return c.writeStream(b, encRespHeader)
	case "udp":
		return c.writePacket(b, encRespHeader)
	default:
		return 0, fmt.Errorf("unsupported network (instruction cmd): %v", c.metadata.Network)
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()
	c.initRead.Do(func() {
		if c.metadata.IsClient {
			bufSize := pool.Get(18) // 2+16
			defer pool.Put(bufSize)
			if _, err = io.ReadFull(c.Conn, bufSize); err != nil {
				err = fmt.Errorf("failed to read response header length: %w", err)
				return
			}
			var ciph cipher.AEAD
			if ciph, err = NewAesGcm(KDF(c.responseBodyKey[:], []byte(KDFSaltConstAEADRespHeaderLenKey))[:16]); err != nil {
				return
			}
			if _, err = ciph.Open(bufSize[:0], KDF(c.responseBodyIV[:], []byte(KDFSaltConstAEADRespHeaderLenIV))[:12], bufSize, nil); err != nil {
				err = fmt.Errorf("failed to decrypt response header length: %w", err)
				return
			}
			headerSize := binary.BigEndian.Uint16(bufSize[:2])
			buf := pool.Get(int(headerSize) + 16)
			defer pool.Put(buf)
			if _, err = io.ReadFull(c.Conn, buf); err != nil {
				err = fmt.Errorf("failed to read response header: %w", err)
				return
			}
			if ciph, err = NewAesGcm(KDF(c.responseBodyKey[:], []byte(KDFSaltConstAEADRespHeaderPayloadKey))[:16]); err != nil {
				return
			}
			if _, err = ciph.Open(buf[:0], KDF(c.responseBodyIV[:], []byte(KDFSaltConstAEADRespHeaderPayloadIV))[:12], buf, nil); err != nil {
				err = fmt.Errorf("failed to decrypt response header: %w", err)
				return
			}
			if buf[0] != c.responseAuth {
				err = fmt.Errorf("unexpected response auth: %v, expect %v", buf[0], c.responseAuth)
				return
			}
			respCmd := buf[2]
			if respCmd != 0 {
				err = fmt.Errorf("unexpected response command: %v", respCmd)
				return
			}
			if c.readBodyCipher, err = c.NewAEAD(c.responseBodyKey[:]); err != nil {
				return
			}

			if ContainOption(c.requestOptions, OptionChunkLengthMasking) {
				c.readChunkSizeParser = NewShakeSizeParser(c.responseBodyIV[:])

				if ContainOption(c.requestOptions, OptionGlobalPadding) {
					c.readPaddingGenerator = c.readChunkSizeParser.(PaddingLengthGenerator)
				}
			} else {
				c.readChunkSizeParser = PlainChunkSizeParser{}
			}
			if c.readPaddingGenerator == nil {
				c.readPaddingGenerator = PlainPaddingGenerator{}
			}
			c.readNonceGenerator = GenerateChunkNonce(c.responseBodyIV[:], uint32(c.readBodyCipher.NonceSize()))
		} else {
			// assume that EAuthID has been read
			buf := pool.Get(26) // len(2) + tag(16) + connection_nonce(8)
			defer pool.Put(buf)
			if _, err = io.ReadFull(c.Conn, buf); err != nil {
				err = fmt.Errorf("failed to read ALength and ConnectionNonce: %w", err)
				return
			}
			connectionNonce := buf[18:26]
			c.cmdKey = c.metadata.authedCmdKey[:]
			var ciph cipher.AEAD
			if ciph, err = NewAesGcm(KDF(c.cmdKey, []byte(KDFSaltConstVMessHeaderPayloadLengthAEADKey), c.metadata.authedEAuthID[:], connectionNonce)[:16]); err != nil {
				return
			}
			if _, err = ciph.Open(buf[:0], KDF(c.cmdKey, []byte(KDFSaltConstVMessHeaderPayloadLengthAEADIV), c.metadata.authedEAuthID[:], connectionNonce)[:12], buf[:18], c.metadata.authedEAuthID[:]); err != nil {
				err = fmt.Errorf("failed to decrypt request header length: %w", err)
				return
			}
			lenInstruction := binary.BigEndian.Uint16(buf)

			instructionData := pool.Get(int(lenInstruction) + 16)
			defer pool.Put(instructionData)
			if _, err = io.ReadFull(c.Conn, instructionData); err != nil {
				err = fmt.Errorf("failed to read instruction data: %w", err)
				return
			}
			if ciph, err = NewAesGcm(KDF(c.cmdKey, []byte(KDFSaltConstVMessHeaderPayloadAEADKey), c.metadata.authedEAuthID[:], connectionNonce)[:16]); err != nil {
				return
			}
			if _, err = ciph.Open(instructionData[:0], KDF(c.cmdKey, []byte(KDFSaltConstVMessHeaderPayloadAEADIV), c.metadata.authedEAuthID[:], connectionNonce)[:12], instructionData, c.metadata.authedEAuthID[:]); err != nil {
				err = fmt.Errorf("failed to decrypt request header: %w", err)
				return
			}
			if err = c.InitContext(instructionData[:lenInstruction]); err != nil {
				return
			}
			if err = c.metadata.CompleteFromInstructionData(instructionData[:lenInstruction]); err != nil {
				return
			}

			if c.readBodyCipher, err = c.NewAEAD(c.requestBodyKey[:]); err != nil {
				return
			}
			if ContainOption(c.requestOptions, OptionChunkLengthMasking) {
				c.readChunkSizeParser = NewShakeSizeParser(c.requestBodyIV[:])

				if ContainOption(c.requestOptions, OptionGlobalPadding) {
					c.readPaddingGenerator = c.readChunkSizeParser.(PaddingLengthGenerator)
				}
			} else {
				c.readChunkSizeParser = PlainChunkSizeParser{}
			}
			if c.readPaddingGenerator == nil {
				c.readPaddingGenerator = PlainPaddingGenerator{}
			}
			c.readNonceGenerator = GenerateChunkNonce(c.requestBodyIV[:], uint32(c.readBodyCipher.NonceSize()))
		}
	})
	if err != nil {
		return 0, err
	}
	if b == nil {
		return 0, nil
	}
	if c.readNonceGenerator == nil {
		// did not initiate successfully
		return 0, net.ErrClosed
	}

	// dump unread data
	if c.indexToRead < len(c.leftToRead) {
		n = copy(b, c.leftToRead[c.indexToRead:])
		c.indexToRead += n
		if c.indexToRead >= len(c.leftToRead) {
			// put the buf back
			pool.Put(c.leftToRead)
		}
		return n, nil
	}

	chunk, err := c.readChunkFromPool()
	if err != nil {
		return 0, err
	}
	//log.Trace("vmess: read len(chunk)=%v", len(chunk))
	n = copy(b, chunk)
	if n < len(chunk) {
		// wait for the next read
		c.leftToRead = chunk
		c.indexToRead = n
	} else {
		// full reading. put the buf back
		pool.Put(chunk)
	}
	return n, nil
}

func (c *Conn) Metadata() Metadata {
	return c.metadata
}

// readSize reads the size and padding from Conn. size=encryptedSize+padding
func (c *Conn) readSize() (size uint16, padding uint16, err error) {
	buf := pool.Get(int(c.readChunkSizeParser.SizeBytes()))
	defer pool.Put(buf)
	if _, err := io.ReadFull(c.Conn, buf); err != nil {
		return 0, 0, err
	}
	padding = c.readPaddingGenerator.NextPaddingLen()
	size, err = c.readChunkSizeParser.Decode(buf)
	if err != nil {
		return size, padding, err
	}
	//log.Warn("read: size: %v, padding: %v", size, padding)
	return size, padding, nil
}

func (c *Conn) readChunkFromPool() (b []byte, err error) {
	size, padding, err := c.readSize()
	if err != nil {
		return nil, err
	}
	// terminal signal
	if size == uint16(c.readBodyCipher.Overhead())+padding {
		return nil, io.EOF
	}
	b = pool.Get(int(size))
	if _, err = io.ReadFull(c.Conn, b); err != nil {
		pool.Put(b)
		return nil, err
	}
	return c.readBodyCipher.Open(b[:0], c.readNonceGenerator(), b[:len(b)-int(padding)], nil)
}

func (c *Conn) EncryptRespHeaderFromPool(header []byte) (b []byte, err error) {
	buf := pool.Get(34 + len(header)) // length(2) + tag(16) + len(header) + tag(16)

	ciph, err := NewAesGcm(KDF(c.responseBodyKey[:], []byte(KDFSaltConstAEADRespHeaderLenKey))[:16])
	if err != nil {
		pool.Put(buf)
		return
	}
	binary.BigEndian.PutUint16(buf, uint16(len(header)))
	ciph.Seal(buf[:0], KDF(c.responseBodyIV[:], []byte(KDFSaltConstAEADRespHeaderLenIV))[:12], buf[:2], nil)

	ciph, err = NewAesGcm(KDF(c.responseBodyKey[:], []byte(KDFSaltConstAEADRespHeaderPayloadKey))[:16])
	if err != nil {
		pool.Put(buf)
		return
	}
	ciph.Seal(buf[18:18], KDF(c.responseBodyIV[:], []byte(KDFSaltConstAEADRespHeaderPayloadIV))[:12], header, nil)

	return buf, nil
}
