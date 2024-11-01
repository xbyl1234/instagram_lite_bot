// Pulled from https://github.com/youtube/vitess 229422035ca0c716ad0c1397ea1351fe62b0d35a
// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package czlib

// NOTE: the routines defined in this file are used for verification in
// czlib_test.go, but you cannot use cgo in test files, so they are
// defined here despite not being exposed.

/*
#include "zlib.h"
*/
import "C"

import (
	"hash"
	"unsafe"
)

type crc32Hash struct {
	crc C.uLong
}

// an empty buffer has an crc32 of '1' by default, so start with that
// (the go hash/crc32 does the same)
func newCrc32() hash.Hash32 {
	c := &crc32Hash{}
	c.Reset()
	return c
}

// Write implements an io.ZipWriter interface
func (a *crc32Hash) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		a.crc = C.crc32(a.crc, (*C.Bytef)(unsafe.Pointer(&p[0])), (C.uInt)(len(p)))
	}
	return len(p), nil
}

// Sum implements a hash.Hash interface
func (a *crc32Hash) Sum(b []byte) []byte {
	s := a.Sum32()
	b = append(b, byte(s>>24))
	b = append(b, byte(s>>16))
	b = append(b, byte(s>>8))
	b = append(b, byte(s))
	return b
}

// Reset resets the hash to default value
func (a *crc32Hash) Reset() {
	a.crc = C.crc32(0, (*C.Bytef)(unsafe.Pointer(nil)), 0)
}

// Size returns the (fixed) size of the hash
func (a *crc32Hash) Size() int {
	return 4
}

// BlockSize returns the (fixed) block size of the hash
func (a *crc32Hash) BlockSize() int {
	return 1
}

// Sum32 implements a hash.Hash32 interface
func (a *crc32Hash) Sum32() uint32 {
	return uint32(a.crc)
}

// helper method for partial checksums. From the zlib.h header:
//
//	Combine two CRC-32 checksums into one.  For two sequences of bytes, seq1
//
// and seq2 with lengths len1 and len2, CRC-32 checksums were calculated for
// each, crc1 and crc2.  crc32_combine() returns the CRC-32 checksum of
// seq1 and seq2 concatenated, requiring only crc1, crc2, and len2.
func crc32Combine(crc1, crc2 uint32, len2 int) uint32 {
	return uint32(C.crc32_combine(C.uLong(crc1), C.uLong(crc2), C.z_off_t(len2)))
}
