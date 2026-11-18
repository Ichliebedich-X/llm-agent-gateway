package httpapi

// ChatMessage 表示 HTTP 请求中的一条消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest 表示客户端发起聊天请求时提交的数据
type ChatCompletionRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
}

// ChatCompletionResponse 表示聊天接口返回给客户端的数据
type ChatCompletionResponse struct {
	Provider string    `json:"provider"`
	Model    string    `json:"model"`
	Content  string    `json:"content"`
	Usage    HTTPUsage `json:"usage"`
}

// HTTPUsage 表示本次请求的 Token 用量
type HTTPUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ErrorResponse 表示统一的错误响应
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError 表示具体的错误信息
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
