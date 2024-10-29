package native

/*
#include "libdeflate.h"
#include "helper.h"

typedef struct libdeflate_compressor comp;
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Compressor compresses data to zlib format at the specified level
type Compressor struct {
	c        *C.comp
	isClosed bool
	lvl      int
}

// NewCompressor returns a new Compressor used to compress data.
// Errors if out of memory or invalid lvl
func NewCompressor(lvl int) (*Compressor, error) {
	if lvl < MinCompressionLevel || lvl > MaxCompressionLevel {
		return nil, errorInvalidLevel
	}

	c := C.libdeflate_alloc_compressor(C.int(lvl))
	if C.isNull(unsafe.Pointer(c)) == 1 {
		return nil, errorOutOfMemory
	}

	return &Compressor{c, false, lvl}, nil
}

// Compress compresses the data from in to out and returns the number
// of bytes written to out, out and an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it.
//
// Notice that for extremely small or already highly compressed data,
// the compressed data could be larger than uncompressed.
// If out == nil: For a too large discrepancy (len(out) > 1000 + 2 * len(in)) Compress will error
func (c *Compressor) Compress(in, out []byte, f compress) (int, []byte, error) {
	if c.isClosed {
		panic(errorAlreadyClosed)
	}
	if len(in) == 0 {
		return 0, out, errorNoInput
	}

	if out != nil {
		n, b, err := c.compress(in, out, f)
		return n, b[:n], err
	}

	out = make([]byte, len(in))
	n, out, err := c.compress(in, out, f)

	if err == errorShortBuffer { // if still doesn't fit (shouldn't happen at all)
		out = make([]byte, 1000+len(in)*2)
		n, _, _ := c.compress(in, out, f)
		return n, out[:n], errors.New("libdeflate: native: compressed data is much larger than uncompressed")
	}

	return n, out[:n], nil
}

func (c *Compressor) compress(in, out []byte, f compress) (int, []byte, error) {
	inAddr := startMemAddr(in)
	outAddr := startMemAddr(out)

	written := f(c.c, inAddr, outAddr, len(in), len(out))

	if written == 0 {
		return written, out, errorShortBuffer
	}
	return written, out, nil
}

// UpperBound works as described in native/libs/libdeflate.h
func (c *Compressor) UpperBound(size int, f bound) int {
	return f(c.c, size)
}

// Close frees the memory allocated by C objects
func (c *Compressor) Close() {
	if c.isClosed {
		panic(errorAlreadyClosed)
	}
	C.libdeflate_free_compressor(c.c)
	c.isClosed = true
}
