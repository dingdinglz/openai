package openai

import (
	"fmt"
	"os"
	"testing"
)

func TestChat(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseUrl: "https://api.deepseek.com/v1",
		ApiKey:  os.Getenv("OPENAI_KEY"),
	})

	// 测试Chat
	res, e := client.Chat("deepseek-chat", []Message{
		{Content: "你是一只可爱的猫娘，你喜欢在说话后加上喵～", Role: "system"},
		{Content: "你是谁？", Role: "user"},
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res.Content)

	// 测试EasyChat
	resMessage, e := client.EasyChat("deepseek-chat", "你是一只可爱的猫娘，你喜欢在说话后加上喵～", "你会干嘛？")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(resMessage)

	// 测试ChatWithConfig
	res, e = client.ChatWithConfig(ChatRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Content: "你是一只可爱的猫娘，你喜欢在说话后加上喵～", Role: "system"},
			{Content: "讲个笑话吧", Role: "user"},
		},
		MaxTokens: 4098,
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}

func TestChatStream(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseUrl: "https://api.deepseek.com/v1",
		ApiKey:  os.Getenv("OPENAI_KEY"),
	})
	e := client.ChatStream("deepseek-chat", []Message{
		{Content: "你是一个golang领域的专家，擅长解释概念", Role: "system"},
		{Content: "什么是反射？", Role: "user"},
	}, func(s string) {
		fmt.Print(s)
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	e = client.ChatStreamWithConfig(ChatRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Content: "你是一只可爱的猫娘，你喜欢在说话后加上喵～", Role: "system"},
			{Content: "讲个笑话吧", Role: "user"},
		},
		MaxTokens: 4098,
	}, func(s string) {
		fmt.Print(s)
	})
	if e != nil {
		t.Error(e.Error())
	}
}

func TestChatReasonStream(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseUrl: "https://api.siliconflow.cn/v1",
		ApiKey:  os.Getenv("OPENAI_KEY"),
	})
	reasoning := true
	e := client.ChatReasonStream("deepseek-ai/DeepSeek-R1", []Message{
		{Content: "你是一名golang大师", Role: "system"},
		{Content: "解释一下context的用处", Role: "user"},
	}, func(s string) {
		// 思考内容
		fmt.Print(s)
	}, func(s string) {
		if reasoning {
			fmt.Println("\n思考结束！\nThink End")
			reasoning = false
		}
		fmt.Print(s)
	})
	if e != nil {
		t.Error(e.Error())
	}
}
