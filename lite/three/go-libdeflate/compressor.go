package libdeflate

import "CentralizedControl/three/go-libdeflate/native"

// Compressor compresses data at the specified compression level.
//
// A single compressor must not not be used across multiple threads concurrently.
// If you want to compress concurrently, create a compressor for each thread.
//
// Always Close() the decompressor to free c memory.
// One Compressor allocates at least 32 KiB.
type Compressor struct {
	c   *native.Compressor
	lvl int
}

// NewCompressor returns a new Compressor used to compress data with compression level DefaultCompressionLevel.
// Errors if out of memory. Allocates 32KiB.
// See NewCompressorLevel for custom compression level
func NewCompressor() (Compressor, error) {
	return NewCompressorLevel(DefaultCompressionLevel)
}

// NewCompressorLevel returns a new Compressor used to compress data.
// Errors if out of memory or if an invalid compression level was passed.
// Allocates 32KiB.
//
// The compression level is legal if and only if:
// MinCompressionLevel <= level <= MaxCompressionLevel
func NewCompressorLevel(level int) (Compressor, error) {
	c, err := native.NewCompressor(level)
	return Compressor{c, level}, err
}

// CompressZlib compresses the data from in to out (in zlib format) and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though).
//
// See c.Compress for further information.
func (c Compressor) CompressZlib(in, out []byte) (int, []byte, error) {
	return c.Compress(in, out, ModeZlib)
}

// Compress compresses the data from in to out and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though).
//
// m specifies which compression format should be used (e.g. ModeZlib)
//
// Notice that for extremely small or already highly compressed data,
// the compressed data could be larger than uncompressed.
// If out == nil: For a too large discrepancy (len(out) > 1000 + 2 * len(in)) Compress will error
func (c Compressor) Compress(in, out []byte, m Mode) (int, []byte, error) {
	switch m {
	case ModeZlib:
		return c.c.Compress(in, out, native.CompressZlib)
	case ModeDEFLATE:
		return c.c.Compress(in, out, native.CompressDEFLATE)
	case ModeGzip:
		return c.c.Compress(in, out, native.CompressGzip)
	default:
		panic(errorInvalidModeCompressor)
	}
}

// Level returns the compression level at which this Compressor compresses.
// May be called after having closed a Compressor.
func (c Compressor) Level() int {
	return c.lvl
}

/*
WorstCaseCompressedSize returns the maximum theoretical size of the data after compressing data of length 'size',
using the given mode of compression.
This prediction is a wild overestimate in most cases, for which holds true: max >= size.
However, it gives a hard maximal bound of the size of compressed data, compressing with the given mode
at the compression level of the this compressor, independent of the actual data.
This method will always return the same max size for the same compressor, input size and mode.
*/
func (c Compressor) WorstCaseCompressedSize(size int, m Mode) (max int) {
	switch m {
	case ModeDEFLATE:
		return c.c.UpperBound(size, native.DeflateBound)
	case ModeZlib:
		return c.c.UpperBound(size, native.ZlibBound)
	case ModeGzip:
		return c.c.UpperBound(size, native.GzipBound)
	default:
		panic(errorInvalidModeCompressor)
	}
}

// Close closes the compressor and releases all occupied resources.
// It is the users responsibility to close compressors in order to free resources,
// as the underlying c objects are not subject to the go garbage collector. They have to be freed manually.
//
// After closing, the compressor must not be used anymore, as the methods will panic (except for the c.Level() method).
func (c Compressor) Close() {
	c.c.Close()
}
