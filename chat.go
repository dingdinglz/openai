package openai

import (
	"bufio"
	"encoding/json"
	"errors"
)

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatRequest struct {
	Messages         []Message `json:"messages"`
	Model            string    `json:"model"`
	FrequencyPenalty int       `json:"frequency_penalty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	PresencePenalty  int       `json:"presence_penalty"`
	ResponseFormat   struct {
		Type string `json:"type"`
	} `json:"response_format"`
	Stop        []string `json:"stop"`
	Stream      bool     `json:"stream"`
	Temperature float32  `json:"temperature"`
	TopP        int      `json:"top_p"`
}

type realChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
		Message      Message `json:"message"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type realChatStreamResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

type realChatReasonStreamResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index int `json:"index"`
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
}

func checkChatRequest(cr *ChatRequest) {
	if cr.ResponseFormat.Type == "" {
		cr.ResponseFormat.Type = "text"
	}
	if cr.Temperature == 0 {
		cr.Temperature = 0.3
	}
	if cr.TopP == 0 {
		cr.TopP = 1
	}
}

// api /chat/completions 的实现
func (client Client) Chat(model string, messages []Message) (*Message, error) {
	reqBody := ChatRequest{}
	reqBody.Messages = messages
	reqBody.Stream = false
	reqBody.Model = model
	checkChatRequest(&reqBody)
	reqClient := client.newHttpClient()
	jsonBody, e := json.Marshal(reqBody)
	if e != nil {
		return nil, e
	}
	reqClient.SetBody(string(jsonBody))
	res := &realChatResponse{}
	reqClient.SetResult(res)
	httpres, e := reqClient.Post(client.Config.BaseUrl + "/chat/completions")
	if e != nil {
		return nil, e
	}
	if httpres.StatusCode() != 200 {
		errorMessage, e := parseRealError(httpres.Body())
		if e != nil {
			return nil, errors.New(string(httpres.Body()))
		}
		return nil, errors.New(errorMessage)
	}
	return &res.Choices[0].Message, nil
}

// api /chat/completions 的傻瓜式实现
// 没有上下文，给入提示词和问题即可获得答案
func (client Client) EasyChat(model string, prompt string, message string) (string, error) {
	reqBody := ChatRequest{}
	reqBody.Messages = []Message{
		{Content: prompt, Role: "system"},
		{Content: message, Role: "user"},
	}
	reqBody.Stream = false
	reqBody.Model = model
	checkChatRequest(&reqBody)
	reqClient := client.newHttpClient()
	jsonBody, e := json.Marshal(reqBody)
	if e != nil {
		return "", e
	}
	reqClient.SetBody(string(jsonBody))
	res := &realChatResponse{}
	reqClient.SetResult(res)
	httpres, e := reqClient.Post(client.Config.BaseUrl + "/chat/completions")
	if e != nil {
		return "", e
	}
	if httpres.StatusCode() != 200 {
		errorMessage, e := parseRealError(httpres.Body())
		if e != nil {
			return "", errors.New(string(httpres.Body()))
		}
		return "", errors.New(errorMessage)
	}
	return res.Choices[0].Message.Content, nil
}

// api /chat/completions 的流式实现
func (client Client) ChatStream(model string, messages []Message, during func(string)) error {
	reqBody := ChatRequest{}
	reqBody.Messages = messages
	reqBody.Stream = true
	reqBody.Model = model
	checkChatRequest(&reqBody)
	reqClient := client.newStreamClient()
	jsonBody, e := json.Marshal(reqBody)
	if e != nil {
		return e
	}
	reqClient.SetBody(string(jsonBody))
	reqClient.SetDoNotParseResponse(true)
	httpres, e := reqClient.Post(client.Config.BaseUrl + "/chat/completions")
	if e != nil {
		return e
	}
	defer httpres.RawBody().Close()
	scanner := bufio.NewScanner(httpres.RawBody())
	initFlag := true
	for scanner.Scan() {
		_res := scanner.Text()
		if _res == "" {
			continue
		}
		if _res == "data: [DONE]" {
			break
		}
		if initFlag {
			resError, e := parseRealError([]byte(_res))
			if e == nil {
				return errors.New(resError)
			}
			initFlag = false
			continue
		}
		_res = _res[6:]
		var _json realChatStreamResponse
		e := json.Unmarshal([]byte(_res), &_json)
		if e != nil {
			return e
		}
		during(_json.Choices[0].Delta.Content)
	}
	return nil
}

