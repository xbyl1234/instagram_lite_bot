package utils

import (
	"fmt"
	"io"
	math_rand "math/rand"
	"reflect"
)

type RandTool struct {
	rand *math_rand.Rand
}

func CreateRandTool(seed int64) *RandTool {
	r := &RandTool{
		rand: math_rand.New(math_rand.NewSource(seed)),
	}
	return r
}

func (this *RandTool) GenString(charSet string, length int) string {
	by := make([]byte, length)
	for index := 0; index < length; index++ {
		by[index] = charSet[this.rand.Intn(len(charSet))]
	}
	return B2s(by)
}

func (this *RandTool) GenNumber(min, max int) int {
	if min == max {
		return max
	}
	return this.rand.Intn(max-min) + min
}

func (this *RandTool) GenFloat(min, max float64) float64 {
	return min + this.rand.Float64()*(max-min)
}

func (this *RandTool) RandIndex(obj interface{}) int {
	if reflect.ValueOf(obj).Len() == 1 {
		return 0
	}
	return this.GenNumber(0, reflect.ValueOf(obj).Len()-1)
}

func (this *RandTool) GenBytes(count int) []byte {
	r := make([]byte, count)
	for i := 0; i < count; i++ {
		r[i] = byte(this.GenNumber(0, 0xff))
	}
	return r
}

func (this *RandTool) NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(this.rand, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func (this *RandTool) GenUUID() string {
	uuid, err := this.NewUUID()
	if err != nil {
		return "cb479ee7-a50d-49e7-8b7b-60cc1a105e22" // default value when error occurred
	}
	return uuid
}

func (this *RandTool) VariantString(s string, threshold int) string {
	var ret = ""
	for _, ch := range s {
		if ch >= 'a' && ch <= 'z' {
			CH := ch
			if CH >= 'a' && CH <= 'z' {
				CH -= 'a' - 'A'
			}
			v := variantMild[CH][0]
			if len(v) == 0 || this.GenNumber(0, 100) > threshold {
				ret += string(ch)
			} else {
				ret += string(v[this.RandIndex(v)])
			}
		} else if ch >= 'A' && ch <= 'Z' {
			v := variantMild[ch][1]
			if len(v) == 0 || this.GenNumber(0, 100) > threshold {
				ret += string(ch)
			} else {
				ret += string(v[this.RandIndex(v)])
			}
		} else {
			ret += string(ch)
		}
	}
	return ret
}

func ChoseOne2[Type any](this *RandTool, array []Type) Type {
	return array[this.GenNumber(0, len(array)-1)]
}
