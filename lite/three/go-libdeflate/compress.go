package libdeflate

// CompressZlib compresses the data from in to out (in zlib format) and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though).
//
// IF YOU WANT TO COMPRESS MORE THAN ONCE, PLEASE REFER TO NewCompressor(),
// as this function creates a new Compressor (alloc 32KiB) which is then closed at the end of the function.
//
// See Compress for further information.
func CompressZlib(in, out []byte) (int, []byte, error) {
	return CompressZlibLevel(in, out, DefaultCompressionLevel)
}

// CompressZlibLevel compresses the data from in to out (in zlib format) at the specified level and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though).
//
// IF YOU WANT TO COMPRESS MORE THAN ONCE, PLEASE REFER TO NewCompressorLevel(),
// as this function creates a new Compressor (alloc 32KiB) which is then closed at the end of the function.
//
// See CompressLevel for further information.
func CompressZlibLevel(in, out []byte, level int) (int, []byte, error) {
	return CompressLevel(in, out, ModeZlib, level)
}

// Compress compresses the data from in to out and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though).
//
// m specifies which compression format should be used (e.g. ModeZlib). Uses default compression level.
//
// IF YOU WANT TO COMPRESS MORE THAN ONCE, PLEASE REFER TO NewCompressor(),
// as this function creates a new Compressor (alloc 32KiB) which is then closed at the end of the function.
//
// Notice that for extremely small or already highly compressed data,
// the compressed data could be larger than uncompressed.
// If out == nil: For a too large discrepancy (len(out) > 1000 + 2 * len(in)) Compress will error
func Compress(in, out []byte, m Mode) (int, []byte, error) {
	return CompressLevel(in, out, m, DefaultCompressionLevel)
}

// CompressLevel compresses the data from in to out using the specified compression level and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though).
//
// m specifies which compression format should be used (e.g. ModeZlib).
// Level defines the compression level.
//
// IF YOU WANT TO COMPRESS MORE THAN ONCE, PLEASE REFER TO NewCompressorLevel(),
// as this function creates a new Compressor (alloc 32KiB) which is then closed at the end of the function.
//
// Notice that for extremely small or already highly compressed data,
// the compressed data could be larger than uncompressed.
// If out == nil: For a too large discrepancy (len(out) > 1000 + 2 * len(in)) Compress will error
func CompressLevel(in, out []byte, m Mode, level int) (int, []byte, error) {
	c, err := NewCompressorLevel(level)
	if err != nil {
		return 0, out, err
	}
	defer c.Close()

	return c.Compress(in, out, m)
}
