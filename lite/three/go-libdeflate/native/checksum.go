package native

/*
#include "libdeflate.h"
*/
import "C"
import "unsafe"

/*
Adler32 works as described in native/libs/libdeflate.h.
Returns 1 for in == nil
*/
func Adler32(adler32 uint32, in []byte) uint32 {
	if in == nil {
		return 1
	}

	addr := startMemAddr(in)

	checksum := C.libdeflate_adler32(
		toUInt32(adler32),
		unsafe.Pointer(addr),
		toInt64(int64(len(in))),
	)

	return uint32(checksum)
}

/*
Crc32 works as described in native/libs/libdeflate.h.
Returns 0 for in == nil
*/
func Crc32(crc32 uint32, in []byte) uint32 {
	if in == nil {
		return 0
	}

	addr := startMemAddr(in)

	checksum := C.libdeflate_crc32(
		toUInt32(crc32),
		unsafe.Pointer(addr),
		toInt64(int64(len(in))),
	)

	return uint32(checksum)
}
