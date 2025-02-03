package dashscope

import (
	"encoding/json"
	"errors"
)

type Text2ImageRequest struct {
	// 模型名，可上：https://help.aliyun.com/zh/model-studio/developer-reference/text-to-image-v2-api-reference 查看
	Model string `json:"model"`

	// 两个提示词字数均不能超过500
	Input      Text2ImageRequestInput     `json:"input"`
	Parameters Text2ImageRequestParameter `json:"parameters"`
}

type Text2ImageRequestInput struct {
	// 正向提示词
	Prompt string `json:"prompt"`

	// 反向提示词
	NegativePrompt string `json:"negative_prompt"`
}

type Text2ImageRequestParameter struct {
	// 图片大小，例如1024*1024
	Size string `json:"size"`

	// 生成图片个数
	N int `json:"n"`

	// 是否开启prompt智能改写。开启后会使用大模型对输入prompt进行智能改写
	PromptExtend bool `json:"prompt_extend"`

	// 是否添加水印标识，水印位于图片右下角，文案为“AI生成”
	Watermark bool `json:"watermark"`
}

type text2imageSuccess struct {
	Output struct {
		TaskStatus string `json:"task_status"`
		TaskID     string `json:"task_id"`
	} `json:"output"`
	RequestID string `json:"request_id"`
}

type text2imageError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type Text2imageResponse struct {
	RequestID string `json:"request_id"`
	Output    struct {
		TaskID        string `json:"task_id"`
		TaskStatus    string `json:"task_status"`
		SubmitTime    string `json:"submit_time"`
		ScheduledTime string `json:"scheduled_time"`
		EndTime       string `json:"end_time"`
		Code          string `json:"code"`
		Message       string `json:"message"`
		Results       []struct {
			OrigPrompt   string `json:"orig_prompt"`
			ActualPrompt string `json:"actual_prompt"`
			URL          string `json:"url"`
		} `json:"results"`
		TaskMetrics struct {
			Total     int `json:"TOTAL"`
			Succeeded int `json:"SUCCEEDED"`
			Failed    int `json:"FAILED"`
		} `json:"task_metrics"`
	} `json:"output"`
	Usage struct {
		ImageCount int `json:"image_count"`
	} `json:"usage"`
}

// 调用文生图接口，成功返回task id
func (client *Client) Text2image(body Text2ImageRequest) (string, error) {
	if body.Parameters.N == 0 {
		body.Parameters.N = 1
	}
	if body.Parameters.Size == "" {
		body.Parameters.Size = "1024*1024"
	}

	httpClient := client.newHttpClient()
	httpClient.SetHeader("X-DashScope-Async", "enable")
	bodyText, e := json.Marshal(body)
	if e != nil {
		return "", e
	}
	httpClient.SetBody(bodyText)
	httpres, e := httpClient.Post("https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis")
	if e != nil {
		return "", e
	}
	if httpres.StatusCode() != 200 {
		var resError text2imageError
		e := json.Unmarshal(httpres.Body(), &resError)
		if e != nil {
			return "", errors.New(string(httpres.Body()))
		}
		return "", errors.New(resError.Message)
	}
	var resSuccess text2imageSuccess
	json.Unmarshal(httpres.Body(), &resSuccess)
	return resSuccess.Output.TaskID, nil
}

func (client *Client) Text2imageResult(taskID string) (*Text2imageResponse, error) {
	httpClient := client.newHttpClient()
	httpres, e := httpClient.Get("https://dashscope.aliyuncs.com/api/v1/tasks/" + taskID)
	if e != nil {
		return nil, e
	}
	if httpres.StatusCode() != 200 {
		var resError text2imageError
		e := json.Unmarshal(httpres.Body(), &resError)
		if e != nil {
			return nil, errors.New(string(httpres.Body()))
		}
		return nil, errors.New(resError.Message)
	}
	var res Text2imageResponse
	json.Unmarshal(httpres.Body(), &res)
	return &res, nil
}
