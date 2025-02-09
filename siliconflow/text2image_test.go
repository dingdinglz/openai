package siliconflow

import (
	"fmt"
	"os"
	"testing"
)

func TestText2Image(t *testing.T) {
	client := New(&Config{
		ApiKey: os.Getenv("APIKEY"),
	})
	res, e := client.Text2Image(Text2ImageRequest{
		Model:             "deepseek-ai/Janus-Pro-7B",
		Prompt:            "画一个美女，脖子上戴着项链",
		PromptEnhancement: true,
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	for _, i := range res {
		fmt.Println(i)
	}
}
