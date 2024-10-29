package native

/*
#include "libdeflate.h"

typedef struct libdeflate_compressor comp;
*/
import "C"
import "unsafe"

type compress func(c *C.comp, inAddr, outAddr *byte, inSize, outSize int) int

// CompressZlib interfaces with c libdeflate for zlib compression
func CompressZlib(c *C.comp, inAddr, outAddr *byte, inSize, outSize int) int {
	return int(C.libdeflate_zlib_compress(c,
		unsafe.Pointer(inAddr), intToInt64(inSize),
		unsafe.Pointer(outAddr), intToInt64(outSize),
	))
}

// CompressDEFLATE interfaces with c libdeflate for DEFLATE compression
func CompressDEFLATE(c *C.comp, inAddr, outAddr *byte, inSize, outSize int) int {
	return int(C.libdeflate_deflate_compress(c,
		unsafe.Pointer(inAddr), intToInt64(inSize),
		unsafe.Pointer(outAddr), intToInt64(outSize),
	))
}

// CompressGzip interfaces with c libdeflate for gzip compression
func CompressGzip(c *C.comp, inAddr, outAddr *byte, inSize, outSize int) int {
	return int(C.libdeflate_gzip_compress(c,
		unsafe.Pointer(inAddr), intToInt64(inSize),
		unsafe.Pointer(outAddr), intToInt64(outSize),
	))
}
