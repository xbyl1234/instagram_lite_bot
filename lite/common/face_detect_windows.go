package common

import "C"
import (
	"errors"
	"syscall"
	"unsafe"
)

//void* CreateEngine()
//void FreeEngine(ArcFaceEngine* engine)
//uint64 DetectFaceFromPath(ArcFaceEngine* engine, const char* path)
//uint64 DetectFace(ArcFaceEngine* engine, void* data, int dataLen)

func IntPtr(n int) uintptr {
	return uintptr(n)
}

func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringBytePtr(s)))
}

var sexTesting *syscall.LazyDLL
var procCreateEngine *syscall.LazyProc
var procFreeEngine *syscall.LazyProc
var procDetectFace *syscall.LazyProc
var procDetectFaceFromPath *syscall.LazyProc

const (
	FaceDetectError               uint64 = 0b00000001
	FaceDetectImageError          uint64 = 0b00000010
	FaceDetectDetectFaceError     uint64 = 0b00000100
	FaceDetectFaceASFProcessError uint64 = 0b00001000
	FaceDetectBoy                 uint64 = 0b00010000
	FaceDetectGirl                uint64 = 0b00100000
	FaceDetectUnknow              uint64 = 0b01000000
	FaceDetectTrueHuman           uint64 = 0b10000000
)

func InitFaceDetect() {
	sexTesting = syscall.NewLazyDLL("./sex_testing.dll")
	procCreateEngine = sexTesting.NewProc("CreateEngine")
	procFreeEngine = sexTesting.NewProc("FreeEngine")
	procDetectFaceFromPath = sexTesting.NewProc("DetectFaceFromPath")
	procDetectFace = sexTesting.NewProc("DetectFace")
}

type FaceDetect struct {
	handle uintptr
}

type FaceDetectResult struct {
	Sex     string
	IsHuman bool
}

func IsSuccessError(err error) bool {
	return err.Error() == "The operation completed successfully."
}

func makeResult(result uintptr) (*FaceDetectResult, error) {
	var ret FaceDetectResult
	var flag = uint64(result)
	if flag&FaceDetectImageError > 0 {
		return nil, errors.New("ImageError")
	} else if flag&FaceDetectDetectFaceError > 0 {
		return nil, errors.New("DetectFaceError")
	} else if flag&FaceDetectFaceASFProcessError > 0 {
		return nil, errors.New("FaceASFProcessError")
	}
	if flag&FaceDetectBoy > 0 {
		ret.Sex = "boy"
	} else if flag&FaceDetectGirl > 0 {
		ret.Sex = "girl"
	} else if flag&FaceDetectUnknow > 0 {
		ret.Sex = "unknow"
	}
	if flag&FaceDetectTrueHuman > 0 {
		ret.IsHuman = true
	}
	return &ret, nil
}

func (this *FaceDetect) DetectFace(data []byte) (*FaceDetectResult, error) {
	result, _, err := procDetectFace.Call(this.handle, uintptr(unsafe.Pointer(&data[0])), IntPtr(len(data)))
	if err != nil && !IsSuccessError(err) {
		return nil, err
	}
	return makeResult(result)
}

func (this *FaceDetect) DetectFaceFromPath(path string) (*FaceDetectResult, error) {
	result, _, err := procDetectFaceFromPath.Call(this.handle, StrPtr(path))
	if err != nil && !IsSuccessError(err) {
		return nil, err
	}
	return makeResult(result)
}

func CreateEngine() (*FaceDetect, error) {
	ret, _, err := procCreateEngine.Call()
	if err != nil && !IsSuccessError(err) {
		return nil, err
	}
	return &FaceDetect{
		handle: ret,
	}, nil
}

func FreeEngine(engine *FaceDetect) error {
	_, _, err := procFreeEngine.Call(engine.handle)
	if !IsSuccessError(err) {
		return err
	}
	return nil
}
