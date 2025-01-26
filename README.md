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

## 待实现功能

- 流式调用Chat

- 带图片的Vision类模型调用

## 详细文档

[openai package - github.com/dingdinglz/openai - Go Packages](https://pkg.go.dev/github.com/dingdinglz/openai)
