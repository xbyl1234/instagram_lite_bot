package io

import (
	"math"
)

type Writer struct {
	BufferWriter
	buffer []byte
	cur    int
}

func CreateWriter(size int) *Writer {
	if size == 0 {
		size = 256
	}
	return &Writer{
		buffer: make([]byte, size),
		cur:    0,
	}
}

func (this *Writer) ensureSize(l int) {
	if this.cur+l > len(this.buffer) {
		newBuff := make([]byte, len(this.buffer)+l+256)
		copy(newBuff, this.buffer)
		this.buffer = newBuff
	}
}

func (this *Writer) GetBytes() []byte {
	return this.buffer[:this.cur]
}

func (this *Writer) WriteBytes(d []byte) {
	this.ensureSize(len(d))
	copy(this.buffer[this.cur:], d)
	this.cur += len(d)
}

func (this *Writer) WriteByte(d byte) {
	this.ensureSize(1)
	this.buffer[this.cur] = d
	this.cur++
}

func (this *Writer) WriteInt(d int32) {
	this.ensureSize(4)
	this.buffer[this.cur] = (byte)(d >> 24)
	this.buffer[this.cur+1] = (byte)(d >> 16)
	this.buffer[this.cur+2] = (byte)(d >> 8)
	this.buffer[this.cur+3] = (byte)(d)
	this.cur += 4
}

func (this *Writer) WriteLong(d int64) {
	this.ensureSize(8)
	this.buffer[this.cur] = (byte)(d >> 56)
	this.buffer[this.cur+1] = (byte)(d >> 48)
	this.buffer[this.cur+2] = (byte)(d >> 40)
	this.buffer[this.cur+3] = (byte)(d >> 32)
	this.buffer[this.cur+4] = (byte)(d >> 24)
	this.buffer[this.cur+5] = (byte)(d >> 16)
	this.buffer[this.cur+6] = (byte)(d >> 8)
	this.buffer[this.cur+7] = (byte)(d)
	this.cur += 8
}

func (this *Writer) WriteFloat32(d float32) {
	this.WriteInt(int32(math.Float32bits(d)))
}

func (this *Writer) WriteFloat64(d float64) {
	this.WriteLong(int64(math.Float64bits(d)))
}

func (this *Writer) WriteShort(d int16) {
	this.ensureSize(2)
	this.buffer[this.cur] = byte(d >> 8)
	this.cur++
	this.buffer[this.cur] = byte(d & 0xff)
	this.cur++
}

func (this *Writer) WriteString(d string) {
	this.WriteShort(int16(len(d)))
	if len(d) != 0 {
		this.WriteBytes([]byte(d))
	}
}

func (this *Writer) WriteVarUInt32(l uint32) {
	if l == 0 {
		this.WriteByte(0)
		return
	}
	var write byte
	for {
		write = byte(l & 0x7F)
		l >>= 7
		if l > 0 {
			write |= 0x80
		}
		this.WriteByte(write)
		if l <= 0 {
			break
		}
	}
}

func (this *Writer) WriteVarInt32(d int32) {
	if d == 0 {
		this.WriteByte(0)
		return
	}

	abs := int64(math.Abs(float64(d)))
	writer := byte(0x3f & abs)
	remain := int32(abs >> 6)
	if remain > 0 {
		writer |= 0x80
	}
	if d < 0 {
		writer |= 0x40
	}

	for {
		this.WriteByte(writer)
		if remain <= 0 {
			return
		}
		writer = byte(remain&0x7f) | 0x80
		remain >>= 7
	}
}

func (this *Writer) WriteVarInt64(v int64, isSigned bool) {
	if !isSigned && v < 0 {
		panic("Received a negative varint64 value")
	}
	if v == 0 {
		this.WriteByte(0)
		return
	}
	if isSigned {
		v = v<<1 ^ (v >> 63)
	}
	for v != 0 {
		w := (byte)((int)(0x7F & v))
		v = v >> 7
		if v > 0 {
			w = w | 128
		}
		this.WriteByte(w)
	}
}

func GetVarUInt32Len(l uint32) int {
	if l == 0 {
		return 1
	}
	var varIntLen = 0
	var write byte
	for {
		write = byte(l & 0x7F)
		l >>= 7
		if l > 0 {
			write |= 0x80
		}
		varIntLen += 1
		if l <= 0 {
			break
		}
	}
	return varIntLen
}
