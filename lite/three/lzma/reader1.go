package lzma

import (
	"errors"
	"fmt"
	"io"
)

type Reader1 struct {
	rangeDec  *rangeDecoder
	outWindow *window

	s             *state
	isEndOfStream bool
}

func NewReader1(inStream io.ByteReader) (*Reader1, error) {
	r := &Reader1{
		rangeDec: newRangeDecoder(inStream),
	}

	return r, r.initializeFull(inStream)
}

func NewReader1ForSevenZip(inStream io.ByteReader, props []byte, unpackSize uint64) (*Reader1, error) {
	lc, pb, lp, err := DecodeProp(props[0])
	if err != nil {
		return nil, err
	}

	dictSize, err := DecodeDictSize(props[1:5])
	if err != nil {
		return nil, err
	}

	r := &Reader1{
		rangeDec:  newRangeDecoder(inStream),
		outWindow: newWindow(dictSize),
	}

	return r, r.initialize(lc, pb, lp, unpackSize)
}

func NewReader1ForReader2(inStream io.ByteReader, prop byte, unpackSize uint64, outWindow *window) (*Reader1, error) {
	lc, pb, lp, err := DecodeProp(prop)
	if err != nil {
		return nil, err
	}

	r := &Reader1{
		outWindow: outWindow,
		rangeDec:  newRangeDecoder(inStream),
	}

	return r, r.initialize(lc, pb, lp, unpackSize)
}

func (r *Reader1) initializeFull(inStream io.ByteReader) error {
	b, err := inStream.ReadByte()
	if err != nil {
		return err
	}

	lc, pb, lp, err := DecodeProp(b)
	if err != nil {
		return fmt.Errorf("decode prop: %w", err)
	}

	dictSize, err := readDictSize(inStream)
	if err != nil {
		return fmt.Errorf("decode dict size: %w", err)
	}

	r.outWindow = newWindow(dictSize)

	unpackSize, err := readUnpackSize(inStream)
	if err != nil {
		return fmt.Errorf("decode unpack size: %w", err)
	}

	return r.initialize(lc, pb, lp, unpackSize)
}

func readDictSize(inStream io.ByteReader) (uint32, error) {
	var (
		b   byte
		err error
	)

	dictSize := uint32(0)
	for i := 0; i < 4; i++ {
		b, err = inStream.ReadByte()
		if err != nil {
			return 0, err
		}

		dictSize |= uint32(b) << (8 * i)
	}

	if dictSize < lzmaDicMin {
		dictSize = lzmaDicMin
	}

	if dictSize > lzmaDicMax {
		return 0, ErrDictOutOfRange
	}

	return dictSize, nil
}

func readUnpackSize(inStream io.ByteReader) (uint64, error) {
	var (
		b          byte
		err        error
		unpackSize uint64
	)

	for i := 0; i < 8; i++ {
		b, err = inStream.ReadByte()
		if err != nil {
			return 0, err
		}

		unpackSize |= uint64(b) << (8 * i)
	}

	return unpackSize, nil
}

func (r *Reader1) initialize(lc, pb, lp uint8, unpackSize uint64) error {
	r.s = newState(lc, pb, lp)
	r.s.SetUnpackSize(unpackSize)

	err := r.rangeDec.Init()
	if err != nil {
		return fmt.Errorf("rangeDec.Init: %w", err)
	}

	return nil
}

func (r *Reader1) Reset() {
	r.s.Reset()
	r.isEndOfStream = false
}

func (r *Reader1) Reopen(inStream io.ByteReader, unpackSize uint64) error {
	r.isEndOfStream = false
	r.s.SetUnpackSize(unpackSize)

	err := r.rangeDec.Reopen(inStream)
	if err != nil {
		return err
	}

	return nil
}

func DecodeUnpackSize(header []byte) uint64 {
	var (
		b          byte
		unpackSize uint64
	)

	for i := 0; i < 8; i++ {
		b = header[i]

		unpackSize |= uint64(b) << (8 * i)
	}

	return unpackSize
}

func DecodeDictSize(properties []byte) (uint32, error) {
	dictSize := uint32(0)
	for i := 0; i < 4; i++ {
		dictSize |= uint32(properties[i]) << (8 * i)
	}

	if dictSize < lzmaDicMin {
		dictSize = lzmaDicMin
	}

	if dictSize > lzmaDicMax {
		return 0, ErrDictOutOfRange
	}

	return dictSize, nil
}

func DecodeProp(d byte) (uint8, uint8, uint8, error) {
	if d >= (9 * 5 * 5) {
		return 0, 0, 0, ErrIncorrectProperties
	}

	lc := d % 9
	d /= 9
	pb := d / 5
	lp := d % 5

	return lc, pb, lp, nil
}

