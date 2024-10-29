package libdeflate

import (
	"CentralizedControl/three/go-libdeflate/native"
)

/*
Adler32 updates the running adler32 checksum by the contents of the slice in.
This function returns the updated checksum.
A new adler32-checksum requires 1 as initial value. This value is also returned if in == nil.
*/
func Adler32(adler32 uint32, in []byte) uint32 {
	return native.Adler32(adler32, in)
}

/*
Crc32 updates the running crc32 checksum by the contents of the slice in.
This function returns the updated checksum.
A new crc32-checksum requires 0 as initial value. This value is also returned if in == nil.
*/
func Crc32(crc32 uint32, in []byte) uint32 {
	return native.Crc32(crc32, in)
}
