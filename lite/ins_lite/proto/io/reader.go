package io

import (
	"bytes"
	"errors"
	"io"
	"math"
)

type Reader struct {
	BufferReader
	reader io.Reader
	buffer *bytes.Buffer
	size   int
}

func CreateReader(buff io.Reader) *Reader {
	return &Reader{
		reader: buff,
	}
}

func CreateReaderBuffer(buff []byte) *Reader {
	b := bytes.NewBuffer(buff)
	return &Reader{
		reader: b,
		buffer: b,
		size:   len(buff),
	}
}

func (this *Reader) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (this *Reader) Offset() int {
	if this.buffer == nil {
		panic("not ReaderBuffer")
	}
	return this.size - this.buffer.Len()
}

func (this *Reader) EOF() bool {
	if this.buffer == nil {
		panic("not ReaderBuffer")
	}
	return this.buffer.Len() == 0
}

func (this *Reader) PeekRemain() []byte {
	if this.buffer == nil {
		panic("not ReaderBuffer")
	}
	return this.buffer.Bytes()
}

func (this *Reader) ReadRemain() []byte {
	all, err := io.ReadAll(this.reader)
	this.checkError(err)
	return all
}

func (this *Reader) ReadBytes(l uint32) []byte {
	b := make([]byte, l)
	cur := uint32(0)
	for true {
		count, err := this.reader.Read(b[cur:l])
		this.checkError(err)
		cur += uint32(count)
		if cur == l {
			break
		}
	}
	//if count != l {
	//	panic(errors.New(fmt.Sprintf("read byte size is %d, not need %d", count, l)))
	//}
	return b
}

func (this *Reader) ReadByte() byte {
	b := make([]byte, 1)
	count, err := this.reader.Read(b)
	this.checkError(err)
	if count != 1 {
		panic(errors.New("read byte size is 0"))
	}
	return b[0]
}

func (this *Reader) ReadFloat32() float32 {
	return math.Float32frombits(uint32(this.ReadInt()))
}

func (this *Reader) ReadFloat64() float64 {
	return math.Float64frombits(uint64(this.ReadLong()))
}

func (this *Reader) ReadInt() int32 {
	readBytes := this.ReadBytes(4)
	var ret int32 = 0
	ret |= int32(readBytes[0]) << 24
	ret |= int32(readBytes[1]) << 16
	ret |= int32(readBytes[2]) << 8
	ret |= int32(readBytes[3])
	return ret
}

func (this *Reader) ReadLong() int64 {
	readBytes := this.ReadBytes(8)
	var ret int64 = 0
	ret |= int64(readBytes[0]) << 56
	ret |= int64(readBytes[1]) << 48
	ret |= int64(readBytes[2]) << 40
	ret |= int64(readBytes[3]) << 32
	ret |= int64(readBytes[4]) << 24
	ret |= int64(readBytes[5]) << 16
	ret |= int64(readBytes[6]) << 8
	ret |= int64(readBytes[7])
	return ret
}

func (this *Reader) ReadUShort() uint16 {
	readBytes := this.ReadBytes(2)
	var ret uint16
	ret = uint16(readBytes[1]) | uint16(readBytes[0])<<8
	return ret
}

func (this *Reader) ReadShort() int16 {
	return int16(this.ReadUShort())
}

func (this *Reader) ReadString() string {
	l := this.ReadShort()
	if l == 0 {
		return ""
	}
	return string(this.ReadBytes(uint32(l)))
}

func (this *Reader) ReadVarUInt32() uint32 {
	var l uint32 = 0
	var reader uint32
	for i := 0; i < 0x20; i += 7 {
		reader = uint32(this.ReadByte())
		l |= (reader & 0x7F) << i
		if (reader & 0x80) == 0 {
			if i == 28 && (reader&0x7F) > 7 {
				this.checkError(errors.New("readVarInt32 error"))
			}
			return l
		}
	}
	this.checkError(errors.New("readVarInt32 error"))
	return 0
}

func (this *Reader) ReadVarInt32() int32 {
	var reader int32
	reader = int32(this.ReadByte())
	if reader == 0 {
		return 0
	}
	var flag int32
	if (reader & 0x40) == 0 {
		flag = 0
	} else {
		flag = 1
	}

	var result int32
	result = reader & 0x3F
	if (reader & 0x80) == 0 {
		if flag == 0 {
			return result
		} else {
			return -result
		}
	}

	for index := 6; index < 32; index += 7 {
		reader = int32(this.ReadByte())
		result |= (reader & 0x7F) << index
		if (reader & 0x80) == 0 {
			if index == 27 {
				if (reader & 0x7F) <= 16 {
					if (reader & 0x7F) <= 15 {
						if flag == 0 {
							return result
						} else {
							return -result
						}
					}
					if flag == 0 {
						this.checkError(errors.New("malformed signed varint32"))
					}
					return -result
				}
				this.checkError(errors.New("malformed signed varint32"))
			}
			if flag == 0 {
				return result
			} else {
				return -result
			}
		}
	}
	this.checkError(errors.New("malformed signed varint32"))
	return 0
}

func (this *Reader) ReadVarUInt64(isBit63 bool) uint64 {
	var result uint64
	var idx = 0
	var bitCount = 0
	if isBit63 {
		bitCount = 63
	} else {
		bitCount = 64
	}
	for {
		r := this.ReadByte()
		result |= uint64(r&0x7F) << idx
		if (r & 128) == 0 {
			if isBit63 {
				return result
			} else {
				return result>>1 ^ -(1 & result)
			}
		}
		idx += 7
		if idx >= bitCount {
			break
		}
	}
	this.checkError(errors.New("malformed signed varint64"))
	return 0
}

//func (this *Reader) ReadStringVarLen() string {
//	l := int(this.ReadVarInt32())
//	return string(this.ReadBytes(l - 1))
//}
//
//func (this *Reader) ReadStrArrayVarLen() []string {
//	l := int(this.ReadVarInt32())
//	if l == 0 {
//		return nil
//	}
//	ret := make([]string, l-1)
//	for i := 0; i < l-1; i++ {
//		ret[i] = this.ReadString()
//	}
//	return ret
//}
//
//func (this *Reader) ReadStrArray() []string {
//	l := int(this.ReadInt())
//	if l == 0 {
//		return nil
//	}
//	ret := make([]string, l-1)
//	for i := 0; i < l-1; i++ {
//		ret[i] = this.ReadString()
//	}
//	return ret
//}
