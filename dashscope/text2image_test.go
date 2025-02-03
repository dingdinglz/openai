package dashscope

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestText2Image(t *testing.T) {
	client := New(&Config{
		ApiKey: os.Getenv("OPENAI_KEY"),
	})
	res, e := client.Text2image(Text2ImageRequest{
		Model: "wanx2.0-t2i-turbo",
		Input: Text2ImageRequestInput{
			Prompt: "你可以把《原神》想象成一个巨大的主题公园，里面有各种不同的区域（如雪山、沙漠、森林），每个区域都有独特的风景、生物和挑战。你扮演一个旅行者，可以自由地在这个公园里探索，解谜，战斗，并与各种角色互动。",
		},
		Parameters: Text2ImageRequestParameter{
			Size:         "1024*1024",
			N:            1,
			PromptExtend: false,
		},
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
	for {
		time.Sleep(time.Second)
		res2, e := client.Text2imageResult(res)
		if e != nil {
			t.Error(e.Error())
			return
		}
		if res2.Output.TaskStatus == "UNKNOWN" || res2.Output.TaskStatus == "FAILED" {
			t.Error(res2.Output.Message)
			return
		}
		if res2.Output.TaskStatus != "SUCCEEDED" {
			fmt.Println("任务执行中...")
		} else {
			fmt.Println(res2.Output.Results[0].URL)
			break
		}
	}
}