func (r *Reader1) Read(p []byte) (n int, err error) {
	var k int

	for {
		if r.outWindow.HasPending() {
			k, err = r.outWindow.ReadPending(p[n:])
			n += k
			if err != nil {
				return n, err
			}

			if n >= len(p) {
				return
			}
		}

		if r.isEndOfStream {
			err = io.EOF
			return
		}

		err = r.decompress(uint32(len(p) - n))
		if errors.Is(err, io.EOF) {
			r.isEndOfStream = true
			err = nil
		}
		if err != nil {
			return
		}
	}
}

func (r *Reader1) DecodeLiteral(state uint32, rep0 uint32) error {
	prevByte := uint32(0)
	if !r.outWindow.IsEmpty() {
		prevByte = uint32(r.outWindow.GetByte(1))
	}

	symbol := uint32(1)
	litState := ((r.outWindow.pos & ((1 << r.s.lp) - 1)) << r.s.lc) + (prevByte >> (8 - r.s.lc))
	probs := r.s.litProbs[(uint32(0x300) * litState):]
	rang := r.rangeDec.Range
	code := r.rangeDec.Code

	var (
		bit uint32
	)

	if state >= 7 {
		matchByte := r.outWindow.GetByte(rep0 + 1)

		var matchBit uint32

		for symbol < 0x100 {
			matchBit = uint32((matchByte >> 7) & 1)
			matchByte <<= 1

			// rc.DecodeBit begin
			v := probs[((1+matchBit)<<8)+symbol]
			bound := (rang >> kNumBitModelTotalBits) * uint32(v)

			if code < bound {
				v += ((1 << kNumBitModelTotalBits) - v) >> kNumMoveBits
				rang = bound
				bit = 0

				// Normalize
				if rang < kTopValue {
					b, err := r.rangeDec.inStream.ReadByte()
					if err != nil {
						return err
					}

					rang <<= 8
					code = (code << 8) | uint32(b)
				}

				probs[((1+matchBit)<<8)+symbol] = v
			} else {
				v -= v >> kNumMoveBits
				code -= bound
				rang -= bound
				bit = 1

				// Normalize
				if rang < kTopValue {
					b, err := r.rangeDec.inStream.ReadByte()
					if err != nil {
						return err
					}

					rang <<= 8
					code = (code << 8) | uint32(b)
				}

				probs[((1+matchBit)<<8)+symbol] = v
			}
			// rc.DecodeBit end

			symbol = (symbol << 1) | bit
			if matchBit != bit {
				break
			}
		}
	}

	for symbol < 0x100 {
		// rc.DecodeBit begin
		v := probs[symbol]
		bound := (rang >> kNumBitModelTotalBits) * uint32(v)

		if code < bound {
			v += ((1 << kNumBitModelTotalBits) - v) >> kNumMoveBits
			rang = bound
			bit = 0

			// Normalize
			if rang < kTopValue {
				b, err := r.rangeDec.inStream.ReadByte()
				if err != nil {
					return err
				}

				rang <<= 8
				code = (code << 8) | uint32(b)
			}

			probs[symbol] = v
		} else {
			v -= v >> kNumMoveBits
			code -= bound
			rang -= bound
			bit = 1

			// Normalize
			if rang < kTopValue {
				b, err := r.rangeDec.inStream.ReadByte()
				if err != nil {
					return err
				}

				rang <<= 8
				code = (code << 8) | uint32(b)
			}

			probs[symbol] = v
		}
		// rc.DecodeBit end

		symbol = (symbol << 1) | bit
	}

	r.rangeDec.Range = rang
	r.rangeDec.Code = code

	r.outWindow.PutByte(byte(symbol - 0x100))

	return nil
}

func (r *Reader1) DecodeDistance(len uint32) (uint32, error) {
	lenState := len
	if lenState > (kNumLenToPosStates - 1) {
		lenState = kNumLenToPosStates - 1
	}

	s := r.s
	posSlot, err := BitTreeDecode(s.posSlotDecoderProbs[lenState][:], posSlotDecoderNumBits, r.rangeDec)
	if err != nil {
		return 0, err
	}

	if posSlot < 4 {
		return posSlot, nil
	}

	numDirectBits := (posSlot >> 1) - 1
	dist := (2 | (posSlot & 1)) << numDirectBits

	var bit uint32

	if posSlot < kEndPosModelIndex {
		bit, err = BitTreeReverseDecode(s.posDecoders[dist-posSlot:], int(numDirectBits), r.rangeDec)
		if err != nil {
			return 0, err
		}
		dist += bit
	} else {
		bit, err = r.rangeDec.DecodeDirectBits(int(numDirectBits - kNumAlignBits))
		if err != nil {
			return 0, err
		}
		dist += bit << kNumAlignBits

		bit, err = BitTreeReverseDecode(s.alignDecoderProbs[:], kNumAlignBits, r.rangeDec)
		if err != nil {
			return 0, err
		}
		dist += bit
	}

	return dist, nil
}
