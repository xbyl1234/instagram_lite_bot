//go:build !noasm || !appengine
// +build !noasm !appengine

// Code generated by asm2asm, DO NOT EDIT.

package avx2

import (
	"github.com/bytedance/sonic/loader"
)

const (
	_entry__f32toa                 = 34624
	_entry__f64toa                 = 320
	_entry__format_significand     = 38736
	_entry__format_integer         = 3168
	_entry__fsm_exec               = 21072
	_entry__advance_ns             = 16928
	_entry__advance_string         = 17664
	_entry__advance_string_default = 40160
	_entry__do_skip_number         = 23696
	_entry__get_by_path            = 28864
	_entry__skip_one_fast          = 25936
	_entry__html_escape            = 10560
	_entry__i64toa                 = 3600
	_entry__u64toa                 = 3712
	_entry__lspace                 = 64
	_entry__quote                  = 5104
	_entry__skip_array             = 21024
	_entry__skip_number            = 25392
	_entry__skip_object            = 23088
	_entry__skip_one               = 25536
	_entry__unquote                = 7888
	_entry__validate_one           = 25584
	_entry__validate_utf8          = 31040
	_entry__validate_utf8_fast     = 31984
	_entry__value                  = 15376
	_entry__vnumber                = 18800
	_entry__atof_eisel_lemire64    = 12624
	_entry__atof_native            = 14768
	_entry__decimal_to_f64         = 13056
	_entry__right_shift            = 39696
	_entry__left_shift             = 39200
	_entry__vsigned                = 20352
	_entry__vstring                = 17424
	_entry__vunsigned              = 20672
)

const (
	_stack__f32toa                 = 48
	_stack__f64toa                 = 80
	_stack__format_significand     = 24
	_stack__format_integer         = 16
	_stack__fsm_exec               = 144
	_stack__advance_ns             = 8
	_stack__advance_string         = 56
	_stack__advance_string_default = 48
	_stack__do_skip_number         = 48
	_stack__get_by_path            = 272
	_stack__skip_one_fast          = 184
	_stack__html_escape            = 72
	_stack__i64toa                 = 16
	_stack__u64toa                 = 8
	_stack__lspace                 = 8
	_stack__quote                  = 56
	_stack__skip_array             = 152
	_stack__skip_number            = 88
	_stack__skip_object            = 152
	_stack__skip_one               = 152
	_stack__unquote                = 72
	_stack__validate_one           = 152
	_stack__validate_utf8          = 48
	_stack__validate_utf8_fast     = 176
	_stack__value                  = 328
	_stack__vnumber                = 240
	_stack__atof_eisel_lemire64    = 32
	_stack__atof_native            = 136
	_stack__decimal_to_f64         = 80
	_stack__right_shift            = 8
	_stack__left_shift             = 24
	_stack__vsigned                = 16
	_stack__vstring                = 112
	_stack__vunsigned              = 8
)

const (
	_size__f32toa                 = 3392
	_size__f64toa                 = 2848
	_size__format_significand     = 464
	_size__format_integer         = 432
	_size__fsm_exec               = 1468
	_size__advance_ns             = 496
	_size__advance_string         = 1088
	_size__advance_string_default = 768
	_size__do_skip_number         = 1360
	_size__get_by_path            = 2176
	_size__skip_one_fast          = 2428
	_size__html_escape            = 2064
	_size__i64toa                 = 48
	_size__u64toa                 = 1248
	_size__lspace                 = 224
	_size__quote                  = 2736
	_size__skip_array             = 48
	_size__skip_number            = 144
	_size__skip_object            = 48
	_size__skip_one               = 48
	_size__unquote                = 2480
	_size__validate_one           = 48
	_size__validate_utf8          = 672
	_size__validate_utf8_fast     = 2608
	_size__value                  = 1004
	_size__vnumber                = 1552
	_size__atof_eisel_lemire64    = 368
	_size__atof_native            = 608
	_size__decimal_to_f64         = 1712
	_size__right_shift            = 400
	_size__left_shift             = 496
	_size__vsigned                = 320
	_size__vstring                = 144
	_size__vunsigned              = 336
)

