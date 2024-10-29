package libdeflate

// DecompressZlib decompresses the given zlib data from in to out and returns the number of consumed bytes c
// from 'in' and 'out' or an error if something went wrong.
//
// c is the number of bytes that were read before the BFINAL flag was
// encountered, which indicates the end of the compressed data.
//
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil to out, this function will allocate a sufficient buffer and return it.
//
// IF YOU WANT TO DECOMPRESS MORE THAN ONCE, PLEASE REFER TO NewDecompressor(),
// as this function creates a new Decompressor (alloc 32KiB) which is then closed at the end of the function.
//
// If error != nil, the data in out is undefined.
func DecompressZlib(in, out []byte) (int, []byte, error) {
	return Decompress(in, out, ModeZlib)
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
// IF YOU WANT TO DECOMPRESS MORE THAN ONCE, PLEASE REFER TO NewDecompressor(),
// as this function creates a new Decompressor (alloc 32KiB) which is then closed at the end of the function.
//
// If error != nil, the data in out is undefined.
func Decompress(in, out []byte, m Mode) (int, []byte, error) {
	dc, err := NewDecompressor()
	if err != nil {
		return 0, out, err
	}
	defer dc.Close()

	return dc.Decompress(in, out, m)
}
