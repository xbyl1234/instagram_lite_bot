// Pulled from https://github.com/youtube/vitess 229422035ca0c716ad0c1397ea1351fe62b0d35a
// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package czlib

import (
	"bytes"
	"io"
)

// err starts out as nil
// we will call inflateEnd when we set err to a value:
// - whatever error is returned by the underlying ZipReader
// - io.EOF if Close was called
type ZipReader struct {
	r      *bytes.Buffer
	in     []byte
	strm   zstream
	err    error
	skipIn bool
}

// NewReader creates a new io.ReadCloser. Reads from the returned io.ReadCloser
// read and decompress data from r. The implementation buffers input and may read
// more data than necessary from r.
// It is the caller's responsibility to call Close on the ReadCloser when done.
func NewReader() (*ZipReader, error) {
	return NewReaderBuffer(bytes.NewBuffer([]byte{}), DEFAULT_COMPRESSED_BUFFER_SIZE)
}

// NewReaderBuffer has the same behavior as NewReader but the user can provides
// a custom buffer size.
func NewReaderBuffer(r *bytes.Buffer, bufferSize int) (*ZipReader, error) {
	z := &ZipReader{r: r, in: make([]byte, bufferSize)}
	if err := z.strm.inflateInit(); err != nil {
		return nil, err
	}
	return z, nil
}

func (z *ZipReader) SetInput(data []byte) {
	z.err = nil
	z.r.Write(data)
}

func (z *ZipReader) Inflate() ([]byte, error) {
	reader := bytes.NewBuffer([]byte{})
	_, err := io.Copy(reader, z)
	if err != nil {
		return nil, err
	}
	return reader.Bytes(), nil
}

func (z *ZipReader) Read(p []byte) (int, error) {
	if z.err != nil {
		return 0, z.err
	}

	if len(p) == 0 {
		return 0, nil
	}

	// read and deflate until the output buffer is full
	z.strm.setOutBuf(p, len(p))

	for {
		// if we have no data to inflate, read more
		if !z.skipIn && z.strm.availIn() == 0 {
			var n int
			n, z.err = z.r.Read(z.in)
			// If we got data and EOF, pretend we didn't get the
			// EOF.  That way we will return the right values
			// upstream.  Note this will trigger another read
			// later on, that should return (0, EOF).
			if n > 0 && z.err == io.EOF {
				z.err = nil
			}

			// FIXME(alainjobart) this code is not compliant with
			// the Reader interface. We should process all the
			// data we got from the ZipReader, and then return the
			// error, whatever it is.
			if (z.err != nil && z.err != io.EOF) || (n == 0 && z.err == io.EOF) {
				return 0, z.err
			}

			z.strm.setInBuf(z.in, n)
		} else {
			z.skipIn = false
		}

		// inflate some
		ret, err := z.strm.inflate(zNoFlush)
		if err != nil {
			z.err = err
			z.strm.inflateEnd()
			return 0, z.err
		}

		// if we read something, we're good
		have := len(p) - z.strm.availOut()
		if have > 0 {
			z.skipIn = ret == Z_OK && z.strm.availOut() == 0
			return have, z.err
		}
	}
}

// Close closes the Reader. It does not close the underlying io.Reader.
func (z *ZipReader) Close() error {
	if z.err != nil {
		if z.err != io.EOF {
			return z.err
		}
		return nil
	}
	z.strm.inflateEnd()
	z.err = io.EOF
	return nil
}
