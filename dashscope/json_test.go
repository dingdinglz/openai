package dashscope

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	testjson := `{
		"code":"InvalidApiKey",
		"message":"Invalid API-key provided.",
		"request_id":"fb53c4ec-1c12-4fc4-a580-xxxxxx"
	}`
	var rTest text2imageSuccess
	e := json.Unmarshal([]byte(testjson), &rTest)
	if e != nil {
		fmt.Println(e.Error())
	}
}
