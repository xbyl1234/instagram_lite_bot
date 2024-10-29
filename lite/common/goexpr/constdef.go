package goexpr

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	Variable string = "Variable"
	Func     string = "function"
	PreOp    string = "prefixOp"
	InfxOp   string = "infixOp"
	Value    string = "value"
	Array    string = "array"
	Args     string = "funcArgs"
)

const (
	Equal      string = "=="
	IN         string = "IN"
	NotIN      string = "NotIN"
	Less       string = "<"
	LessEqual  string = "<="
	AboveEqual string = ">="
	Above      string = ">"
	NotEqual   string = "!="
	Add        string = "+"
	Sub        string = "-"
	Mult       string = "*"
	Div        string = "/"
	Rest       string = "%"
	And        string = "&&"
	Or         string = "||"
	MathAnd    string = "&"
	MathOr     string = "|"
	Not        string = "!"
	BraktLeft  string = "("
	BraktRight string = ")"
	ArrayLeft  string = "["
	ArrayRight string = "]"
	ItemSpit   string = ","
)

var OpPriority = map[string]int32{
	Mult:       60,
	Div:        60,
	Rest:       60,
	MathAnd:    60,
	MathOr:     60,
	Add:        50,
	Sub:        50,
	Above:      40,
	AboveEqual: 40,
	Less:       40,
	LessEqual:  40,
	Equal:      40,
	NotEqual:   40,
	IN:         30,
	NotIN:      30,
	And:        20,
	Or:         10,
}

type InputArgs any

// 函数运算
type FunctionOp func(inputArgs InputArgs, v ...interface{}) interface{}

// 前缀运算
type PrefixOp func(v interface{}) interface{}

// 中缀运算
type InfixOp func(v1, v2 interface{}) interface{}

var PrefixOpSet = map[string]PrefixOp{
	Not: func(v1 interface{}) interface{} {
		return !v1.(bool)
	},
	Sub: func(v interface{}) interface{} {
		return -FloatVal(v)
	},
}

var InfixOpSet = map[string]InfixOp{
	MathAnd: func(v1, v2 interface{}) interface{} {
		return Int64Val(v1) & Int64Val(v2)
	},
	MathOr: func(v1, v2 interface{}) interface{} {
		return Int64Val(v1) | Int64Val(v2)
	},
	Equal: func(v1, v2 interface{}) interface{} {
		return fmt.Sprint(v1) == fmt.Sprint(v2)
	},
	Add: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) + FloatVal(v2)
	},
	Sub: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) - FloatVal(v2)
	},
	Mult: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) * FloatVal(v2)
	},
	Div: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) / FloatVal(v2)
	},
	Rest: func(v1, v2 interface{}) interface{} {
		return int64(FloatVal(v1)) % int64(FloatVal(v2))
	},
	And: func(v1, v2 interface{}) interface{} {
		return v1.(bool) && v2.(bool)
	},
	Or: func(v1, v2 interface{}) interface{} {
		return v1.(bool) || v2.(bool)
	},
	Less: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) < FloatVal(v2)
	},
	LessEqual: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) <= FloatVal(v2)
	},
	AboveEqual: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) >= FloatVal(v2)
	},
	Above: func(v1, v2 interface{}) interface{} {
		return FloatVal(v1) > FloatVal(v2)
	},
	NotEqual: func(v1, v2 interface{}) interface{} {
		return fmt.Sprint(v1) != fmt.Sprint(v2)
	},
	IN: func(v1 interface{}, v2 interface{}) interface{} {
		return in(v1, v2)
	},
	NotIN: func(v1 interface{}, v2 interface{}) interface{} {
		return notIn(v1, v2)
	},
}

func notIn(a, b interface{}) interface{} {
	if b == nil {
		return true
	}
	aStr := fmt.Sprint(a)
	array := reflect.ValueOf(b)
	length := array.Len()
	for i := 0; i < length; i++ {
		bStr := fmt.Sprint(array.Index(i).Interface())
		if bStr == aStr {
			return false
		}
	}
	return true
}

func in(a, b interface{}) interface{} {
	if b == nil {
		return false
	}
	aStr := fmt.Sprint(a)
	array := reflect.ValueOf(b)
	length := array.Len()
	for i := 0; i < length; i++ {
		bStr := fmt.Sprint(array.Index(i).Interface())
		if bStr == aStr {
			return true
		}
	}
	return false
}

func FloatVal(v interface{}) float64 {
	s := fmt.Sprint(v)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func Int64Val(v interface{}) int64 {
	s := fmt.Sprint(v)
	f, _ := strconv.ParseInt(s, 10, 64)
	return f
}
