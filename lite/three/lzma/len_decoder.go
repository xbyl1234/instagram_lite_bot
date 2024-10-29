package lzma

type lenDecoder struct {
	choice  prob
	choice2 prob

	lowCoder  [1 << kNumPosBitsMax][1 << lenLowCoderNumBits]prob
	midCoder  [1 << kNumPosBitsMax][1 << lenMidCoderNumBits]prob
	highCoder [1 << lenHighCoderNumBits]prob
}

func newLenDecoder() *lenDecoder {
	d := &lenDecoder{
		choice:  probInitVal,
		choice2: probInitVal,
	}

	d.Reset()

	return d
}

func (d *lenDecoder) Reset() {
	d.choice = probInitVal
	d.choice2 = probInitVal
	initProbs(d.highCoder[:])

	for i := 0; i < len(d.lowCoder); i++ {
		initProbs(d.lowCoder[i][:])
		initProbs(d.midCoder[i][:])
	}
}

func (d *lenDecoder) Decode(rc *rangeDecoder, posState uint32) (uint32, error) {
	bit, err := rc.DecodeBit(&d.choice)
	if err != nil {
		return 0, err
	}
	if bit == 0 {
		return BitTreeDecode(d.lowCoder[posState][:], lenLowCoderNumBits, rc)
	}

	bit, err = rc.DecodeBit(&d.choice2)
	if err != nil {
		return 0, err
	}
	if bit == 0 {
		bit, err = BitTreeDecode(d.midCoder[posState][:], lenMidCoderNumBits, rc)
		if err != nil {
			return 0, err
		}
		return 8 + bit, nil
	}

	bit, err = BitTreeDecode(d.highCoder[:], lenHighCoderNumBits, rc)
	if err != nil {
		return 0, err
	}
	return 16 + bit, nil
}
