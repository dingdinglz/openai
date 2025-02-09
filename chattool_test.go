package openai

import (
	"fmt"
	"os"
	"testing"
)

func TestChatToolStream(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseUrl: "https://api.siliconflow.cn/v1",
		ApiKey:  os.Getenv("APIKEY"),
	})
	client.ChatWithTools("deepseek-ai/DeepSeek-V3", []ToolMessage{
		{Role: "system", Content: "你是智能机器人，可以辅助用户操作，如果调用的函数出现了错误，直接把错误信息原样返回"},
		{Role: "user", Content: "帮我删除test.txt"},
	}, []ChatToolFunction{
		{Type: "function", Function: ChatToolFuctionDetail{
			Name:        "deleteFile",
			Description: "给入文件名，可以删除对应的文件，返回删除的结果",
			Parameters: ChatToolParameters{
				Type: "object",
				Properties: map[string]ChatToolFuctionPropertie{
					"file": {
						Type:        "string",
						Description: "文件名",
					},
				},
				Required: []string{"location"},
			},
		}},
	}, map[string]func(map[string]interface{}) string{
		"deleteFile": func(i map[string]interface{}) string {
			fmt.Println("假装真的删除了文件" + i["file"].(string))
			return "因为服务器正在被dinglz拿来玩原神，所以无法删除"
		},
	}, func(s string) {
		fmt.Println("回复：" + s)
	})
}
