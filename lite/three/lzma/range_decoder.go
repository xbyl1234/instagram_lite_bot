package lzma

import (
	"io"
)

type rangeDecoder struct {
	inStream io.ByteReader

	Range     uint32
	Code      uint32
	Corrupted bool
}

func newRangeDecoder(inStream io.ByteReader) *rangeDecoder {
	return &rangeDecoder{
		inStream: inStream,

		Range: 0xFFFFFFFF,
	}
}

func (d *rangeDecoder) IsFinishedOK() bool {
	return d.Code == 0
}

func (d *rangeDecoder) Init() error {
	b, err := d.inStream.ReadByte()
	if err != nil {
		return err
	}
	if b != 0 {
		return ErrResultError
	}

	for i := 0; i < 4; i++ {
		b, err = d.inStream.ReadByte()
		if err != nil {
			return err
		}

		d.Code = (d.Code << 8) | uint32(b)
	}

	return nil
}

func (d *rangeDecoder) Reopen(inStream io.ByteReader) error {
	d.inStream = inStream
	d.Corrupted = false
	d.Range = 0xFFFFFFFF
	d.Code = 0

	return d.Init()
}

func (d *rangeDecoder) DecodeBit(v *prob) (uint32, error) {
	bound := (d.Range >> kNumBitModelTotalBits) * uint32(*v)

	if d.Code < bound {
		*v += ((1 << kNumBitModelTotalBits) - *v) >> kNumMoveBits
		d.Range = bound

		// Normalize
		if d.Range < kTopValue {
			b, err := d.inStream.ReadByte()
			if err != nil {
				return 0, err
			}

			d.Range <<= 8
			d.Code = (d.Code << 8) | uint32(b)

			return 0, nil
		} else {
			return 0, nil
		}
	} else {
		*v -= *v >> kNumMoveBits
		d.Code -= bound
		d.Range -= bound

		// Normalize
		if d.Range < kTopValue {
			b, err := d.inStream.ReadByte()
			if err != nil {
				return 0, err
			}

			d.Range <<= 8
			d.Code = (d.Code << 8) | uint32(b)

			return 1, nil
		} else {
			return 1, nil
		}
	}
}

func (d *rangeDecoder) DecodeDirectBits(numBits int) (uint32, error) {
	var res uint32
	rang := d.Range
	code := d.Code

	for ; numBits > 0; numBits-- {
		rang >>= 1
		code -= rang
		t := 0 - (code >> 31)
		code += rang & t

		if code == rang {
			d.Corrupted = true
		}

		res <<= 1
		res += t + 1

		// Normalize
		if rang < kTopValue {
			b, err := d.inStream.ReadByte()
			if err != nil {
				return 0, err
			}

			rang <<= 8
			code = (code << 8) | uint32(b)
		}
	}

	d.Range = rang
	d.Code = code

	return res, nil
}
