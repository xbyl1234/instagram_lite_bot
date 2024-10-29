package libdeflate

import "errors"

var (
	errorInvalidModeCompressor   = errors.New("libdeflate: compressor: invalid mode")
	errorInvalidModeDecompressor = errors.New("libdeflate: decompressor: invalid mode")
)
