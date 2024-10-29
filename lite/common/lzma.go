package common

import (
	"bytes"
	"errors"
	"io"
)

const (
	STATES               = 12
	LIT_STATES           = 7
	LIT_LIT              = 0
	MATCH_LIT_LIT        = 1
	REP_LIT_LIT          = 2
	SHORTREP_LIT_LIT     = 3
	MATCH_LIT            = 4
	REP_LIT              = 5
	SHORTREP_LIT         = 6
	LIT_MATCH            = 7
	LIT_LONGREP          = 8
	LIT_SHORTREP         = 9
	NONLIT_MATCH         = 10
	NONLIT_REP           = 11
	SHIFT_BITS           = 8
	TOP_MASK             = 0xFF000000
	BIT_MODEL_TOTAL_BITS = 11
	BIT_MODEL_TOTAL      = 1 << BIT_MODEL_TOTAL_BITS
	PROB_INIT            = BIT_MODEL_TOTAL / 2
	MOVE_BITS            = 5
	INIT_SIZE            = 5
	POS_STATES_MAX       = 1 << 4
	MATCH_LEN_MIN        = 2
	LOW_SYMBOLS          = 1 << 3
	MID_SYMBOLS          = 1 << 3
	HIGH_SYMBOLS         = 1 << 8
	MATCH_LEN_MAX        = MATCH_LEN_MIN + LOW_SYMBOLS + MID_SYMBOLS + HIGH_SYMBOLS - 1
	DIST_STATES          = 4
	DIST_SLOTS           = 1 << 6
	DIST_MODEL_START     = 4
	DIST_MODEL_END       = 14
	FULL_DISTANCES       = 1 << (DIST_MODEL_END / 2)
	ALIGN_BITS           = 4
	ALIGN_SIZE           = 1 << ALIGN_BITS
	ALIGN_MASK           = ALIGN_SIZE - 1
	REPS                 = 4
	DICT_SIZE_MIN        = 4096
	DICT_SIZE_MAX        = int(^uint(0)>>1) & ^15
	COMPRESSED_SIZE_MAX  = 1 << 16
)

type XZIOException struct {
	message string
}

func (e *XZIOException) Error() string {
	return e.message
}

type State struct {
	state int
}

func NewState() *State {
	return &State{}
}

func (s *State) Reset() {
	s.state = LIT_LIT
}

func (s *State) Get() int {
	return s.state
}

func (s *State) Set(other *State) {
	s.state = other.state
}

func (s *State) UpdateLiteral() {
	if s.state <= SHORTREP_LIT_LIT {
		s.state = LIT_LIT
	} else if s.state <= LIT_SHORTREP {
		s.state -= 3
	} else {
		s.state -= 6
	}
}

func (s *State) UpdateMatch() {
	if s.state < LIT_STATES {
		s.state = LIT_MATCH
	} else {
		s.state = NONLIT_MATCH
	}
}

func (s *State) UpdateLongRep() {
	if s.state < LIT_STATES {
		s.state = LIT_LONGREP
	} else {
		s.state = NONLIT_REP
	}
}

func (s *State) UpdateShortRep() {
	if s.state < LIT_STATES {
		s.state = LIT_SHORTREP
	} else {
		s.state = NONLIT_REP
	}
}

func (s *State) IsLiteral() bool {
	return s.state < LIT_STATES
}

type RangeDecoder struct {
	buf    []byte
	pos    int
	end    int
	range_ uint32
	code   uint32
}

func NewRangeDecoder(inputSizeMax int) *RangeDecoder {
	return &RangeDecoder{
		buf: make([]byte, inputSizeMax-INIT_SIZE),
	}
}

func (rd *RangeDecoder) PrepareInputBuffer(in io.Reader, len int) error {
	if len < INIT_SIZE {
		return &XZIOException{"Invalid input size"}
	}
	buf := make([]byte, len)
	if _, err := io.ReadFull(in, buf); err != nil {
		return err
	}
	if buf[0] != 0 {
		return &XZIOException{"buf[0] != 0"}
	}
	rd.code = uint32(buf[1])<<24 | uint32(buf[2])<<16 | uint32(buf[3])<<8 | uint32(buf[4])
	rd.range_ = 0xffffffff
	rd.pos = 0
	rd.end = len - INIT_SIZE
	copy(rd.buf, buf[INIT_SIZE:])
	return nil
}

