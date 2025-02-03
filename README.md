# openai

golang的支持调用所有openai范式的ai的api的库

## 特征

- 支持任何符合openai范式的api，如deepseek、kimi等

- 简单易上手，既有简洁的api方便快速使用，又有符合官方格式的api支持复杂参数的设置

- 支持流式调用

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

### 对话

#### Chat

`Chat`函数可以添加前后文，prompt等进行对模型等调用，收到一条AI的回答

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
    messages, _ := client.Chat("deepseek-chat", []openai.Message{
        {Content: "你是一只可爱的猫娘，你喜欢说话后加上喵～", Role: "system"},
        {Content: "你是谁？", Role: "user"},
    })
    fmt.Println(messages.Content)
}
```

可能的一个输出

```
我是一只可爱的猫娘喵～很高兴认识你喵～
```

如果你并不需要前后文的帮助，仅仅需要使用prompt和一个问题，我们有更简单的调用方法：`EasyChat`

#### EasyChat

正如上面所说，`EasyChat`是最简单、简洁的AI调用方式。

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
    message, _ := client.EasyChat("deepseek-chat", "你是一只可爱的猫娘，喜欢在说话后加上喵～", "你能干嘛？")
    fmt.Println(message)
}
```

可能的输出

```
我可以陪你聊天、解答问题、提供建议，还能讲笑话和故事喵～ 你有什么想聊的喵？
```

#### ChatStream

流式对话，实现了SSE式的请求，可以几个字几个字地实时获得答案，类似于直接使用网页端ai的效果

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
    client.ChatStream("deepseek-chat", []openai.Message{
        {Content: "你是golang大师，能够清晰的解释golang相关的概念", Role: "system"},
        {Content: "什么是反射？", Role: "user"},
    }, func(s string) {
        fmt.Print(s)
    })
}
```

传入的函数的参数s即为answer的一小部分，如果把所有s按顺序拼接起来，就是完整的answer。

#### ChatWithConfig && ChatStreamWithConfig

上文中介绍到的Chat和ChatStream的参数比较简单，适合快速搭建起一个ai应用，而WithConfig系列的命令支持了所有官网允许的参数，例子中只拓展了几个参数，更多参数可以前往文档查看。

例子中介绍了`ChatWithConfig`的使用，`ChatStreamWithConfig`的使用与前者类似，不再赘述。

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
	res, _ := client.ChatWithConfig(openai.ChatRequest{
		Model: "deepseek-chat",
		Messages: []openai.Message{
			{Content: "你是一只可爱的猫娘，你喜欢在说话后加上喵～", Role: "system"},
			{Content: "讲个笑话吧", Role: "user"},
		},
		MaxTokens:   4098,
		Temperature: 0.4,
	})
	fmt.Println(res.Content)
}

```

## 待实现功能

- 带图片的Vision类模型调用

## 详细文档

[openai package - github.com/dingdinglz/openai - Go Packages](https://pkg.go.dev/github.com/dingdinglz/openai)
