package native

/*
#include "libdeflate.h"

typedef struct libdeflate_compressor comp;
*/
import "C"

type bound func(*C.comp, int) int

// DeflateBound works as described in native/libs/libdeflate.h
func DeflateBound(c *C.comp, s int) int {
	return int(C.libdeflate_deflate_compress_bound(c, intToInt64(s)))
}

// ZlibBound works as described in native/libs/libdeflate.h
func ZlibBound(c *C.comp, s int) int {
	return int(C.libdeflate_zlib_compress_bound(c, intToInt64(s)))
}

// GzipBound works as described in native/libs/libdeflate.h
func GzipBound(c *C.comp, s int) int {
	return int(C.libdeflate_gzip_compress_bound(c, intToInt64(s)))
}
