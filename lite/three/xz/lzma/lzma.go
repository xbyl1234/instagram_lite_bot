// Package lzma is a thin wrapper around the C lzma2 library.
//
// The emphasis is on the word "thin". This package does not provide an
// idiomatic Go API; rather, it simply wraps C functions and types with
// analogous Go functions and types.
// A nice Go API should be built on top of this package.
//
// The documentation for each type and function in this package generally just
// contains a reference to
// to the underlying C type or function in the /src/liblzma/api/ directory of the
// upstream C repository. Full documentation for the type and function can be found
// by looking at the excellent documentation on the C side.
package lzma

/*
#cgo CFLAGS: -Isrc/common
#cgo CFLAGS: -Isrc/liblzma/api
#cgo CFLAGS: -Isrc/liblzma/common
#cgo CFLAGS: -Isrc/liblzma/check
#cgo CFLAGS: -Isrc/liblzma/delta
#cgo CFLAGS: -Isrc/liblzma/lz
#cgo CFLAGS: -Isrc/liblzma/lzma
#cgo CFLAGS: -Isrc/liblzma/rangecoder
#cgo CFLAGS: -Isrc/liblzma/simple

// This block of flags specify which lzma2 features to link in.
#cgo CFLAGS: -DHAVE_ENCODER_LZMA2 -DHAVE_DECODER_LZMA2
#cgo CFLAGS: -DHAVE_CHECK_CRC32 -DHAVE_CHECK_CRC64
// The following 3 flags were determined by inspecting lzma_encoder_presets.c
#cgo CFLAGS: -DHAVE_MF_HC3 -DHAVE_MF_HC4 -DHAVE_MF_BT4

// We need to manually specify whether the architecture is 32-bit or 64-bit using the
// SIZEOF_SIZE_T macro. We do this by manually enumerating all 32-bit architectures
// supported by Go, marking those as 32-bit, and assuming all others are 64-bit.
// This has the caveat that if support for a new 32-bit architecture is added the
// package will not work for that architecture until it is added here. However, a new
// 32-bit architecture is unlikely. The canonical list of architectures supported by
// Go is here:
// https://github.com/golang/go/blob/master/src/go/build/syslist.go
#cgo  386  amd64p32  arm  armbe  mips  mipsle  mips64p32  mips64p32le  ppc  riscv  s390  sparc CFLAGS: -DSIZEOF_SIZE_T=4
#cgo !386,!amd64p32,!arm,!armbe,!mips,!mipsle,!mips64p32,!mips64p32le,!ppc,!riscv,!s390,!sparc CFLAGS: -DSIZEOF_SIZE_T=8

// We assume these C standard libraries are available.
#cgo CFLAGS: -DHAVE_STDBOOL_H -DHAVE_STDINT_H -DHAVE_INTTYPES_H

// Performance improvement for 32-bit and 64-bit x86.
#cgo 386 amd64 CFLAGS: -DTUKLIB_FAST_UNALIGNED_ACCESS

#include <stdlib.h>
#include <string.h>

#include "src/liblzma/api/lzma.h"

// The lzma library requires that the stream be initialized to the value of the macro
// LZMA_STREAM_INIT. Because this is a macro it has no type. This function exists to cast the
// macro to the stream type.
lzma_stream new_stream() {
	lzma_stream strm = LZMA_STREAM_INIT;
	return strm;
}
*/
import "C"
import (
	"unsafe"
)

const (
	Ok               Return = 0
	StreamEnd               = 1
	NoCheck                 = 2
	UnsupportedCheck        = 3
	GetCheck                = 4
	MemoryError             = 5
	MemoryLimitError        = 6
	FormatError             = 7
	OptionsError            = 8
	DataError               = 9
	BufferError             = 10
	ProgrammingError        = 11
	SeekNeeded              = 12
)

const (
	Run         Action = 0
	SyncFlush          = 1
	FullFlush          = 2
	Finish             = 3
	FullBarrier        = 4
)

// Return corresponds to the lzma_ret type in base.h.
type Return int

func (r Return) String() string {
	switch r {
	case Ok:
		return "OK"
	case StreamEnd:
		return "STREAM_END"
	case NoCheck:
		return "NO_CHECK"
	case UnsupportedCheck:
		return "UNSUPPORTED_CHECK"
	case GetCheck:
		return "GET_CHECK"
	case MemoryError:
		return "MEMORY_ERROR"
	case MemoryLimitError:
		return "MEMORY_LIMIT_ERROR"
	case FormatError:
		return "FORMAT_ERROR"
	case OptionsError:
		return "OPTIONS_ERROR"
	case DataError:
		return "DATA_ERROR"
	case BufferError:
		return "BUFFER_ERROR"
	case ProgrammingError:
		return "PROGRAMMING_ERROR"
	case SeekNeeded:
		return "SEEK_NEEDED"
	}
	return "UNKNOWN_RESULT"
}

