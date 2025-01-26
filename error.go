package openai

import "encoding/json"

type realError struct {
	Error struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   interface{} `json:"param"`
		Code    string      `json:"code"`
	} `json:"error"`
}

// 解析错误
func parseRealError(data []byte) (string, error) {
	var r realError
	e := json.Unmarshal(data, &r)
	if e != nil {
		return "", e
	}
	return r.Error.Message, nil
}
