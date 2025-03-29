package model

// ALiAiReq 阿里内容生成入参
type ALiAiReq struct {
	Model      string        `json:"model"`      // 指定用于对话的通义千问模型名，目前可选择qwen-turbo、qwen-plus、qwen-max、qwen-max-0403、qwen-max-0107、qwen-max-1201和qwen-max-longcontext。
	Messages   []*AiMessage  `json:"messages"`   // 文本信息
	Parameters *AiParameters `json:"parameters"` // 视觉模型参数
	Stream     bool          `json:"stream"`
}

// AiParameters 视觉模型参数
type AiParameters struct {
	ALiReqParameters
	PresencePenalty float64           `json:"presence_penalty"` // 控制模型生成文本时的内容重复度。取值范围：[-2.0, 2.0]。正数会减少重复度，负数会增加重复度
	ResponseFormat  ALiResponseFormat `json:"response_format"`  // 返回内容的格式
}

// ALiResponseFormat 返回内容的格式
type ALiResponseFormat struct {
	Type string `json:"type"` // 返回内容的格式 可选值：{"type": "text"}或{"type": "json_object"}。设置为{"type": "json_object"}时会输出标准格式的JSON字符串。
}

// ALiReqParameters 用于控制模型生成的参数
type ALiReqParameters struct {
	ResultFormat      string  `json:"result_format"`      // 用于指定返回结果的格式，默认为text，也可设置为message。当设置为message时，输出格式请参考返回结果。推荐优先使用message格式。
	IncrementalOutput bool    `json:"incremental_output"` // 控制在流式输出模式下是否开启增量输出，即后续输出内容是否包含已输出的内容。设置为True时，将开启增量输出模式，后面输出不会包含已经输出的内容，您需要自行拼接整体输出；设置为False则会包含已输出的内容
	MaxTokens         int32   `json:"max_tokens"`         // 用于限制模型生成token的数量，表示生成token个数的上限。其中qwen-turbo最大值和默认值为1500，qwen-max、qwen-max-1201 、qwen-max-longcontext 和 qwen-plus最大值和默认值均为2000。
	Temperature       float64 `json:"temperature"`        // 要使用的采样温度，介于 0 和 2 之间。
	TopP              float64 `json:"top_p"`              // 温度采样的替代方法，称为核心采样，其中模型将考虑具有 top_p 概率质量的令牌的结果。
}

// AiMessage 文本信息
type AiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AiContent 内容
type AiContent struct {
	Type     string `json:"type"`      // text或者image_url
	Text     string `json:"text"`      // type为text时填入该字段
	ImageUrl AiImg  `json:"image_url"` // type为image_url时填入该字段
}

// AiImg 图片
type AiImg struct {
	Url string `json:"url"`
}

// ALiOpenStreamUsage ali文案生成流式结果usage
type ALiOpenStreamUsage struct {
	TotalTokens  int32 `json:"total_tokens"`
	InputTokens  int32 `json:"input_tokens"`
	OutputTokens int32 `json:"output_tokens"`
}

// ALiStreamRsp 阿里流式返回
// {"choices":[{"delta":{"content":"","role":"assistant"},"index":0,"logprobs":null,"finish_reason":null}],"object":"chat.completion.chunk","usage":null,"created":1743265691,"system_fingerprint":null,"model":"qwen-plus-2025-01-25","id":"chatcmpl-5606b696-0150-9406-bb20-06ac69b0ed5c"}
//
//	{
//	 "choices" : [ {
//	   "finish_reason" : "stop",
//	   "delta" : {
//	     "content" : ""
//	   },
//	   "index" : 0,
//	   "logprobs" : null
//	 } ],
//	 "object" : "chat.completion.chunk",
//	 "usage" : null,
//	 "created" : 1743265691,
//	 "system_fingerprint" : null,
//	 "model" : "qwen-plus-2025-01-25",
//	 "id" : "chatcmpl-5606b696-0150-9406-bb20-06ac69b0ed5c"
//	}
type ALiStreamRsp struct {
	Choices           []*ALiStreamChoice `json:"choices"`
	Object            string             `json:"object"`
	Usage             interface{}        `json:"usage"`
	Created           int                `json:"created"`
	SystemFingerprint interface{}        `json:"system_fingerprint"`
	Model             string             `json:"model"`
	Id                string             `json:"id"`
}

type ALiStreamChoice struct {
	FinishReason interface{}     `json:"finish_reason"`
	Delta        *ALiStreamDelta `json:"delta"`
	Index        int             `json:"index"`
}

type ALiStreamDelta struct {
	Content string `json:"content"`
}

// GetContent 获取文案
func (a *ALiStreamRsp) GetContent() string {
	if a == nil || a.Choices == nil || len(a.Choices) <= 0 {
		return ""
	}
	content := ""
	for i := range a.Choices {
		if a.Choices[i].Delta != nil {
			content += a.Choices[i].Delta.Content
		}
	}
	return content
}

// GetFinish 获取是否结束
func (a *ALiStreamRsp) GetFinish() bool {
	if a == nil || a.Choices == nil || len(a.Choices) <= 0 {
		return false
	}
	for i := range a.Choices {
		if a.Choices[i].FinishReason == "stop" {
			return true
		}
	}
	return false
}