// IsErr returns true if the return code indicates an error condition.
func (r Return) IsErr() bool {
	return r != Ok && r != StreamEnd
}

// Action corresponds to the lzma_action type in base.h.
type Action int

func (a Action) String() string {
	switch a {
	case Run:
		return "RUN"
	case SyncFlush:
		return "SYNC_FLUSH"
	case FullFlush:
		return "FULL_FLUSH"
	case Finish:
		return "FINISH"
	case FullBarrier:
		return "FULL_BARRIER"
	}
	return "UNKNOWN_ACTION"
}

type cBuffer struct {
	start *C.uint8_t
	len   C.size_t
	cap   C.size_t
}

func (buf *cBuffer) set(p []byte) {
	if len(p) == 0 {
		buf.len = 0
		return
	}
	buf.grow(len(p))
	C.memcpy(unsafe.Pointer(buf.start), unsafe.Pointer(&p[0]), C.size_t(len(p)))
	buf.len = C.size_t(len(p))
}

func (buf *cBuffer) read(length int) []byte {
	return C.GoBytes(unsafe.Pointer(buf.start), C.int(length))
}

func (buf *cBuffer) grow(n int) {
	if n <= int(buf.cap) {
		return
	}
	buf.clear()
	buf.start = (*C.uint8_t)(C.malloc(C.size_t(n)))
	buf.len = 0
	buf.cap = C.size_t(n)
}

func (buf *cBuffer) clear() {
	if buf.start != nil {
		C.free(unsafe.Pointer(buf.start))
	}
	buf.start = nil
	buf.len = 0
	buf.cap = 0
}

// This was chosen arbitrarily but seems to work fine in practice
const outputBufferLength = 1024

// Stream wraps lzma_stream in base.h and the input and output buffers that the lzma_stream type
// requires to exist.
//
// The lzma_stream type operates on the two buffers but does not take ownership of them. This
// type thus contains handling for these buffers. This part of the package is the most Go-like
// because it needs to map from Go slices to C arrays, and ultimately hide the C implementation
// details.
type Stream struct {
	cStream C.lzma_stream
	input   cBuffer
	output  cBuffer
}

// NewStream returns a new stream.
func NewStream() *Stream {
	stream := Stream{
		cStream: C.new_stream(),
	}
	stream.output.grow(outputBufferLength)
	stream.output.len = outputBufferLength
	stream.cStream.next_out = stream.output.start
	stream.cStream.avail_out = stream.output.len
	return &stream
}

// AvailIn returns the number of bytes that have been placed in the input buffer using the SetInput
// method that have yet to be processed by the stream.
func (stream *Stream) AvailIn() int {
	return int(stream.cStream.avail_in)
}

// TotalIn returns the total number of bytes that have been read from the input buffer.
func (stream *Stream) TotalIn() int {
	return int(stream.cStream.total_in)
}

// AvailOut returns the number of bytes that the stream has written into the output buffer that
// have yet to be read using the Output method.
func (stream *Stream) AvailOut() int {
	return int(stream.cStream.avail_out)
}

// TotalOut returns the total number of bytes that have been written to the input buffer
func (stream *Stream) TotalOut() int {
	return int(stream.cStream.total_out)
}

// SetInput sets the input buffer of the stream to be the provided bytes. Note this overwrites
// any data that is already in the input buffer, so before calling SetInput it's best to verify
// that AvailIn returns 0.
func (stream *Stream) SetInput(p []byte) {
	stream.input.set(p)
	stream.cStream.next_in = stream.input.start
	stream.cStream.avail_in = stream.input.len
}

// Output returns all bytes that have been written to the output buffer by the stream, and resets
// the output buffer.
func (stream *Stream) Output() []byte {
	b := stream.output.read(int(stream.output.len - stream.cStream.avail_out))
	stream.cStream.next_out = stream.output.start
	stream.cStream.avail_out = stream.output.len
	return b
}

// Close closes the stream and releases C memory that has been allocated by the type.
func (stream *Stream) Close() {
	stream.input.clear()
	stream.output.clear()
	End(stream)
}

// End wraps lzma_end in base.h.
func End(stream *Stream) {
	C.lzma_end(&stream.cStream)
}

// EasyEncoder wraps lzma_easy_encoder in container.h.
func EasyEncoder(stream *Stream, preset int) Return {
	return Return(C.lzma_easy_encoder(&stream.cStream, C.uint(preset), C.LZMA_CHECK_CRC64))
}

// StreamDecoder wraps lzma_stream_decoder in container.h.
func StreamDecoder(stream *Stream) Return {
	return Return(C.lzma_stream_decoder(&stream.cStream, C.UINT64_MAX, C.LZMA_TELL_UNSUPPORTED_CHECK))
}

// Code wraps lzma_code in base.h.
func Code(stream *Stream, action Action) Return {
	return Return(C.lzma_code(&stream.cStream, C.lzma_action(action)))
}
