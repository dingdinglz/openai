package dashscope

import "github.com/go-resty/resty/v2"

type Client struct {
	config *Config
}

func (client *Client) newHttpClient() *resty.Request {
	return resty.New().R().SetHeader("Content-Type", "application/json").SetHeader("Authorization", "Bearer "+client.config.ApiKey)
}

func New(config *Config) *Client {
	return &Client{
		config: config,
	}
}