func (rd *RangeDecoder) IsInBufferOK() bool {
	return rd.pos <= rd.end
}

func (rd *RangeDecoder) IsFinished() bool {
	return rd.pos == rd.end && rd.code == 0
}

func (rd *RangeDecoder) Normalize() error {
	if (rd.range_ & TOP_MASK) == 0 {
		if rd.pos >= rd.end {
			return &XZIOException{"Input buffer overflow"}
		}
		rd.code = rd.code<<SHIFT_BITS | uint32(rd.buf[rd.pos])
		rd.range_ <<= SHIFT_BITS
		rd.pos++
	}
	return nil
}

func (rd *RangeDecoder) DecodeBit(probs []int16, index int) (int, error) {
	if err := rd.Normalize(); err != nil {
		return 0, err
	}
	prob := probs[index]
	bound := (rd.range_ >> BIT_MODEL_TOTAL_BITS) * uint32(prob)
	var bit int
	if rd.code < bound {
		rd.range_ = bound
		probs[index] = prob + ((BIT_MODEL_TOTAL - prob) >> MOVE_BITS)
		bit = 0
	} else {
		rd.range_ -= bound
		rd.code -= bound
		probs[index] = prob - (prob >> MOVE_BITS)
		bit = 1
	}
	return bit, nil
}

func (rd *RangeDecoder) DecodeBitTree(probs []int16) (int, error) {
	symbol := 1
	for symbol < len(probs) {
		bit, err := rd.DecodeBit(probs, symbol)
		if err != nil {
			return 0, err
		}
		symbol = (symbol << 1) | bit
	}
	return symbol - len(probs), nil
}

func (rd *RangeDecoder) DecodeReverseBitTree(probs []int16) (int, error) {
	symbol := 1
	i := 0
	result := 0
	for symbol < len(probs) {
		bit, err := rd.DecodeBit(probs, symbol)
		if err != nil {
			return 0, err
		}
		symbol = (symbol << 1) | bit
		result |= bit << i
		i++
	}
	return result, nil
}

func (rd *RangeDecoder) DecodeDirectBits(count int) (int, error) {
	result := uint32(0)
	for count > 0 {
		if err := rd.Normalize(); err != nil {
			return 0, err
		}
		rd.range_ >>= 1
		t := (rd.code - rd.range_) >> 31
		rd.code -= rd.range_ & (t - 1)
		result = ((result) << 1) | (1 - uint32(t))
		count--
	}
	return int(result), nil
}

type LiteralSubdecoder struct {
	probs []int16
	obj   *LZMADecoder
}

func NewLiteralSubdecoder(obj *LZMADecoder) *LiteralSubdecoder {
	return &LiteralSubdecoder{
		probs: make([]int16, 0x300),
		obj:   obj,
	}
}

func (ls *LiteralSubdecoder) Reset() {
	for i := range ls.probs {
		ls.probs[i] = PROB_INIT
	}
}

func (ls *LiteralSubdecoder) Decode() error {
	symbol := 1
	if ls.obj.state.IsLiteral() {
		for symbol < 0x100 {
			bit, err := ls.obj.rc.DecodeBit(ls.probs, symbol)
			if err != nil {
				return err
			}
			symbol = (symbol << 1) | bit
		}
	} else {
		matchByte := ls.obj.lz.GetByte(ls.obj.reps[0])
		offset := 0x100
		for symbol < 0x100 {
			matchByte <<= 1
			matchBit := matchByte & offset
			bit, err := ls.obj.rc.DecodeBit(ls.probs, offset+matchBit+symbol)
			if err != nil {
				return err
			}
			symbol = (symbol << 1) | bit
			offset &= (0 - bit) ^ ^matchBit
		}
	}
	ls.obj.lz.PutByte(byte(symbol))
	ls.obj.state.UpdateLiteral()
	return nil
}

