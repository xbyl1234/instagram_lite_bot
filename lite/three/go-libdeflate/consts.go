package libdeflate

import "CentralizedControl/three/go-libdeflate/native"

// These constants specify several special compression levels
const (
	MinCompressionLevel        = native.MinCompressionLevel
	MaxStdZlibCompressionLevel = native.MaxStdZlibCompressionLevel
	MaxCompressionLevel        = native.MaxCompressionLevel
	DefaultCompressionLevel    = native.DefaultCompressionLevel
)