var (
	_pcsp__f32toa = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{3350, 48},
		{3351, 40},
		{3353, 32},
		{3355, 24},
		{3357, 16},
		{3359, 8},
		{3363, 0},
		{3385, 48},
	}
	_pcsp__f64toa = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{2788, 56},
		{2792, 48},
		{2793, 40},
		{2795, 32},
		{2797, 24},
		{2799, 16},
		{2801, 8},
		{2805, 0},
		{2843, 56},
	}
	_pcsp__format_significand = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{452, 24},
		{453, 16},
		{455, 8},
		{457, 0},
	}
	_pcsp__format_integer = [][2]uint32{
		{1, 0},
		{4, 8},
		{412, 16},
		{413, 8},
		{414, 0},
		{423, 16},
		{424, 8},
		{426, 0},
	}
	_pcsp__fsm_exec = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{1157, 88},
		{1161, 48},
		{1162, 40},
		{1164, 32},
		{1166, 24},
		{1168, 16},
		{1170, 8},
		{1171, 0},
		{1468, 88},
	}
	_pcsp__advance_ns = [][2]uint32{
		{1, 0},
		{453, 8},
		{457, 0},
		{481, 8},
		{486, 0},
	}
	_pcsp__advance_string = [][2]uint32{
		{14, 0},
		{18, 8},
		{20, 16},
		{22, 24},
		{24, 32},
		{26, 40},
		{27, 48},
		{433, 56},
		{437, 48},
		{438, 40},
		{440, 32},
		{442, 24},
		{444, 16},
		{446, 8},
		{450, 0},
		{1078, 56},
	}
	_pcsp__advance_string_default = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{332, 48},
		{333, 40},
		{335, 32},
		{337, 24},
		{339, 16},
		{341, 8},
		{345, 0},
		{757, 48},
	}
	_pcsp__do_skip_number = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{1274, 48},
		{1275, 40},
		{1277, 32},
		{1279, 24},
		{1281, 16},
		{1283, 8},
		{1287, 0},
		{1360, 48},
	}
	_pcsp__get_by_path = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{2049, 88},
		{2053, 48},
		{2054, 40},
		{2056, 32},
		{2058, 24},
		{2060, 16},
		{2062, 8},
		{2063, 0},
		{2170, 88},
	}
	_pcsp__skip_one_fast = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{418, 176},
		{419, 168},
		{421, 160},
		{423, 152},
		{425, 144},
		{427, 136},
		{431, 128},
		{2428, 176},
	}
	_pcsp__html_escape = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{2045, 72},
		{2049, 48},
		{2050, 40},
		{2052, 32},
		{2054, 24},
		{2056, 16},
		{2058, 8},
		{2063, 0},
	}
	_pcsp__i64toa = [][2]uint32{
		{14, 0},
		{34, 8},
		{36, 0},
	}
	_pcsp__u64toa = [][2]uint32{
		{1, 0},
		{161, 8},
		{162, 0},
		{457, 8},
		{458, 0},
		{758, 8},
		{759, 0},
		{1225, 8},
		{1227, 0},
	}
	_pcsp__lspace = [][2]uint32{
		{1, 0},
		{184, 8},
		{188, 0},
		{204, 8},
		{208, 0},
		{215, 8},
		{220, 0},
	}
	_pcsp__quote = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{2687, 56},
		{2691, 48},
		{2692, 40},
		{2694, 32},
		{2696, 24},
		{2698, 16},
		{2700, 8},
		{2704, 0},
		{2731, 56},
	}
	_pcsp__skip_array = [][2]uint32{
		{1, 0},
		{28, 8},
		{34, 0},
	}
	_pcsp__skip_number = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{100, 40},
		{101, 32},
		{103, 24},
		{105, 16},
		{107, 8},
		{108, 0},
		{139, 40},
	}
	_pcsp__skip_object = [][2]uint32{
		{1, 0},
		{28, 8},
		{34, 0},
	}
	_pcsp__skip_one = [][2]uint32{
		{1, 0},
		{30, 8},
		{36, 0},
	}
	_pcsp__unquote = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{79, 72},
		{83, 48},
		{84, 40},
		{86, 32},
		{88, 24},
		{90, 16},
		{92, 8},
		{96, 0},
		{2464, 72},
	}
	_pcsp__validate_one = [][2]uint32{
		{1, 0},
		{35, 8},
		{41, 0},
	}
	_pcsp__validate_utf8 = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{11, 40},
		{623, 48},
		{627, 40},
		{628, 32},
		{630, 24},
		{632, 16},
		{634, 8},
		{635, 0},
		{666, 48},
	}
	_pcsp__validate_utf8_fast = [][2]uint32{
		{1, 0},
		{4, 8},
		{5, 16},
		{1738, 176},
		{1739, 168},
		{1743, 160},
		{2018, 176},
		{2019, 168},
		{2023, 160},
		{2600, 176},
	}
	_pcsp__value = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{988, 88},
		{992, 48},
		{993, 40},
		{995, 32},
		{997, 24},
		{999, 16},
		{1001, 8},
		{1004, 0},
	}
	_pcsp__vnumber = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{803, 104},
		{807, 48},
		{808, 40},
		{810, 32},
		{812, 24},
		{814, 16},
		{816, 8},
		{817, 0},
		{1547, 104},
	}
	_pcsp__atof_eisel_lemire64 = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{292, 32},
		{293, 24},
		{295, 16},
		{297, 8},
		{298, 0},
		{362, 32},
	}
	_pcsp__atof_native = [][2]uint32{
		{1, 0},
		{4, 8},
		{587, 56},
		{591, 8},
		{593, 0},
	}
	_pcsp__decimal_to_f64 = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{12, 40},
		{13, 48},
		{1673, 56},
		{1677, 48},
		{1678, 40},
		{1680, 32},
		{1682, 24},
		{1684, 16},
		{1686, 8},
		{1690, 0},
		{1702, 56},
	}
	_pcsp__right_shift = [][2]uint32{
		{1, 0},
		{318, 8},
		{319, 0},
		{387, 8},
		{388, 0},
		{396, 8},
		{398, 0},
	}
	_pcsp__left_shift = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{363, 24},
		{364, 16},
		{366, 8},
		{367, 0},
		{470, 24},
		{471, 16},
		{473, 8},
		{474, 0},
		{486, 24},
	}
	_pcsp__vsigned = [][2]uint32{
		{1, 0},
		{4, 8},
		{112, 16},
		{113, 8},
		{114, 0},
		{125, 16},
		{126, 8},
		{127, 0},
		{260, 16},
		{261, 8},
		{262, 0},
		{266, 16},
		{267, 8},
		{268, 0},
		{306, 16},
		{307, 8},
		{308, 0},
		{316, 16},
		{317, 8},
		{319, 0},
	}
	_pcsp__vstring = [][2]uint32{
		{1, 0},
		{4, 8},
		{6, 16},
		{8, 24},
		{10, 32},
		{11, 40},
		{105, 56},
		{109, 40},
		{110, 32},
		{112, 24},
		{114, 16},
		{116, 8},
		{118, 0},
	}
	_pcsp__vunsigned = [][2]uint32{
		{1, 0},
		{71, 8},
		{72, 0},
		{83, 8},
		{84, 0},
		{107, 8},
		{108, 0},
		{273, 8},
		{274, 0},
		{312, 8},
		{313, 0},
		{320, 8},
		{322, 0},
	}
)

