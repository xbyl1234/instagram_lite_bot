package facebook

import (
	"CentralizedControl/common/utils"
)

func SerializeForm(userSet map[string]string, seq []string, temp map[string]string, mode utils.Encoding) string {
	data := ""
	for _, key := range seq {
		data += key + "="
		if value, ok := userSet[key]; ok {
			data += utils.Escape(value, mode)
		} else {
			data += utils.Escape(temp[key], mode)
		}
		data += "&"
	}
	return data[:len(data)-1]
}

type Body interface {
	Serialize() string
}

type BaseBody struct {
	req *ApiRequest
}
