package openai

import (
	"errors"
)

type realModels struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

// api /models 的实现
func (client *Client) Models() ([]Model, error) {
	req := client.newHttpClient()
	result := &realModels{}
	req.SetResult(result)
	res, e := req.Get(client.Config.BaseUrl + "/models")
	if e != nil {
		return []Model{}, e
	}
	if res.StatusCode() != 200 {
		resError, e := parseRealError(res.Body())
		if e != nil {
			return []Model{}, errors.New(string(res.Body()))
		}
		return []Model{}, errors.New(resError)
	}
	return result.Data, nil
}
