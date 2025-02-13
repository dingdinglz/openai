package openai

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type VisionMessage struct {
	Content []VisionContent `json:"content"`
	Role    string          `json:"role"`
}

type VisionContent struct {
	// VISION_MESSAGE_
	Type     string                 `json:"type"`
	Text     string                 `json:"text,omitempty"`
	ImageUrl *VisionContentImageUrl `json:"image_url,omitempty"`
}

type VisionContentImageUrl struct {
	Url string `json:"url"`
	// 图片细节：
	Detail string `json:"detail,omitempty"`
}

type VisionChatRequest struct {
	Messages         []VisionMessage `json:"messages"`
	Model            string          `json:"model"`
	FrequencyPenalty int             `json:"frequency_penalty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	PresencePenalty  int             `json:"presence_penalty"`
	ResponseFormat   struct {
		Type string `json:"type"`
	} `json:"response_format"`
	Stop        []string `json:"stop"`
	Stream      bool     `json:"stream"`
	Temperature float32  `json:"temperature"`
	TopP        int      `json:"top_p"`
}

type realVisionChatStreamError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func checkVisionChatRequest(cr *VisionChatRequest) {
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

func GenerateImageUrlBase64(file []byte) string {
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(file)
}

func (client Client) ChatVisionStream(model string, messages []VisionMessage, during func(string)) error {
	reqBody := VisionChatRequest{}
	reqBody.Messages = messages
	reqBody.Stream = true
	reqBody.Model = model
	checkVisionChatRequest(&reqBody)
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
			var resError realVisionChatStreamError
			json.Unmarshal([]byte(_res), &resError)
			if resError.Code != 0 {
				return errors.New(resError.Message)
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
