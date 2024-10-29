package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	math_rand "math/rand"
	"reflect"
)

func GenString(charSet string, length int) string {
	by := make([]byte, length)
	//math_rand.Seed(time.Now().UnixNano())
	for index := 0; index < length; index++ {
		by[index] = charSet[math_rand.Intn(len(charSet))]
	}
	return B2s(by)
}

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GenUUID() string {
	uuid, err := NewUUID()
	if err != nil {
		return "cb479ee7-a50d-49e7-8b7b-60cc1a105e22" // default value when error occurred
	}
	return uuid
}

func RandIndex(obj interface{}) int {
	if reflect.ValueOf(obj).Len() == 1 {
		return 0
	}
	return GenNumber(0, reflect.ValueOf(obj).Len()-1)
}

func GenBytes(count int) []byte {
	r := make([]byte, count)
	for i := 0; i < count; i++ {
		r[i] = byte(GenNumber(0, 0xff))
	}
	return r
}

func GenNumber(min, max int) int {
	//math_rand.Seed(time.Now().Unix())
	if min == max {
		return max
	}
	return math_rand.Intn(max-min) + min
}

func GenFloat(min, max float64) float64 {
	return min + math_rand.Float64()*(max-min)
}

func ChoseOne[Type any](array []Type) Type {
	return array[GenNumber(0, len(array)-1)]
}

func VariantString(s string, threshold int) string {
	var ret = ""
	for _, ch := range s {
		if ch >= 'a' && ch <= 'z' {
			CH := ch
			if CH >= 'a' && CH <= 'z' {
				CH -= 'a' - 'A'
			}
			v := variantMild[CH][0]
			if len(v) == 0 || GenNumber(0, 100) > threshold {
				ret += string(ch)
			} else {
				ret += string(v[RandIndex(v)])
			}
		} else if ch >= 'A' && ch <= 'Z' {
			v := variantMild[ch][1]
			if len(v) == 0 || GenNumber(0, 100) > threshold {
				ret += string(ch)
			} else {
				ret += string(v[RandIndex(v)])
			}
		} else {
			ret += string(ch)
		}
	}
	return ret
}
