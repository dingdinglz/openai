该包用于调用通义千问下的非对话大模型类功能

模型列表：[阿里云百炼](https://bailian.console.aliyun.com/#/home)

对话类例如`qwen-plus`，openai包均支持

## 快速上手

### 文生图

```go
package main

import (
	"fmt"
	"os"

	"github.com/dingdinglz/openai/dashscope"
)

func main() {
	client := dashscope.New(&dashscope.Config{
		ApiKey: os.Getenv("DASHSCOPE_KEY"),
	})
	task_id, _ := client.Text2image(dashscope.Text2ImageRequest{
		Model: "wanx2.1-t2i-turbo",
		Input: dashscope.Text2ImageRequestInput{
			Prompt: "画一只可爱的小猫",
		},
		Parameters: dashscope.Text2ImageRequestParameter{
			N: 1,
		},
	})
	for {
		res, _ := client.Text2imageResult(task_id)
		if res.Output.TaskStatus == "FAILED" {
			fmt.Println("生成失败")
			break
		}
		if res.Output.TaskStatus == "SUCCEEDED" {
			fmt.Println(res.Output.Results[0].URL)
		}
	}
}

```
