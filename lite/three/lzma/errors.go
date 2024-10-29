package lzma

import "errors"

var (
	ErrCorrupted           = errors.New("corrupted")
	ErrIncorrectProperties = errors.New("incorrect LZMA properties")
	ErrResultError         = errors.New("result error")
	ErrDictOutOfRange      = errors.New("dictionary capacity is out of range")
	ErrUnexpectedLZMA2Code = errors.New("unexpected lzma2 code")
	ErrNoLZMAReader        = errors.New("no lzma reader on chunkLZMAResetState")
)
