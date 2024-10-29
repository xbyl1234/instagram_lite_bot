package native

/*
#include "libdeflate.h"

typedef struct libdeflate_decompressor decomp;
typedef enum libdeflate_result res;

size_t* mkPtr(size_t s) {
	return (size_t*) s;
}
*/
import "C"
import "unsafe"

type decompress func(dc *C.decomp, inAddr, outAddr *byte, inSize, outSize int, consPtr, sPtr uintptr) error

// DecompressZlib interfaces with c libdeflate for zlib decompression
func DecompressZlib(dc *C.decomp, inAddr, outAddr *byte, inSize, outSize int, consPtr, sPtr uintptr) error {
	return parseResult(C.res(C.libdeflate_zlib_decompress_ex(dc,
		unsafe.Pointer(inAddr), intToInt64(inSize),
		unsafe.Pointer(outAddr), intToInt64(outSize),
		C.mkPtr(C.size_t(consPtr)),
		C.mkPtr(C.size_t(sPtr)),
	)))
}

// DecompressDEFLATE interfaces with c libdeflate for DEFLATE decompression
func DecompressDEFLATE(dc *C.decomp, inAddr, outAddr *byte, inSize, outSize int, consPtr, sPtr uintptr) error {
	return parseResult(C.res(C.libdeflate_deflate_decompress_ex(dc,
		unsafe.Pointer(inAddr), intToInt64(inSize),
		unsafe.Pointer(outAddr), intToInt64(outSize),
		C.mkPtr(C.size_t(consPtr)),
		C.mkPtr(C.size_t(sPtr)),
	)))
}

// DecompressGzip interfaces with c libdeflate for gzip decompression
func DecompressGzip(dc *C.decomp, inAddr, outAddr *byte, inSize, outSize int, consPtr, sPtr uintptr) error {
	return parseResult(C.res(C.libdeflate_gzip_decompress_ex(dc,
		unsafe.Pointer(inAddr), intToInt64(inSize),
		unsafe.Pointer(outAddr), intToInt64(outSize),
		C.mkPtr(C.size_t(consPtr)),
		C.mkPtr(C.size_t(sPtr)),
	)))
}

func parseResult(r C.res) error {
	switch r {
	case C.LIBDEFLATE_SUCCESS:
		return nil
	case C.LIBDEFLATE_BAD_DATA:
		return errorBadData
	case C.LIBDEFLATE_SHORT_OUTPUT:
		return errorShortOutput
	case C.LIBDEFLATE_INSUFFICIENT_SPACE:
		return errorInsufficientSpace
	default:
		return errorUnknown
	}
}
