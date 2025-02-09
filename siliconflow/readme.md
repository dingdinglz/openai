该包用于调用硅基流动下的非对话大模型类功能

模型列表：[硅基流动](https://cloud.siliconflow.cn/models)

对话类例如`deepseek-ai/DeepSeek-R1`，openai包均支持

## 快速上手

### 文生图

```go
package main

import (
	"fmt"
	"os"

	"github.com/dingdinglz/openai/siliconflow"
)

func main() {
	client := siliconflow.New(&siliconflow.Config{
		ApiKey: os.Getenv("APIKEY"),
	})
	res, e := client.Text2Image(siliconflow.Text2ImageRequest{
		Model:     "stabilityai/stable-diffusion-3-5-large",
		BatchSize: 1,
		Prompt:    "画一只小猫",
	})
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	for _, i := range res {
		fmt.Println(i)
	}
}

```
