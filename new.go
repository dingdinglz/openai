package openai

func NewClient(config *ClientConfig) *Client {
	return &Client{
		Config: config,
	}
}
