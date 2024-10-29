package lzma

type state struct {
	litProbs []prob

	posSlotDecoderProbs [kNumLenToPosStates][1 << posSlotDecoderNumBits]prob
	posDecoders         [1 + kNumFullDistances - kEndPosModelIndex]prob
	alignDecoderProbs   [1 << kNumAlignBits]prob

	lenDecoderChoice    prob
	lenDecoderChoice2   prob
	lenDecoderLowCoder  [1 << kNumPosBitsMax][1 << lenLowCoderNumBits]prob
	lenDecoderMidCoder  [1 << kNumPosBitsMax][1 << lenMidCoderNumBits]prob
	lenDecoderHighCoder [1 << lenHighCoderNumBits]prob

	repLenDecoderChoice    prob
	repLenDecoderChoice2   prob
	repLenDecoderLowCoder  [1 << kNumPosBitsMax][1 << lenLowCoderNumBits]prob
	repLenDecoderMidCoder  [1 << kNumPosBitsMax][1 << lenMidCoderNumBits]prob
	repLenDecoderHighCoder [1 << lenHighCoderNumBits]prob

	isMatch    [kNumStates << kNumPosBitsMax]prob
	isRep      [kNumStates]prob
	isRepG0    [kNumStates]prob
	isRepG1    [kNumStates]prob
	isRepG2    [kNumStates]prob
	isRep0Long [kNumStates << kNumPosBitsMax]prob

	markerIsMandatory bool
	unpackSizeDefined bool
	lc, pb, lp        uint8

	unpackSize             uint64
	bytesLeft              uint64
	posMask                uint32
	rep0, rep1, rep2, rep3 uint32
	state                  uint32
	posState               uint32
}

func newState(lc, pb, lp uint8) *state {
	s := &state{
		litProbs: make([]prob, uint32(0x300)<<(lc+lp)),

		lc: lc,
		pb: pb,
		lp: lp,

		posMask: (1 << pb) - 1,
	}

	s.Reset()

	return s
}

func (s *state) Renew(lc, pb, lp uint8) {
	s.lc = lc
	s.pb = pb
	s.lp = lp
	s.posMask = (1 << pb) - 1

	litProbsCount := int(0x300) << (lc + lp)
	if litProbsCount > cap(s.litProbs) {
		s.litProbs = make([]prob, litProbsCount)
	} else {
		s.litProbs = s.litProbs[:litProbsCount]
	}

	s.Reset()
}

func (s *state) Reset() {
	initProbs(s.litProbs)

	for i := 0; i < len(s.posSlotDecoderProbs); i++ {
		initProbs(s.posSlotDecoderProbs[i][:])
	}

	initProbs(s.alignDecoderProbs[:])
	initProbs(s.posDecoders[:])

	initProbs(s.isMatch[:])
	initProbs(s.isRep[:])
	initProbs(s.isRepG0[:])
	initProbs(s.isRepG1[:])
	initProbs(s.isRepG2[:])
	initProbs(s.isRep0Long[:])

	{ // lenDecoder
		s.lenDecoderChoice = probInitVal
		s.lenDecoderChoice2 = probInitVal
		initProbs(s.lenDecoderHighCoder[:])

		for i := 0; i < len(s.lenDecoderLowCoder); i++ {
			initProbs(s.lenDecoderLowCoder[i][:])
			initProbs(s.lenDecoderMidCoder[i][:])
		}
	}

	{ // repLenDecoder
		s.repLenDecoderChoice = probInitVal
		s.repLenDecoderChoice2 = probInitVal
		initProbs(s.repLenDecoderHighCoder[:])

		for i := 0; i < len(s.repLenDecoderLowCoder); i++ {
			initProbs(s.repLenDecoderLowCoder[i][:])
			initProbs(s.repLenDecoderMidCoder[i][:])
		}
	}

	s.rep0, s.rep1, s.rep2, s.rep3 = 0, 0, 0, 0
	s.state = 0
	s.posState = 0
}

func (s *state) SetUnpackSize(unpackSize uint64) {
	s.bytesLeft = unpackSize
	s.unpackSize = unpackSize

	s.unpackSizeDefined = isUnpackSizeDefined(unpackSize)
	s.markerIsMandatory = !s.unpackSizeDefined
}

func (s *state) decompressed() uint64 {
	return s.unpackSize - s.bytesLeft
}

func isUnpackSizeDefined(unpackSize uint64) bool {
	var (
		b                 byte
		unpackSizeDefined bool
	)

	for i := 0; i < 8; i++ {
		b = byte(unpackSize & 0xFF)
		if b != 0xFF {
			unpackSizeDefined = true
		}

		unpackSize >>= 8
	}

	return unpackSizeDefined
}

func stateUpdateLiteral(state uint32) uint32 {
	if state < 4 {
		return 0
	}

	if state < 10 {
		return state - 3
	}

	return state - 6
}

func stateUpdateMatch(state uint32) uint32 {
	if state < 7 {
		return 7
	}

	return 10
}

func stateUpdateRep(state uint32) uint32 {
	if state < 7 {
		return 8
	}

	return 11
}

func stateUpdateShortRep(state uint32) uint32 {
	if state < 7 {
		return 9
	}

	return 11
}
