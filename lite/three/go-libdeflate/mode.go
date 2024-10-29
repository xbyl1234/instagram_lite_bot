package libdeflate

// Mode specifies the type of compression/decompression such as zlib, gzip and raw DEFLATE
type Mode int

// The constants that specify a certain mode of compression/decompression
const (
	ModeDEFLATE Mode = iota
	ModeZlib
	ModeGzip
)
