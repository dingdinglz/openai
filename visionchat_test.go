package openai

import (
	"fmt"
	"os"
	"testing"
)

func TestVisionChatStream(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseUrl: "https://api.siliconflow.cn/v1",
		ApiKey:  os.Getenv("APIKEY"),
	})
	imageData, _ := os.ReadFile("test.png")
	e := client.ChatVisionStream("deepseek-ai/deepseek-vl2", []VisionMessage{
		{
			Role: "user",
			Content: []VisionContent{
				{
					Type: VISION_MESSAGE_IMAGE_URL,
					ImageUrl: &VisionContentImageUrl{
						Url: GenerateImageUrlBase64(imageData),
					},
				},
				{
					Type: VISION_MESSAGE_TEXT,
					Text: "这只猫在干嘛",
				},
			},
		},
	}, func(s string) {
		fmt.Print(s)
	})
	if e != nil {
		t.Error(e.Error())
	}
}