// api /chat/completions 的高自定义度实现
func (client Client) ChatWithConfig(config ChatRequest) (*Message, error) {
	config.Stream = false
	checkChatRequest(&config)
	reqClient := client.newHttpClient()
	jsonBody, e := json.Marshal(config)
	if e != nil {
		return nil, e
	}
	reqClient.SetBody(string(jsonBody))
	res := &realChatResponse{}
	reqClient.SetResult(res)
	httpres, e := reqClient.Post(client.Config.BaseUrl + "/chat/completions")
	if e != nil {
		return nil, e
	}
	if httpres.StatusCode() != 200 {
		errorMessage, e := parseRealError(httpres.Body())
		if e != nil {
			return nil, errors.New(string(httpres.Body()))
		}
		return nil, errors.New(errorMessage)
	}
	return &res.Choices[0].Message, nil
}

// api /chat/completions 的高自定义度流式实现
func (client Client) ChatStreamWithConfig(config ChatRequest, during func(string)) error {
	config.Stream = true
	checkChatRequest(&config)
	reqClient := client.newStreamClient()
	jsonBody, e := json.Marshal(config)
	if e != nil {
		return e
	}
	reqClient.SetBody(string(jsonBody))
	reqClient.SetDoNotParseResponse(true)
	httpres, e := reqClient.Post(client.Config.BaseUrl + "/chat/completions")
	if e != nil {
		return e
	}
	defer httpres.RawBody().Close()
	scanner := bufio.NewScanner(httpres.RawBody())
	initFlag := true
	for scanner.Scan() {
		_res := scanner.Text()
		if _res == "" {
			continue
		}
		if _res == "data: [DONE]" {
			break
		}
		if initFlag {
			resError, e := parseRealError([]byte(_res))
			if e == nil {
				return errors.New(resError)
			}
			initFlag = false
			continue
		}
		_res = _res[6:]
		var _json realChatStreamResponse
		e := json.Unmarshal([]byte(_res), &_json)
		if e != nil {
			return e
		}
		during(_json.Choices[0].Delta.Content)
	}
	return nil
}

// 相比于ChatStream，ChatReasonStream支持了深度思考的模型
func (client Client) ChatReasonStream(model string, messages []Message, think func(string), during func(string)) error {
	reqBody := ChatRequest{}
	reqBody.Messages = messages
	reqBody.Stream = true
	reqBody.Model = model
	checkChatRequest(&reqBody)
	reqClient := client.newStreamClient()
	jsonBody, e := json.Marshal(reqBody)
	if e != nil {
		return e
	}
	reqClient.SetBody(string(jsonBody))
	reqClient.SetDoNotParseResponse(true)
	httpres, e := reqClient.Post(client.Config.BaseUrl + "/chat/completions")
	if e != nil {
		return e
	}
	defer httpres.RawBody().Close()
	scanner := bufio.NewScanner(httpres.RawBody())
	initFlag := true
	for scanner.Scan() {
		_res := scanner.Text()
		if _res == "" {
			continue
		}
		if _res == "data: [DONE]" {
			break
		}
		if initFlag {
			resError, e := parseRealError([]byte(_res))
			if e == nil {
				return errors.New(resError)
			}
			initFlag = false
			continue
		}
		_res = _res[6:]
		var _json realChatReasonStreamResponse
		e := json.Unmarshal([]byte(_res), &_json)
		if e != nil {
			return e
		}
		if _json.Choices[0].Delta.ReasoningContent != "" {
			think(_json.Choices[0].Delta.ReasoningContent)
		}
		if _json.Choices[0].Delta.Content != "" {
			during(_json.Choices[0].Delta.Content)
		}
	}
	return nil
}
