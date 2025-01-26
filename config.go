package openai

type ClientConfig struct {
	// 模型api的地址，例如：https://api.deepseek.com
	//
	// 请不要在结尾处加上/ 例如：https://api.deepseek.com合法，而https://api.deepseek.com/不合法
	BaseUrl string
	ApiKey  string // api key
}
