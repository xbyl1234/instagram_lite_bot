package lzma

import "unsafe"

type bitTreeDecoder struct {
	probs   []prob
	numBits int
}

func newBitTreeDecoder(numBits int) *bitTreeDecoder {
	d := &bitTreeDecoder{
		numBits: numBits,
		probs:   make([]prob, uint32(1)<<numBits),
	}
	d.Reset()

	return d
}

func (d *bitTreeDecoder) Reset() {
	initProbs(d.probs)
}

func (d *bitTreeDecoder) Decode(rc *rangeDecoder) (uint32, error) {
	return BitTreeDecode(d.probs, d.numBits, rc)
}

func BitTreeDecode(probs []prob, numBits int, rc *rangeDecoder) (uint32, error) {
	m := uint32(1)

	rang := rc.Range
	code := rc.Code
	probsPtr := uintptr(unsafe.Pointer(&probs[0]))

	for i := 0; i < numBits; i++ {
		probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
		// rc.DecodeBit begin
		bound := (rang >> kNumBitModelTotalBits) * uint32(*probPtr)

		if code < bound {
			*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
			rang = bound
			m <<= 1

			// Normalize
			if rang < kTopValue {
				b, err := rc.inStream.ReadByte()
				if err != nil {
					return 0, err
				}

				rang <<= 8
				code = (code << 8) | uint32(b)
			}
		} else {
			*probPtr -= *probPtr >> kNumMoveBits
			code -= bound
			rang -= bound
			m = (m << 1) + 1

			// Normalize
			if rang < kTopValue {
				b, err := rc.inStream.ReadByte()
				if err != nil {
					return 0, err
				}

				rang <<= 8
				code = (code << 8) | uint32(b)
			}
		}
		// rc.DecodeBit end
	}

	rc.Range = rang
	rc.Code = code

	return m - (uint32(1) << numBits), nil
}

func (d *bitTreeDecoder) ReverseDecode(rc *rangeDecoder) (uint32, error) {
	return BitTreeReverseDecode(d.probs, d.numBits, rc)
}

func BitTreeReverseDecode(probs []prob, numBits int, rc *rangeDecoder) (uint32, error) {
	rang := rc.Range
	code := rc.Code

	m := uint32(1)
	symbol := uint32(0)
	probsPtr := uintptr(unsafe.Pointer(&probs[0]))

	for i := 0; i < numBits; i++ {
		probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
		// rc.DecodeBit begin
		bound := (rang >> kNumBitModelTotalBits) * uint32(*probPtr)

		if code < bound {
			*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
			rang = bound
			m <<= 1
			symbol |= 0 << i

			// Normalize
			if rang < kTopValue {
				b, err := rc.inStream.ReadByte()
				if err != nil {
					return 0, err
				}

				rang <<= 8
				code = (code << 8) | uint32(b)
			}
		} else {
			*probPtr -= *probPtr >> kNumMoveBits
			code -= bound
			rang -= bound
			m = (m << 1) | 1
			symbol |= 1 << i

			// Normalize
			if rang < kTopValue {
				b, err := rc.inStream.ReadByte()
				if err != nil {
					return 0, err
				}

				rang <<= 8
				code = (code << 8) | uint32(b)
			}
		}
		// rc.DecodeBit end
	}

	rc.Range = rang
	rc.Code = code

	return symbol, nil
}
