# openai

golang的支持调用所有openai范式的ai的api的库

## 特征

- 支持任何符合openai范式的api，如deepseek、kimi等

- 简单易上手，既有简洁的api方便快速使用，又有符合官方格式的api支持复杂参数的设置

- 支持流式调用

- 支持视觉模型和深度思考模型

- 支持功能齐全

## 说明

仓库主体均为llm大模型调用内容且只支持openai的api范式，仓库主体类似于`python`的`openai`包

随着业务的增广，增加了对于其他ai厂商api格式的支持（包括且不限于对话，可能包括文生图等等），具体实现内容请前往对应的板块查看

- [DashScope（阿里云百炼）](https://github.com/dingdinglz/openai/tree/main/dashscope)

- [SiliconFlow（硅基流动）](https://github.com/dingdinglz/openai/tree/main/siliconflow)

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

#### ChatReasonStream

由于deepseek-r1的爆火，部分模型同样支持深度思考，`ChatStream`并不能获取到深度思考的内容，因此添加`ChatReasonStream`，可以流式获取到思考的内容，下面的例子是对deepseek-r1的调用。

> 需要注意的是，由于服务提供商遍地开花，部分厂商把深度思考的内容加在了message里，那么ChatStream是可以直接获取到的

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
    client.ChatReasonStream("deepseek-reasoner", []openai.Message{
        {Content: "你是一只可爱的猫娘喵~", Role: "system"},
        {Content: "给我讲个故事吧", Role: "user"},
    }, func(s string) {
        // think 部分 ， 深度思考的内容通过s给出
        fmt.Print(s)
    }, func(s string) {
        // message 部分 ， 真正的回答通过s给出
        fmt.Print(s)
    })
}
```

#### ChatVisionStream

调用视觉模型对话，下面的例子是上传了一张本地图片test.png，并要求ai解题

Url参数同样可以参数真实的url，也可以像例子中一样用base64编码上传本地图片

```go
package main

import (
	"fmt"
	"os"

	"github.com/dingdinglz/openai"
)

func main() {
	client := openai.NewClient(&openai.ClientConfig{
		BaseUrl: "https://api.siliconflow.cn/v1",
		ApiKey:  os.Getenv("APIKEY"),
	})
	imageData, _ := os.ReadFile("test.png")
	client.ChatVisionStream("deepseek-ai/deepseek-vl2", []openai.VisionMessage{
		{
			Role: "user",
			Content: []openai.VisionContent{
				{
					Type: openai.VISION_MESSAGE_IMAGE_URL,
					ImageUrl: &openai.VisionContentImageUrl{
						Url: openai.GenerateImageUrlBase64(imageData),
					},
				},
				{
					Type: openai.VISION_MESSAGE_TEXT,
					Text: "请帮我解一下这道题",
				},
			},
		},
	}, func(s string) {
		fmt.Print(s)
	})
}
```

#### Fuction Calling

[什么是function calling](https://docs.siliconflow.cn/cn/userguide/guides/function-calling)

[视频 - 如何用本包使用function calling && 什么是function calling](https://www.bilibili.com/video/BV14wNSeWEGR/?share_source=copy_web&vd_source=48d9e62f9891701ebeb6dd853a402b14)

上面的视频对该功能有详细的介绍，和示例代码的编写，下面是示例代码

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/openai"
)

func main() {
    client := openai.NewClient(&openai.ClientConfig{
        BaseUrl: "https://api.siliconflow.cn/v1",
        ApiKey:  os.Getenv("APIKEY"),
    })
    client.ChatWithTools("Qwen/Qwen2.5-14B-Instruct", []openai.ToolMessage{
        {Role: "user", Content: "南京现在天气怎么样"},
    }, []openai.ChatToolFunction{
        {
            Type: "function",
            Function: openai.ChatToolFuctionDetail{
                Name:        "weather",
                Description: "通过传入地区名，获取该地区的天气状况",
                Parameters: openai.ChatToolParameters{
                    Type: "object",
                    Properties: map[string]openai.ChatToolFuctionPropertie{
                        "location": {
                            Type:        "string",
                            Description: "地区名",
                        },
                    },
                },
            },
        },
    }, map[string]func(map[string]interface{}) string{
        "weather": func(m map[string]interface{}) string {
            fmt.Println(m["location"].(string) + "的数据被请求了")
            return "不好"
        },
    }, func(s string) {
        fmt.Println("回复内容:" + s)
    })
}
```

## 详细文档

[openai package - github.com/dingdinglz/openai - Go Packages](https://pkg.go.dev/github.com/dingdinglz/openai)