type LiteralDecoder struct {
	lc             int
	literalPosMask int
	subdecoders    []*LiteralSubdecoder
	obj            *LZMADecoder
}

func NewLiteralDecoder(lc, lp int, obj *LZMADecoder) *LiteralDecoder {
	subdecoders := make([]*LiteralSubdecoder, 1<<(lc+lp))
	for i := range subdecoders {
		subdecoders[i] = NewLiteralSubdecoder(obj)
	}
	return &LiteralDecoder{
		lc:             lc,
		literalPosMask: (1 << lp) - 1,
		subdecoders:    subdecoders,
		obj:            obj,
	}
}

func (ld *LiteralDecoder) GetSubcoderIndex(prevByte int, pos int) int {
	low := prevByte >> (8 - ld.lc)
	high := (pos & ld.literalPosMask) << ld.lc
	return low + high
}

func (ld *LiteralDecoder) Reset() {
	for _, subdecoder := range ld.subdecoders {
		subdecoder.Reset()
	}
}

func (ld *LiteralDecoder) Decode() error {
	i := ld.GetSubcoderIndex(ld.obj.lz.GetByte(0), ld.obj.lz.GetPos())
	return ld.subdecoders[i].Decode()
}

type LengthDecoder struct {
	choice []int16
	low    [][]int16
	mid    [][]int16
	high   []int16
	obj    *LZMADecoder
}

func NewLengthDecoder(obj *LZMADecoder) *LengthDecoder {
	choice := make([]int16, 2)
	low := make([][]int16, POS_STATES_MAX)
	for i := range low {
		low[i] = make([]int16, LOW_SYMBOLS)
	}
	mid := make([][]int16, POS_STATES_MAX)
	for i := range mid {
		mid[i] = make([]int16, MID_SYMBOLS)
	}
	high := make([]int16, HIGH_SYMBOLS)
	return &LengthDecoder{
		choice: choice,
		low:    low,
		mid:    mid,
		high:   high,
		obj:    obj,
	}
}

func (ld *LengthDecoder) Reset() {
	for _, probs := range ld.low {
		for i := range probs {
			probs[i] = PROB_INIT
		}
	}
	for _, probs := range ld.mid {
		for i := range probs {
			probs[i] = PROB_INIT
		}
	}
	for i := range ld.high {
		ld.high[i] = PROB_INIT
	}
	ld.choice[0] = PROB_INIT
	ld.choice[1] = PROB_INIT
}

func (ld *LengthDecoder) Decode(posState int) (int, error) {
	if bit, err := ld.obj.rc.DecodeBit(ld.choice, 0); err != nil {
		return 0, err
	} else if bit == 0 {
		tree, err := ld.obj.rc.DecodeBitTree(ld.low[posState])
		return tree + MATCH_LEN_MIN, err
	} else if bit, err := ld.obj.rc.DecodeBit(ld.choice, 1); err != nil {
		return 0, err
	} else if bit == 0 {
		tree, err := ld.obj.rc.DecodeBitTree(ld.mid[posState])
		return tree + MATCH_LEN_MIN + LOW_SYMBOLS, err
	} else {
		tree, err := ld.obj.rc.DecodeBitTree(ld.high)
		return tree + MATCH_LEN_MIN + LOW_SYMBOLS + MID_SYMBOLS, err
	}
}

type LZDecoder struct {
	buf         []byte
	start       int
	pos         int
	full        int
	limit       int
	pendingLen  int
	pendingDist int
}

func NewLZDecoder(dictSize int, presetDict []byte) *LZDecoder {
	buf := make([]byte, dictSize)
	if presetDict != nil {
		pos := len(presetDict)
		copy(buf[dictSize-pos:], presetDict)
		return &LZDecoder{
			buf:   buf,
			start: dictSize - pos,
			pos:   dictSize,
			full:  dictSize,
		}
	}
	return &LZDecoder{
		buf:   buf,
		start: 0,
		pos:   0,
		full:  0,
	}
}

