package libdeflate

import "CentralizedControl/three/go-libdeflate/native"

// Decompressor decompresses any DEFLATE, zlib or gzip compressed data at any level
//
// A single decompressor must not not be used across multiple threads concurrently.
// If you want to decompress concurrently, create a decompressor for each thread.
//
// Always Close() the decompressor to free c memory.
// One Decompressor allocates at least 32KiB.
type Decompressor struct {
	dc *native.Decompressor
}

// NewDecompressor returns a new Decompressor used to decompress data at any compression level and with any Mode.
// Errors if out of memory. Allocates 32KiB.
func NewDecompressor() (Decompressor, error) {
	dc, err := native.NewDecompressor()
	return Decompressor{dc}, err
}

// NewDecompressorWithExtendedDecompression returns a new Decompressor used to decompress data at any compression level and with any Mode.
// maxDecompressionFactor customizes how much larger your output than your input may be. This is set to a sensible default,
// however, it might need some tweaking if you have a huge compression factor. Usually, NewDecompressor should suffice.
// Errors if out of memory. Allocates 32KiB.
func NewDecompressorWithExtendedDecompression(maxDecompressionFactor int) (Decompressor, error) {
	dc, err := native.NewDecompressorWithExtendedDecompression(maxDecompressionFactor)
	return Decompressor{dc}, err
}

// DecompressZlib decompresses the given zlib data from in to out and returns the number of consumed bytes c
// from 'in' and 'out' or an error if something went wrong.
//
// c is the number of bytes that were read before the BFINAL flag was
// encountered, which indicates the end of the compressed data.
//
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil to out, this function will allocate a sufficient buffer and return it.
//
// If error != nil, the data in out is undefined.
func (dc Decompressor) DecompressZlib(in, out []byte) (int, []byte, error) {
	return dc.Decompress(in, out, ModeZlib)
}

// Decompress decompresses the given data from in to out and returns the number of consumed bytes c from 'in' and 'out'
// or an error if something went wrong.
// Mode m specifies the format (e.g. zlib) of the data within in.
//
// c is the number of bytes that were read before the BFINAL flag was
// encountered, which indicates the end of the compressed data.
//
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil to out, this function will allocate a sufficient buffer and return it.
//
// If error != nil, the data in out is undefined.
func (dc Decompressor) Decompress(in, out []byte, m Mode) (int, []byte, error) {
	switch m {
	case ModeZlib:
		return dc.dc.Decompress(in, out, native.DecompressZlib)
	case ModeDEFLATE:
		return dc.dc.Decompress(in, out, native.DecompressDEFLATE)
	case ModeGzip:
		return dc.dc.Decompress(in, out, native.DecompressGzip)
	default:
		panic(errorInvalidModeDecompressor)
	}
}

// Close closes the decompressor and releases all occupied resources.
// It is the users responsibility to close decompressors in order to free resources,
// as the underlying c objects are not subject to the go garbage collector. They have to be freed manually.
//
// After closing, the decompressor must not be used anymore, as the methods will panic.
func (dc Decompressor) Close() {
	dc.dc.Close()
}