var Funcs = []loader.CFunc{
	{"__native_entry__", 0, 67, 0, nil},
	{"_f32toa", _entry__f32toa, _size__f32toa, _stack__f32toa, _pcsp__f32toa},
	{"_f64toa", _entry__f64toa, _size__f64toa, _stack__f64toa, _pcsp__f64toa},
	{"_format_significand", _entry__format_significand, _size__format_significand, _stack__format_significand, _pcsp__format_significand},
	{"_format_integer", _entry__format_integer, _size__format_integer, _stack__format_integer, _pcsp__format_integer},
	{"_fsm_exec", _entry__fsm_exec, _size__fsm_exec, _stack__fsm_exec, _pcsp__fsm_exec},
	{"_advance_ns", _entry__advance_ns, _size__advance_ns, _stack__advance_ns, _pcsp__advance_ns},
	{"_advance_string", _entry__advance_string, _size__advance_string, _stack__advance_string, _pcsp__advance_string},
	{"_advance_string_default", _entry__advance_string_default, _size__advance_string_default, _stack__advance_string_default, _pcsp__advance_string_default},
	{"_do_skip_number", _entry__do_skip_number, _size__do_skip_number, _stack__do_skip_number, _pcsp__do_skip_number},
	{"_get_by_path", _entry__get_by_path, _size__get_by_path, _stack__get_by_path, _pcsp__get_by_path},
	{"_skip_one_fast", _entry__skip_one_fast, _size__skip_one_fast, _stack__skip_one_fast, _pcsp__skip_one_fast},
	{"_html_escape", _entry__html_escape, _size__html_escape, _stack__html_escape, _pcsp__html_escape},
	{"_i64toa", _entry__i64toa, _size__i64toa, _stack__i64toa, _pcsp__i64toa},
	{"_u64toa", _entry__u64toa, _size__u64toa, _stack__u64toa, _pcsp__u64toa},
	{"_lspace", _entry__lspace, _size__lspace, _stack__lspace, _pcsp__lspace},
	{"_quote", _entry__quote, _size__quote, _stack__quote, _pcsp__quote},
	{"_skip_array", _entry__skip_array, _size__skip_array, _stack__skip_array, _pcsp__skip_array},
	{"_skip_number", _entry__skip_number, _size__skip_number, _stack__skip_number, _pcsp__skip_number},
	{"_skip_object", _entry__skip_object, _size__skip_object, _stack__skip_object, _pcsp__skip_object},
	{"_skip_one", _entry__skip_one, _size__skip_one, _stack__skip_one, _pcsp__skip_one},
	{"_unquote", _entry__unquote, _size__unquote, _stack__unquote, _pcsp__unquote},
	{"_validate_one", _entry__validate_one, _size__validate_one, _stack__validate_one, _pcsp__validate_one},
	{"_validate_utf8", _entry__validate_utf8, _size__validate_utf8, _stack__validate_utf8, _pcsp__validate_utf8},
	{"_validate_utf8_fast", _entry__validate_utf8_fast, _size__validate_utf8_fast, _stack__validate_utf8_fast, _pcsp__validate_utf8_fast},
	{"_value", _entry__value, _size__value, _stack__value, _pcsp__value},
	{"_vnumber", _entry__vnumber, _size__vnumber, _stack__vnumber, _pcsp__vnumber},
	{"_atof_eisel_lemire64", _entry__atof_eisel_lemire64, _size__atof_eisel_lemire64, _stack__atof_eisel_lemire64, _pcsp__atof_eisel_lemire64},
	{"_atof_native", _entry__atof_native, _size__atof_native, _stack__atof_native, _pcsp__atof_native},
	{"_decimal_to_f64", _entry__decimal_to_f64, _size__decimal_to_f64, _stack__decimal_to_f64, _pcsp__decimal_to_f64},
	{"_right_shift", _entry__right_shift, _size__right_shift, _stack__right_shift, _pcsp__right_shift},
	{"_left_shift", _entry__left_shift, _size__left_shift, _stack__left_shift, _pcsp__left_shift},
	{"_vsigned", _entry__vsigned, _size__vsigned, _stack__vsigned, _pcsp__vsigned},
	{"_vstring", _entry__vstring, _size__vstring, _stack__vstring, _pcsp__vstring},
	{"_vunsigned", _entry__vunsigned, _size__vunsigned, _stack__vunsigned, _pcsp__vunsigned},
}