func (ld *LZDecoder) Reset() {
	ld.start = 0
	ld.pos = 0
	ld.full = 0
	ld.limit = 0
	ld.buf[len(ld.buf)-1] = 0x00
}

func (ld *LZDecoder) SetLimit(outMax int) {
	if len(ld.buf)-ld.pos <= outMax {
		ld.limit = len(ld.buf)
	} else {
		ld.limit = ld.pos + outMax
	}
}

func (ld *LZDecoder) HasSpace() bool {
	return ld.pos < ld.limit
}

func (ld *LZDecoder) HasPending() bool {
	return ld.pendingLen > 0
}

func (ld *LZDecoder) GetPos() int {
	return ld.pos
}

func (ld *LZDecoder) GetByte(dist int) int {
	offset := ld.pos - dist - 1
	if dist >= ld.pos {
		offset += len(ld.buf)
	}
	return int(ld.buf[offset]) & 0xFF
}

func (ld *LZDecoder) PutByte(b byte) {
	ld.buf[ld.pos] = b
	ld.pos++
	if ld.full < ld.pos {
		ld.full = ld.pos
	}
}

func (ld *LZDecoder) Repeat(dist, length int) error {
	if dist < 0 || dist >= ld.full {
		return &XZIOException{"Invalid distance"}
	}
	left := min(ld.limit-ld.pos, length)
	ld.pendingLen = length - left
	ld.pendingDist = dist
	back := ld.pos - dist - 1
	if dist >= ld.pos {
		back += len(ld.buf)
	}
	for left > 0 {
		ld.buf[ld.pos] = ld.buf[back]
		ld.pos++
		back++
		if back == len(ld.buf) {
			back = 0
		}
		left--
	}
	if ld.full < ld.pos {
		ld.full = ld.pos
	}
	return nil
}

func (ld *LZDecoder) RepeatPending() error {
	if ld.pendingLen > 0 {
		if err := ld.Repeat(ld.pendingDist, ld.pendingLen); err != nil {
			return err
		}
	}
	return nil
}

func (ld *LZDecoder) CopyUncompressed(in io.Reader, length int) error {
	copySize := min(len(ld.buf)-ld.pos, length)
	if _, err := io.ReadFull(in, ld.buf[ld.pos:ld.pos+copySize]); err != nil {
		return err
	}
	ld.pos += copySize
	if ld.full < ld.pos {
		ld.full = ld.pos
	}
	return nil
}

func (ld *LZDecoder) Flush(out []byte, outOff int) int {
	copySize := ld.pos - ld.start
	if ld.pos == len(ld.buf) {
		ld.pos = 0
	}
	copy(out[outOff:][:copySize], ld.buf[ld.start:][:copySize])
	ld.start = ld.pos
	return copySize

}

type LZMADecoder struct {
	lz              *LZDecoder
	rc              *RangeDecoder
	literalDecoder  *LiteralDecoder
	matchLenDecoder *LengthDecoder
	repLenDecoder   *LengthDecoder
	posMask         int
	reps            [REPS]int
	state           *State
	isMatch         [][]int16
	isRep           []int16
	isRep0          []int16
	isRep1          []int16
	isRep2          []int16
	isRep0Long      [][]int16
	distSlots       [][]int16
	distSpecial     [][]int16
	distAlign       []int16
}

func Make2DArray(d1, d2 int16) [][]int16 {
	var arr = make([][]int16, d1)
	for i := int16(0); i < d1; i++ {
		arr[i] = make([]int16, d2)
	}
	return arr
}

