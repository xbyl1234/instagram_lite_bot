package utils

import (
	"CentralizedControl/common/log"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/klauspost/compress/gzip"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
	"unsafe"
)

const (
	volatileSeed   = "12345"
	CharSet_abc    = "abcdefghijklmnopqrstuvwxyz"
	CharSet_ABC    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharSet_123    = "0123456789"
	CharSet_16_Num = "0123456789abcdef"
	CharSet_All    = CharSet_abc + CharSet_ABC + CharSet_123
)

func LoadJsonFile(path string, ret interface{}) error {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(by, ret)
	return err
}

func Dumps(path string, obj interface{}) error {
	by, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write(by)
	if err != nil {
		return err
	}
	_ = file.Close()
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func B2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetRMidByte(src []byte, start, end string) string {
	e := bytes.LastIndex(src, []byte(start))
	if e == -1 {
		return ""
	}
	s := bytes.LastIndex(src[:e], []byte(end))
	if s == -1 {
		return ""
	}
	return string(src[s+len(end) : e])
}

func GetRStartMidByte(src []byte, offset, rStart, rEnd string) string {
	p := bytes.LastIndex(src, []byte(offset))
	if p == -1 {
		return ""
	}
	return GetRMidByte(src[:p], rStart, rEnd)
}

func GetLMidByte(src []byte, start, end string) string {
	ps := bytes.Index(src, []byte(start))
	if ps == -1 {
		return ""
	}
	pe := bytes.Index(src[ps+len(start):], []byte(end))
	if pe == -1 {
		return ""
	}
	return string(src[ps+len(start) : ps+len(start)+pe])
}

func GetLStartMidByte(src []byte, offset, lStart, lEnd string) string {
	p := bytes.Index(src, []byte(offset))
	if p == -1 {
		return ""
	}
	return GetLMidByte(src[p+len(offset):], lStart, lEnd)
}

func GetMidString(src, start, end string) string {
	ps := strings.Index(src, start)
	if ps == -1 {
		return ""
	}
	pe := strings.Index(src[ps+len(start):], end)
	if pe == -1 {
		return ""
	}
	return src[ps+len(start) : ps+len(start)+pe]
}

func Json2String(params map[string]string) (string, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return B2s(data), nil
}

func GZipCompress(in []byte) []byte {
	var input bytes.Buffer
	_gzip := gzip.NewWriter(&input)
	_gzip.Write(in)
	_gzip.Close()
	return input.Bytes()
}

func GZipDecompress(in []byte) ([]byte, error) {
	input := bytes.NewReader(in)
	_gzip, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	_, err = _gzip.WriteTo(&output)
	return output.Bytes(), err
}

func Base64Encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func DecodeBase64(s string) ([]byte, error) {
	miss := len(s) % 4
	if miss != 0 {
		s = s + strings.Repeat("=", 4-miss)
	}
	return base64.StdEncoding.DecodeString(s)
}

func GetCode(msg string) string {
	var index = 0
	find := false
	for index = range msg {
		if msg[index] >= '0' && msg[index] <= '9' {
			find = true
			break
		}
	}
	if find {
		code := strings.ReplaceAll(msg[index:index+7], " ", "")
		if len(code) != 6 {
			return ""
		}
		return code
	} else {
		return ""
	}
}

type Encoding int

const (
	EscapeEncodePath Encoding = 1 + iota
	EscapeEncodePathSegment
	EscapeEncodeHost
	EscapeEncodeZone
	EscapeEncodeUserPassword
	EscapeEncodeQueryComponent
	EscapeEncodeFragment
	EscapeEncodeNone
)

//encodePath: 用于对URL的路径部分进行编码。在这种模式下，除了一些特殊字符（如/）之外，所有的非字母数字字符都会被编码。
//encodePathSegment: 用于对URL的路径段进行编码。在这种模式下，所有的非字母数字字符都会被编码。
//encodeHost: 用于对URL的主机部分进行编码。在这种模式下，除了一些特殊字符（如.）之外，所有的非字母数字字符都会被编码。
//encodeZone: 用于对URL的区域部分进行编码。在这种模式下，除了一些特殊字符（如-和_）之外，所有的非字母数字字符都会被编码。
//encodeUserPassword: 用于对URL的用户密码部分进行编码。在这种模式下，除了一些特殊字符（如:）之外，所有的非字母数字字符都会被编码。
//encodeQueryComponent: 用于对URL的查询部分进行编码。在这种模式下，除了一些特殊字符（如=和&）之外，所有的非字母数字字符都会被编码。
//encodeFragment: 用于对URL的片段部分进行编码。在这种模式下，所有的非字母数字字符都会被编码。
//EncodeQueryPath函数调用了Escape函数，并且传入了encodePath作为mode参数。这意味着EncodeQueryPath函数会对URL的路径部分进行编码，除了一些特殊字符（如/）之外，所有的非字母数字字符都会被编码。

const upperhex = "0123456789ABCDEF"

func shouldEscape(c byte, mode Encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == EscapeEncodeHost || mode == EscapeEncodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case EscapeEncodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case EscapeEncodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case EscapeEncodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case EscapeEncodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case EscapeEncodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	if mode == EscapeEncodeFragment {
		// RFC 3986 §2.2 allows not escaping sub-delims. A subset of sub-delims are
		// included in reserved from RFC 2396 §2.2. The remaining sub-delims do not
		// need to be escaped. To minimize potential breakage, we apply two restrictions:
		// (1) we always escape sub-delims outside of the fragment, and (2) we always
		// escape single quote to avoid breaking callers that had previously assumed that
		// single quotes would be escaped. See issue #19917.
		switch c {
		case '!', '(', ')', '*':
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

func Escape(s string, mode Encoding) string {
	if mode == EscapeEncodeNone {
		return s
	}
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			//if c == ' ' && mode == EscapeEncodeQueryComponent {
			//	spaceCount++
			//} else {
			hexCount++
			//}
		}
	}
	_ = spaceCount
	//if spaceCount == 0 && hexCount == 0 {
	if hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	if hexCount == 0 {
		copy(t, s)
		for i := 0; i < len(s); i++ {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		//case c == ' ' && mode == EscapeEncodeQueryComponent:
		//	t[j] = '+'
		//	j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func EncodeQueryPath(s string) string {
	return Escape(s, EscapeEncodeQueryComponent)
}

func EncodeQueryMap(v map[string][]string) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := Escape(k, EscapeEncodePath)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(Escape(v, EscapeEncodePath))
		}
	}
	return buf.String()
}

const TimeLayout = "2006-01-02 15:04:05"

func GetNewYorkTimeString() string {
	location, _ := time.LoadLocation("America/New_York")
	return time.Now().In(location).Format(TimeLayout)
}

func GetShanghaiTimeString() string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(location).Format(TimeLayout)
}

