// github.com/dingdinglz/openai 致力于解决对支持openai协议的ai的调用
// 如国内的deepseek，kimi等，以及gemini通过特定的程式可以转为openai协议
package openai

import "github.com/go-resty/resty/v2"

// Client 是请求客户端，可以通过client调用不同的api
type Client struct {
	Config *ClientConfig
}

// 内部生成一个用于访问的内容
func (client Client) newHttpClient() *resty.Request {
	return resty.New().R().SetHeader("Content-Type", "application/json").SetHeader("Authorization", "Bearer "+client.Config.ApiKey)
}

// 用于SSE的请求
func (client Client) newStreamClient() *resty.Request {
	return resty.New().R().SetHeader("Content-Type", "application/json").SetHeader("Authorization", "Bearer "+client.Config.ApiKey)
}
