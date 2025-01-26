# openai

golang的支持调用所有openai范式的ai的api的库

> [!WARNING]
> 本库处于积极更新状态

## 特征

- 支持任何符合openai范式的api，如deepseek、kimi等

- 简单易上手，既有简洁的api方便快速使用，又有符合官方格式的api支持复杂参数的设置

- 支持功能齐全

## 安装

```shell
go get github.com/dingdinglz/openai
```

## 快速上手

### 查看模型列表

```go
package main

import (
	"fmt"
	"os"

	"github.com/dingdinglz/openai"
)

func main() {
	client := openai.NewClient(&openai.ClientConfig{
		BaseUrl: "https://api.deepseek.com/v1",
		ApiKey:  os.Getenv("DEEPSEEK_APIKEY"),
	})
	models, _ := client.Models()
	fmt.Println(models)
}
```

`model`是Model类型的数组

输出内容可能为：

```
[{deepseek-chat model deepseek} {deepseek-reasoner model deepseek}]
```
