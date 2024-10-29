package lzma

const lzmaHeaderLen = 13

const (
	lzmaDicMin = 1 << 12
	lzmaDicMax = 1<<32 - 1
)

const (
	kNumBitModelTotalBits = 11
	kNumMoveBits          = 5
	probInitVal           = (1 << kNumBitModelTotalBits) / 2

	kNumPosBitsMax = 4

	kNumStates          = 12
	kNumLenToPosStates  = 4
	kNumAlignBits       = 4
	kStartPosModelIndex = 4
	kEndPosModelIndex   = 14
	kNumFullDistances   = 1 << (kEndPosModelIndex >> 1)
	kMatchMinLen        = 2

	kTopValue = uint32(1) << 24
)

const (
	posSlotDecoderNumBits = 6

	lenLowCoderNumBits  = 3
	lenMidCoderNumBits  = 3
	lenHighCoderNumBits = 8
)

const lzmaRequiredInputMax = 20

const rangeDecoderHeaderLen = 5

// minMatchLen and maxMatchLen give the minimum and maximum values for
// encoding and decoding length values. minMatchLen is also used as base
// for the encoded length values.
const (
	minMatchLen = 2
	maxMatchLen = minMatchLen + 16 + 256 - 1
)

type chunkType int

const (
	chunkEndOfStream chunkType = iota
	chunkUncompressedResetDict
	chunkUncompressedNoResetDict

	chunkLZMANoReset
	chunkLZMAResetState
	chunkLZMAResetStateNewProp
	chunkLZMAResetStateNewPropResetDict
)

var isChunkResetDict = map[chunkType]bool{
	chunkUncompressedResetDict:          true,
	chunkLZMAResetStateNewPropResetDict: true,
}

var isChunkNewProp = map[chunkType]bool{
	chunkLZMAResetStateNewProp:          true,
	chunkLZMAResetStateNewPropResetDict: true,
}

var isChunkUncompressed = map[chunkType]bool{
	chunkUncompressedResetDict:   true,
	chunkUncompressedNoResetDict: true,
}

var isChunkLZMA = map[chunkType]bool{
	chunkLZMANoReset:                    true,
	chunkLZMAResetState:                 true,
	chunkLZMAResetStateNewProp:          true,
	chunkLZMAResetStateNewPropResetDict: true,
}

const (
	endOfStreamCode         = 0
	uncompressedResetDict   = 0b01
	uncompressedNoResetDict = 0b10

	maskLZMANoReset                    = 0b100
	maskLZMAResetState                 = 0b101
	maskLZMAResetStateNewProp          = 0b110
	maskLZMAResetStateNewPropResetDict = 0b111

	maskLZMAUncompressedSize = 0b11111
)

type prob uint16