func NewLZMADecoder(lz *LZDecoder, rc *RangeDecoder, lc, lp, pb int) *LZMADecoder {
	ret := &LZMADecoder{
		lz:         lz,
		rc:         rc,
		posMask:    (1 << pb) - 1,
		state:      NewState(),
		isMatch:    Make2DArray(STATES, POS_STATES_MAX),
		isRep:      make([]int16, STATES),
		isRep0:     make([]int16, STATES),
		isRep1:     make([]int16, STATES),
		isRep2:     make([]int16, STATES),
		isRep0Long: Make2DArray(STATES, POS_STATES_MAX),
		distSlots:  Make2DArray(DIST_STATES, DIST_SLOTS),
		distSpecial: [][]int16{make([]int16, 2),
			make([]int16, 2),
			make([]int16, 4),
			make([]int16, 4),
			make([]int16, 8),
			make([]int16, 8),
			make([]int16, 16),
			make([]int16, 16),
			make([]int16, 32),
			make([]int16, 32)},
		distAlign: make([]int16, ALIGN_SIZE),
	}

	ret.literalDecoder = NewLiteralDecoder(lc, lp, ret)
	ret.matchLenDecoder = NewLengthDecoder(ret)
	ret.repLenDecoder = NewLengthDecoder(ret)
	ret.Reset()
	return ret
}

func (ld *LZMADecoder) Reset() {
	ld.reps[0] = 0
	ld.reps[1] = 0
	ld.reps[2] = 0
	ld.reps[3] = 0
	ld.state.Reset()
	for _, probs := range ld.isMatch {
		for i := range probs {
			probs[i] = PROB_INIT
		}
	}
	for i := range ld.isRep {
		ld.isRep[i] = PROB_INIT
	}
	for i := range ld.isRep0 {
		ld.isRep0[i] = PROB_INIT
	}
	for i := range ld.isRep1 {
		ld.isRep1[i] = PROB_INIT
	}
	for i := range ld.isRep2 {
		ld.isRep2[i] = PROB_INIT
	}
	for _, probs := range ld.isRep0Long {
		for i := range probs {
			probs[i] = PROB_INIT
		}
	}
	for _, probs := range ld.distSlots {
		for i := range probs {
			probs[i] = PROB_INIT
		}
	}
	for _, probs := range ld.distSpecial {
		for i := range probs {
			probs[i] = PROB_INIT
		}
	}
	for i := range ld.distAlign {
		ld.distAlign[i] = PROB_INIT
	}
	ld.literalDecoder.Reset()
	ld.matchLenDecoder.Reset()
	ld.repLenDecoder.Reset()
}

func (ld *LZMADecoder) Decode() error {
	var err error
	err = ld.lz.RepeatPending()
	if err != nil {
		return err
	}
	for ld.lz.HasSpace() {
		posState := ld.lz.GetPos() & ld.posMask
		bit, err := ld.rc.DecodeBit(ld.isMatch[ld.state.Get()], posState)
		if err != nil {
			return err
		}
		if bit == 0 {
			err = ld.literalDecoder.Decode()
			if err != nil {
				return err
			}
		} else {
			var l int
			bit, err = ld.rc.DecodeBit(ld.isRep, ld.state.Get())
			if err != nil {
				return err
			}
			if bit == 0 {
				l, err = ld.decodeMatch(posState)
				if err != nil {
					return err
				}
			} else {
				l, err = ld.decodeRepMatch(posState)
				if err != nil {
					return err
				}
			}
			if err = ld.lz.Repeat(ld.reps[0], l); err != nil {
				return err
			}
		}
	}
	if err := ld.rc.Normalize(); err != nil {
		return err
	}
	if !ld.rc.IsInBufferOK() {
		return &XZIOException{"Input buffer overflow"}
	}
	return nil
}

func (ld *LZMADecoder) decodeMatch(posState int) (int, error) {
	ld.state.UpdateMatch()
	ld.reps[3] = ld.reps[2]
	ld.reps[2] = ld.reps[1]
	ld.reps[1] = ld.reps[0]
	len, _ := ld.matchLenDecoder.Decode(posState)
	distSlot, _ := ld.rc.DecodeBitTree(ld.distSlots[getDistState(len)])
	if distSlot < DIST_MODEL_START {
		ld.reps[0] = distSlot
	} else {
		limit := (distSlot >> 1) - 1
		ld.reps[0] = (2 | (distSlot & 1)) << limit
		if distSlot < DIST_MODEL_END {
			d, err := ld.rc.DecodeReverseBitTree(ld.distSpecial[distSlot-DIST_MODEL_START])
			if err != nil {
				return 0, err
			}
			ld.reps[0] |= d
		} else {
			d, err := ld.rc.DecodeDirectBits(limit - ALIGN_BITS)
			if err != nil {
				return 0, err
			}
			ld.reps[0] |= d << ALIGN_BITS
			d, err = ld.rc.DecodeReverseBitTree(ld.distAlign)
			if err != nil {
				return 0, err
			}
			ld.reps[0] |= d
		}
	}
	return len, nil
}

