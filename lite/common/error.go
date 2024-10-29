package common

type MakeMoneyError struct {
	ErrStr    string
	ErrType   ErrType
	ExternErr error
	HttpCode  int
}

// 实现 `error` 接口
func (this *MakeMoneyError) Error() string {
	if this.ErrStr != "" {
		return this.ErrStr
	}
	if this.ExternErr != nil {
		return this.ExternErr.Error()
	}

	var errHead string
	switch this.ErrType {
	case ApiError:
		errHead = "api error"
		break
	case PorxyError:
		errHead = "proxy error"
		break
	case NoMoreError:
		errHead = "no more error"
		break
	case OtherError:
		errHead = "other error"
		break
	case RequestError:
		errHead = "request error"
		break
	}
	return errHead
}

type ErrType int

var (
	ApiError               ErrType = 0
	PorxyError             ErrType = 1
	NoMoreError            ErrType = 2
	OtherError             ErrType = 3
	RequestError           ErrType = 4
	ChallengeRequiredError ErrType = 5
	LoginRequiredError     ErrType = 6
	RecvPhoneCodeError     ErrType = 7
	FeedbackError          ErrType = 9
	ParamsError            ErrType = 10
)

func IsHttpErrorCode(err error, code int) bool {
	e, ok := err.(*MakeMoneyError)
	if ok {
		return e.HttpCode == code
	}
	return false
}

func IsNoMoreError(err error) bool {
	e, ok := err.(*MakeMoneyError)
	if ok {
		return e.ErrType == NoMoreError
	}
	return false
}

func IsError(err error, errType ErrType) bool {
	e, ok := err.(*MakeMoneyError)
	if ok {
		return e.ErrType == errType
	}
	return false
}

func GetErrorMsg(err error) string {
	e, ok := err.(*MakeMoneyError)
	if ok {
		return e.ErrStr
	}
	return ""
}

func NerError(msg string) *MakeMoneyError {
	return &MakeMoneyError{
		ErrStr: msg,
	}
}
