package native

/*
#include "libdeflate.h"
#include "helper.h"
#include <stddef.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct libdeflate_decompressor decomp;
*/
import "C"
import "unsafe"

// Decompressor decompresses any DEFLATE, zlib or gzip compressed data at any level
type Decompressor struct {
	dc                     *C.decomp
	isClosed               bool
	maxDecompressionFactor int
}

// NewDecompressor returns a new Decompressor with maxDecompressionFactor = 30 or and error if out of memory
func NewDecompressor() (*Decompressor, error) {
	return NewDecompressorWithExtendedDecompression(30)
}

// NewDecompressorWithExtendedDecompression returns a new Decompressor with maxDecompressionFactor or and error if out of memory
func NewDecompressorWithExtendedDecompression(maxDecompressionFactor int) (*Decompressor, error) {
	dc := C.libdeflate_alloc_decompressor()
	if C.isNull(unsafe.Pointer(dc)) == 1 {
		return nil, errorOutOfMemory
	}

	return &Decompressor{dc, false, maxDecompressionFactor}, nil
}

// Decompress decompresses the given data from in to out and returns out and an error if something went wrong.
// If error != nil, then the data in out is undefined.
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil as out, this function will allocate a sufficient buffer and return it.
// Returns the number of consumed bytes from 'in'
func (dc *Decompressor) Decompress(in, out []byte, f decompress) (int, []byte, error) {
	if dc.isClosed {
		panic(errorAlreadyClosed)
	}
	if len(in) == 0 {
		return 0, out, errorNoInput
	}

	if out != nil {
		cons, _, err := dc.decompress(in, out, true, f)
		return cons, out, err
	}

	cons := 0
	n := 0
	decompFactor := 6
	err := errorInsufficientSpace
	for err == errorInsufficientSpace {
		out = make([]byte, len(in)*decompFactor)
		cons, n, err = dc.decompress(in, out, false, f)

		if decompFactor > dc.maxDecompressionFactor {
			return cons, out, errorInsufficientDecompressionFactor
		}

		if decompFactor >= 16 {
			decompFactor += 3
			continue
		}
		decompFactor += 5
	}

	return cons, out[:n], err
}

func (dc *Decompressor) decompress(in, out []byte, fit bool, f decompress) (int, int, error) {
	inAddr := startMemAddr(in)
	outAddr := startMemAddr(out)

	var (
		cons int
		n    int
	)

	consPtr := uintptr(unsafe.Pointer(&cons))
	sPtr := uintptr(unsafe.Pointer(&n))
	if fit {
		sPtr = 0
	}

	err := f(dc.dc, inAddr, outAddr, len(in), len(out), consPtr, sPtr)

	if fit {
		n = len(out)
	}

	return cons, n, err
}

// Close frees the memory allocated by C objects
func (dc *Decompressor) Close() {
	if dc.isClosed {
		panic(errorAlreadyClosed)
	}
	C.libdeflate_free_decompressor(dc.dc)
	dc.isClosed = true
}