func (ld *LZMADecoder) decodeRepMatch(posState int) (int, error) {
	if bit, err := ld.rc.DecodeBit(ld.isRep0, ld.state.Get()); err != nil {
		return 0, err
	} else if bit == 0 {
		if bit, err := ld.rc.DecodeBit(ld.isRep0Long[ld.state.Get()], posState); err != nil {
			return 0, err
		} else if bit == 0 {
			ld.state.UpdateShortRep()
			return 1, nil
		}
	} else {
		var tmp int
		if bit, err := ld.rc.DecodeBit(ld.isRep1, ld.state.Get()); err != nil {
			return 0, err
		} else if bit == 0 {
			tmp = ld.reps[1]
		} else {
			if bit, err := ld.rc.DecodeBit(ld.isRep2, ld.state.Get()); err != nil {
				return 0, err
			} else if bit == 0 {
				tmp = ld.reps[2]
			} else {
				tmp = ld.reps[3]
				ld.reps[3] = ld.reps[2]
			}
			ld.reps[2] = ld.reps[1]
		}
		ld.reps[1] = ld.reps[0]
		ld.reps[0] = tmp
	}
	ld.state.UpdateLongRep()
	return ld.repLenDecoder.Decode(posState)
}

func getDictSize(dictSize int) int {
	if dictSize < DICT_SIZE_MIN || dictSize > DICT_SIZE_MAX {
		panic("Unsupported dictionary size")
	}
	return (dictSize + 15) & ^15
}

