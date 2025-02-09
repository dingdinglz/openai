package siliconflow

import (
	"encoding/json"
	"errors"
	"math/rand"
)

type Text2ImageRequest struct {
	// 输出图片个数
	BatchSize int `json:"batch_size"`

	// 用于控制生成图像与给定提示（Prompt）的匹配程度，该值越高，生成的图像越倾向于严格匹配文本提示的内容；该值越低，则生成的图像会更加具有创造性和多样性，可能包含更多的意外元素。
	GuidanceScale float64 `json:"guidance_scale"`

	// 图像尺寸，格式：[宽度]x[高度]
	ImageSize string `json:"image_size"`

	// 模型名
	Model string `json:"model"`

	// 推理步骤数，其中stable-diffusion-3-5-large-turbo是固定值4
	NumInferenceSteps int `json:"num_inference_steps"`

	// 提示词
	Prompt string `json:"prompt"`
	// 负向提示词
	NegativePrompt string `json:"negative_prompt"`

	// 提示增强开关，当开启时，将提示优化为详细的、模型友好的版本
	PromptEnhancement bool `json:"prompt_enhancement"`

	// 随机种子
	Seed int64 `json:"seed"`
}

type text2ImageResponse struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	Timings struct {
		Inference float64 `json:"inference"`
	} `json:"timings"`
	Seed     int64  `json:"seed"`
	SharedID string `json:"shared_id"`
	Data     []struct {
		URL string `json:"url"`
	} `json:"data"`
	Created int `json:"created"`
}

type text2ImageError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func checkText2ImageRequest(r *Text2ImageRequest) {
	if r.BatchSize == 0 {
		r.BatchSize = 1
	}
	if r.GuidanceScale == 0 {
		r.GuidanceScale = 7.5
	}
	if r.ImageSize == "" {
		r.ImageSize = "1024x1024"
	}
	if r.NumInferenceSteps == 0 {
		r.NumInferenceSteps = 20
	}
	if r.Seed == 0 {
		r.Seed = rand.Int63n(9999999999)
	}
}

// 调用文生图接口，成功返回生成的图片的url数组
func (client *Client) Text2Image(req Text2ImageRequest) ([]string, error) {
	httpReq := client.newHttpClient()
	checkText2ImageRequest(&req)
	body, e := json.Marshal(req)
	if e != nil {
		return []string{}, e
	}
	httpReq.SetBody(body)
	httpRes, e := httpReq.Post("https://api.siliconflow.cn/v1/images/generations")
	if e != nil {
		return []string{}, e
	}
	if httpRes.StatusCode() != 200 {
		var errMessage text2ImageError
		e := json.Unmarshal(httpRes.Body(), &errMessage)
		if e != nil {
			return []string{}, errors.New(string(httpRes.Body()))
		}
		return []string{}, errors.New(errMessage.Message)
	}
	var resData text2ImageResponse
	e = json.Unmarshal(httpRes.Body(), &resData)
	if e != nil {
		return []string{}, e
	}
	var resUrls []string
	for _, i := range resData.Images {
		resUrls = append(resUrls, i.URL)
	}
	return resUrls, nil
}
