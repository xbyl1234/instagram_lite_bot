package native

/* #include <stdint.h>

static size_t convertToUInt64(unsigned long long integer) {
	return (size_t) integer;
}

static uint32_t convertToUInt32(unsigned long integer) {
	return (uint32_t) integer;
}
*/
import "C"

func toInt64(in int64) C.size_t {
	return C.convertToUInt64(C.ulonglong(in))
}

func intToInt64(in int) C.size_t {
	return toInt64(int64(in))
}

func toUInt32(in uint32) C.uint32_t {
	return C.convertToUInt32(C.ulong(in))
}
