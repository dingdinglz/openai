package openai

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	type Test struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Num  int
		O    bool
	}
	res, _ := json.Marshal(&Test{
		Name: "11",
	})
	fmt.Println(string(res))
}