var variantMild = map[int32][][]rune{
	'A': [][]rune{[]rune("ɑàáâãäåāáǎàа"), []rune("ÅÄÃÂÁÀАΑ")},
	'B': [][]rune{[]rune("ЬЪ"), []rune("вВΒ")},
	'C': [][]rune{[]rune("çс"), []rune("ÇС")},
	'D': [][]rune{[]rune(""), []rune("")},
	'E': [][]rune{[]rune("ēéěèèéêëêеё"), []rune("ËÊÉÈ")},
	'F': [][]rune{[]rune(""), []rune("")},
	'G': [][]rune{[]rune(""), []rune("")},
	'H': [][]rune{[]rune(""), []rune("нН")},
	'I': [][]rune{[]rune("¡ⅰīíǐììíîï"), []rune("｜ⅠÍÌÏ")},
	'J': [][]rune{[]rune(""), []rune("")},
	'K': [][]rune{[]rune("кΚК"), []rune("")},
	'L': [][]rune{[]rune(""), []rune("")},
	'M': [][]rune{[]rune("м"), []rune("М")},
	'N': [][]rune{[]rune("ńňп"), []rune("ий")},
	'O': [][]rune{[]rune("ооòóõôöōóǒò"), []rune("○ÒÓÔÕÖΟО")},
	'P': [][]rune{[]rune("р"), []rune("Р")},
	'Q': [][]rune{[]rune(""), []rune("")},
	'R': [][]rune{[]rune(""), []rune("")},
	'S': [][]rune{[]rune("š"), []rune("")},
	'T': [][]rune{[]rune("т"), []rune("Т")},
	'U': [][]rune{[]rune("ùǔúūǜǚǘǖùúûüü"), []rune("")},
	'V': [][]rune{[]rune("ⅴ"), []rune("ⅤⅤ")},
	'W': [][]rune{[]rune(""), []rune("")},
	'X': [][]rune{[]rune("×ΧХх"), []rune("ⅹ✖")},
	'Y': [][]rune{[]rune("ýÿу"), []rune("")},
	'Z': [][]rune{[]rune(""), []rune("")},
}

func WriteFile(_path string, data []byte) error {
	_path = strings.ReplaceAll(_path, "\\", "/")
	err := os.MkdirAll(path.Dir(_path), 0777)
	if err != nil {
		return err
	}
	return os.WriteFile(_path, data, 0777)
}

func GetTimezoneOffset(timezone string) int {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Warn("time zone %s is error!", timezone)
		return 0
	}
	_, offset := time.Now().In(loc).Zone()
	return offset
}

func GetDefault[ValueType any](value ValueType, exist bool, _default ValueType) any {
	if exist {
		return value
	} else {
		return _default
	}
}

func DecodeHex(data string) []byte {
	data = strings.ReplaceAll(data, " ", "")
	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\r", "")
	data = strings.ReplaceAll(data, "\t", "")
	decodeString, _ := hex.DecodeString(data)
	return decodeString
}

func Sha256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
