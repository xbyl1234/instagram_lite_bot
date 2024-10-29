package goexpr

//
//func Contain(a, b interface{}) interface{} {
//	bStr := fmt.Sprint(b)
//	array := reflect.ValueOf(a)
//	length := array.Len()
//	for i := 0; i < length; i++ {
//		aStr := fmt.Sprint(array.Index(i).Interface())
//		if bStr == aStr {
//			return true
//		}
//	}
//	return false
//}
//
//func TestEngine(t *testing.T) {
//	exprs := `(2)-4+3>(-9)&&5<4+5&&3NotIN[1,2,4]&&ADD(1,2)<4&&-(#-3-4)<=30&&4>1&&[1,2,4] Contain 4 && ADD(1,2)!=1 && user.name=='kiteee' && user_count>20`
//	//exprs = `-------1`
//	//exprs=`user.name=='kiteee' && user_count>20`
//	//exprs=`user_count>20 && user_count>20`
//	//exprs=`#--3*-4-#2`
//	//exprs=`-4-#2`
//	//exprs = `3NotIN([1,2,3])&&ADD(1,2)<4`
//	//exprs = `[1,2,4] Contain 4 IN [true] NotIN [false]`
//	eg := NewEngine()
//	eg.AddFunc("ADD", func(v ...interface{}) interface{} {
//		return FloatVal(v[0]) + FloatVal(v[1])
//	})
//	eg.AddPrefix("#", func(v interface{}) interface{} {
//		return FloatVal(v) * FloatVal(v)
//	})
//	eg.AddInfix("Contain", 30, func(v1, v2 interface{}) interface{} {
//		return Contain(v1, v2)
//	})
//	var params = map[string]interface{}{
//		"user": map[string]interface{}{
//			"name": "kiteee",
//			"age":  50,
//		},
//		"user_count": 30,
//	}
//	//eg.SetPriority("NotIN", 30)
//	result := eg.Execute(exprs, params)
//	fmt.Println(result)
//}
//
//func TestSpitExpr(t *testing.T) {
//	exprs := `('ssss','aaa',['aaa',bb],[aaa,game notIn values],funa bsna ( sa,ssdf), atype contains (90),type notIn [1,2,4],value >= images ,-add(),[name,add()],-otEs(),[[-aaam,bb],bbb])`
//	exprs = `()`
//	result := SpitExpr(exprs)
//	for _, v := range result {
//		fmt.Println(v)
//	}
//}
//
//func TestGetArgs(t *testing.T) {
//	var mp = map[string]interface{}{
//		"user": map[string]interface{}{
//			"name": "kiteee",
//			"age":  50,
//		},
//	}
//	va := GetArg("user.name", mp)
//	fmt.Println(va)
//}
