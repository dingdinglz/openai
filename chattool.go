package openai

import (
	"bufio"
	"encoding/json"
)

type ChatToolRequest struct {
	Messages         []ToolMessage `json:"messages"`
	Model            string        `json:"model"`
	FrequencyPenalty int           `json:"frequency_penalty"`
	MaxTokens        int           `json:"max_tokens,omitempty"`
	PresencePenalty  int           `json:"presence_penalty"`
	ResponseFormat   struct {
		Type string `json:"type"`
	} `json:"response_format"`
	Stop        []string           `json:"stop"`
	Stream      bool               `json:"stream"`
	Temperature float32            `json:"temperature"`
	TopP        int                `json:"top_p"`
	Tools       []ChatToolFunction `json:"tools"`
}

type ChatToolFunction struct {
	Type     string                `json:"type"`
	Function ChatToolFuctionDetail `json:"function"`
}

type ChatToolFuctionDetail struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  ChatToolParameters `json:"parameters"`
}

type ChatToolParameters struct {
	Type       string                              `json:"type"`
	Properties map[string]ChatToolFuctionPropertie `json:"properties"`
	Required   []string                            `json:"required"`
}

type ChatToolFuctionPropertie struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type realChatToolStreamMessage struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role      string          `json:"role"`
			Content   string          `json:"content"`
			ToolCalls []realToolCalls `json:"tool_calls"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
		PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
		PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

type ToolMessage struct {
	Role       string          `json:"role"`
	Content    string          `json:"content"`
	ToolCalls  []realToolCalls `json:"tool_calls"`
	ToolCallID string          `json:"tool_call_id"`
}

type realToolCalls struct {
	Index    int    `json:"index"`
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

func checkChatToolRequest(cr *ChatToolRequest) {
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

func (client *Client) ChatWithTools(model string, messages []ToolMessage, toolInfo []ChatToolFunction, functionMap map[string](func(map[string]interface{}) string), during func(string)) error {
	reqBody := ChatToolRequest{}
	reqBody.Tools = toolInfo
	reqBody.Messages = messages
	reqBody.Model = model
	checkChatToolRequest(&reqBody)
	body, e := json.Marshal(reqBody)
	if e != nil {
		return e
	}
	reqClient := client.newStreamClient()
	reqClient.SetBody(body)
	reqClient.SetDoNotParseResponse(true)
	httpres, e := reqClient.Post(client.Config.BaseUrl + "/chat/completions")
	if e != nil {
		return e
	}
	defer httpres.RawBody().Close()
	scanner := bufio.NewScanner(httpres.RawBody())
	for scanner.Scan() {
		_res := scanner.Text()
		if _res == "" {
			continue
		}
		if _res == "data: [DONE]" {
			break
		}
		var _Tooljson realChatToolStreamMessage
		e := json.Unmarshal([]byte(_res), &_Tooljson)
		if e == nil {
			if _Tooljson.Choices[0].FinishReason == "tool_calls" {
				messages = append(messages, ToolMessage{
					Role:      _Tooljson.Choices[0].Message.Role,
					Content:   _Tooljson.Choices[0].Message.Content,
					ToolCalls: _Tooljson.Choices[0].Message.ToolCalls,
				})
				for _, tool_item := range _Tooljson.Choices[0].Message.ToolCalls {
					if tool_item.Type == "function" {
						var arguments_data map[string]interface{}
						e := json.Unmarshal([]byte(tool_item.Function.Arguments), &arguments_data)
						if e != nil {
							return e
						}
						toolRes := functionMap[tool_item.Function.Name](arguments_data)
						messages = append(messages, ToolMessage{
							Role:       "tool",
							ToolCallID: tool_item.ID,
							Content:    toolRes,
						})
					}
				}
				return client.ChatWithTools(model, messages, toolInfo, functionMap, during)
			}
		}
		var _json realChatToolStreamMessage
		e = json.Unmarshal([]byte(_res), &_json)
		if e != nil {
			return e
		}
		during(_json.Choices[0].Message.Content)
	}
	return nil
}