func getDistState(len int) int {
	if len < DIST_STATES+MATCH_LEN_MIN {
		return len - MATCH_LEN_MIN
	}
	return DIST_STATES - 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type LZMA2InputStream struct {
	in               *bytes.Buffer
	lz               *LZDecoder
	rc               *RangeDecoder
	lzma             *LZMADecoder
	uncompressedSize int
	isLZMAChunk      bool
	needDictReset    bool
	needProps        bool
	endReached       bool
	exception        error
}

func NewLZMA2InputStream(dictSize int) *LZMA2InputStream {
	lz := NewLZDecoder(getDictSize(dictSize), nil)
	needDictReset := true
	//if presetDict != nil && len(presetDict) > 0 {
	//	needDictReset = false
	//}
	return &LZMA2InputStream{
		in:            bytes.NewBuffer([]byte{}),
		lz:            lz,
		rc:            NewRangeDecoder(COMPRESSED_SIZE_MAX),
		needDictReset: needDictReset,
	}
}

func (s *LZMA2InputStream) Write(b []byte) {
	s.in.Write(b)
}

func (s *LZMA2InputStream) ReadAll() ([]byte, error) {
	var err error
	var l int
	var rl int
	l, err = s.Available()
	if err != nil {
		return nil, err
	}
	b := make([]byte, l)
	rl, err = s.Read(b, 0, l)
	if err != nil || rl != l {
		return b, err
	}
	return b, nil
}

func (s *LZMA2InputStream) ReadByte() (int, error) {
	buf := make([]byte, 1)
	n, err := s.Read(buf, 0, 1)
	if err != nil {
		return -1, err
	}
	if n == -1 {
		return -1, nil
	}
	return int(buf[0]) & 0xFF, nil
}

func (s *LZMA2InputStream) Read(buf []byte, off, length int) (int, error) {
	if off < 0 || length < 0 || off+length < 0 || off+length > len(buf) {
		return 0, errors.New("Invalid buffer parameters")
	}
	if length == 0 {
		return 0, nil
	}
	if s.in == nil {
		return 0, errors.New("Stream closed")
	}
	if s.exception != nil {
		return 0, s.exception
	}
	if s.endReached {
		return -1, nil
	}
	size := 0
	for length > 0 {
		if s.uncompressedSize == 0 {
			err := s.decodeChunkHeader()
			if err != nil {
				if err == io.EOF {
					if size == 0 {
						return -1, nil
					}
					return size, nil
				}
				return size, err
			}
		}
		copySizeMax := min(s.uncompressedSize, length)
		if !s.isLZMAChunk {
			err := s.lz.CopyUncompressed(s.in, copySizeMax)
			if err != nil {
				return size, err
			}
		} else {
			s.lz.SetLimit(copySizeMax)
			err := s.lzma.Decode()
			if err != nil {
				return size, err
			}
		}
		copiedSize := s.lz.Flush(buf, off)
		off += copiedSize
		length -= copiedSize
		size += copiedSize
		s.uncompressedSize -= copiedSize
		if s.uncompressedSize == 0 {
			if !s.rc.IsFinished() || s.lz.HasPending() {
				return size, errors.New("Invalid LZMA2 stream")
			}
		}
	}
	return size, nil
}

func (s *LZMA2InputStream) decodeChunkHeader() error {
	control := make([]byte, 1)
	_, err := io.ReadFull(s.in, control)
	if err != nil {
		return err
	}
	if control[0] == 0x00 {
		s.endReached = true
		return nil
	}
	if control[0] >= 0xE0 || control[0] == 0x01 {
		s.needProps = true
		s.needDictReset = false
		s.lz.Reset()
	} else if s.needDictReset {
		return errors.New("Invalid LZMA2 stream")
	}
	if control[0] >= 0x80 {
		s.isLZMAChunk = true
		s.uncompressedSize = int(control[0]&0x1F) << 16
		s.uncompressedSize += int(s.readUnsignedShort() + 1)
		compressedSize := int(s.readUnsignedShort() + 1)
		if control[0] >= 0xC0 {
			s.needProps = false
			err := s.decodeProps()
			if err != nil {
				return err
			}
		} else if s.needProps {
			return errors.New("Invalid LZMA2 stream")
		} else if control[0] >= 0xA0 {
			s.lzma.Reset()
		}
		err := s.rc.PrepareInputBuffer(s.in, compressedSize)
		if err != nil {
			return err
		}
	} else if control[0] > 0x02 {
		return errors.New("Invalid LZMA2 stream")
	} else {
		s.isLZMAChunk = false
		s.uncompressedSize = int(s.readUnsignedShort() + 1)
	}
	return nil
}

func (s *LZMA2InputStream) decodeProps() error {
	props := make([]byte, 1)
	_, err := io.ReadFull(s.in, props)
	if err != nil {
		return err
	}
	if props[0] > (4*5+4)*9+8 {
		return errors.New("Invalid LZMA2 stream")
	}
	pb := props[0] / (9 * 5)
	props[0] -= pb * 9 * 5
	lp := props[0] / 9
	lc := props[0] - lp*9
	if lc+lp > 4 {
		return errors.New("invalid LZMA2 stream")
	}
	s.lzma = NewLZMADecoder(s.lz, s.rc, int(lc), int(lp), int(pb))
	return nil
}

func (s *LZMA2InputStream) readUnsignedShort() uint16 {
	buf := make([]byte, 2)
	_, _ = io.ReadFull(s.in, buf)
	return uint16(buf[0])<<8 | uint16(buf[1])
}

func (s *LZMA2InputStream) Available() (int, error) {
	if s.in == nil {
		return 0, errors.New("stream closed")
	}
	if s.exception != nil {
		return 0, s.exception
	}
	if s.uncompressedSize == 0 {
		err := s.decodeChunkHeader()
		if err != nil {
			return 0, err
		}
	}
	return s.uncompressedSize, nil
}
