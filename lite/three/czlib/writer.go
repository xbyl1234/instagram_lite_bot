// Pulled from https://github.com/youtube/vitess 229422035ca0c716ad0c1397ea1351fe62b0d35a
// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package czlib

import (
	"bytes"
	"fmt"
	"io"
)

// Allowed flush values
const (
	Z_NO_FLUSH      = 0
	Z_PARTIAL_FLUSH = 1
	Z_SYNC_FLUSH    = 2
	Z_FULL_FLUSH    = 3
	Z_FINISH        = 4
	Z_BLOCK         = 5
	Z_TREES         = 6
)

// Return codes
const (
	Z_OK            = 0
	Z_STREAM_END    = 1
	Z_NEED_DICT     = 2
	Z_ERRNO         = -1
	Z_STREAM_ERROR  = -2
	Z_DATA_ERROR    = -3
	Z_MEM_ERROR     = -4
	Z_BUF_ERROR     = -5
	Z_VERSION_ERROR = -6
)

// our default buffer size
// most go io functions use 32KB as buffer size, so 32KB
// works well here for compressed data buffer
const (
	DEFAULT_COMPRESSED_BUFFER_SIZE = 32 * 1024
)

// ZipWriter implements a io.WriteCloser
// we will call deflateEnd when we set err to a value:
// - whatever error is returned by the underlying writer
// - io.EOF if Close was called
type ZipWriter struct {
	w    *bytes.Buffer
	out  []byte
	strm zstream
	err  error
}

// NewWriter returns a new zlib writer that writes to the underlying writer
func NewWriter() *ZipWriter {
	z, _ := NewWriterLevelBuffer(bytes.NewBuffer([]byte{}), DefaultCompression, DEFAULT_COMPRESSED_BUFFER_SIZE)
	return z
}

// NewWriterLevel let the user provide a compression level value
func NewWriterLevel(level int) (*ZipWriter, error) {
	return NewWriterLevelBuffer(bytes.NewBuffer([]byte{}), level, DEFAULT_COMPRESSED_BUFFER_SIZE)
}

// NewWriterLevelBuffer let the user provide compression level and buffer size values
func NewWriterLevelBuffer(w *bytes.Buffer, level, bufferSize int) (*ZipWriter, error) {
	z := &ZipWriter{w: w, out: make([]byte, bufferSize)}
	if err := z.strm.deflateInit(level); err != nil {
		return nil, err
	}
	return z, nil
}

// this is the main function: it advances the write with either
// new data or something else to do, like a flush
func (z *ZipWriter) write(p []byte, flush int) int {
	if len(p) == 0 {
		z.strm.setInBuf(nil, 0)
	} else {
		z.strm.setInBuf(p, len(p))
	}
	// we loop until we don't get a full output buffer
	// each loop completely writes the output buffer to the underlying
	// writer
	for {
		// deflate one buffer
		z.strm.setOutBuf(z.out, len(z.out))
		z.strm.deflate(flush)

		// write everything
		from := 0
		have := len(z.out) - int(z.strm.availOut())
		for have > 0 {
			var n int
			n, z.err = z.w.Write(z.out[from:have])
			if z.err != nil {
				z.strm.deflateEnd()
				return 0
			}
			from += n
			have -= n
		}

		// we stop trying if we get a partial response
		if z.strm.availOut() != 0 {
			break
		}
	}
	// the library guarantees this
	if z.strm.availIn() != 0 {
		panic(fmt.Errorf("cgzip: Unexpected error (2)"))
	}
	return len(p)
}

// Write implements the io.ZipWriter interface
func (z *ZipWriter) Write(p []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}
	n = z.write(p, Z_NO_FLUSH)
	return n, z.err
}

// Flush let the user flush the zlib buffer to the underlying writer buffer
func (z *ZipWriter) Flush() error {
	if z.err != nil {
		return z.err
	}
	z.write(nil, Z_SYNC_FLUSH)
	return z.err
}

func (z *ZipWriter) Read() ([]byte, error) {
	return io.ReadAll(z.w)
}

// Close closes the zlib buffer but does not close the wrapped io.ZipWriter originally
// passed to NewWriterX.
func (z *ZipWriter) Close() error {
	if z.err != nil {
		return z.err
	}
	z.write(nil, Z_FINISH)
	if z.err != nil {
		return z.err
	}
	z.strm.deflateEnd()
	z.err = io.EOF
	return nil
}
