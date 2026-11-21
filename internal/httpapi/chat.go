package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"llm-agent-gateway/internal/llm"
)

const (
	maxRequestBodySize = 1 << 20 // 2^20 = 1048576Bytes = 1 MB
	maxTokensLimit     = 4096
)

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

// ChatGateway 是 Handler 所需的最小业务接口
type ChatGateway interface {
	Chat(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error)
}

// NewChatHandler 创建聊天接口 Handler
func NewChatHandler(gw ChatGateway) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeError(w, http.StatusMethodNotAllowed,
				"method_not_allowed", "only POST is allowed")
			return
		}

		// 限制客户端最多只能提交 1MB的请求 Body,避免内存损耗
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

		var input ChatCompletionRequest

		//限制未知字段，返回解析错误
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&input); err != nil {
			var maxBytesErr *http.MaxBytesError

			if errors.As(err, &maxBytesErr) {
				writeError(
					w,
					http.StatusRequestEntityTooLarge,
					"request_too_large",
					"request body too large",
				)
				return
			}

			writeError(
				w,
				http.StatusBadRequest,
				"invalid_json",
				"invalid JSON request body",
			)
			return
		}

		// 确保请求 Body 中只有一个 JSON 对象
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			writeError(w, http.StatusBadRequest,
				"invalid_json", "request body must contain only one JSON object")
			return
		}

		// 校验模型名称
		if strings.TrimSpace(input.Model) == "" {
			writeError(w, http.StatusBadRequest,
				"invalid_request", "model must not be empty")
			return
		}

		// 校验消息数组
		if len(input.Messages) == 0 {
			writeError(w, http.StatusBadRequest,
				"invalid_request", "messages must not be empty")
			return
		}

		// 校验 max_tokens
		if input.MaxTokens <= 0 || input.MaxTokens > maxTokensLimit {
			writeError(w, http.StatusBadRequest,
				"invalid_request", "max_tokens must be between 1 and 4096")
			return
		}

		// 将 HTTP 消息转换成内部消息
		messages := make([]llm.Message, 0, len(input.Messages))

		for _, msg := range input.Messages {
			role := llm.MessageRole(msg.Role)

			switch role {
			case llm.RoleSystem, llm.RoleUser,
				llm.RoleAssistant, llm.RoleTool:
			default:
				writeError(w, http.StatusBadRequest,
					"invalid_request", "message role is invalid")
				return
			}

			if strings.TrimSpace(msg.Content) == "" {
				writeError(w, http.StatusBadRequest,
					"invalid_request", "message content must not be empty")
				return
			}

			messages = append(messages, llm.Message{
				Role:    role,
				Content: msg.Content,
			})
		}

		// 构造内部业务请求
		req := llm.ChatRequest{
			Model:     input.Model,
			Messages:  messages,
			MaxTokens: input.MaxTokens,
		}

		// 调用 Gateway
		resp, err := gw.Chat(r.Context(), req)
		if err != nil {
			log.Printf("chat gateway failed: %v", err)

			if errors.Is(err, context.DeadlineExceeded) {
				writeError(w, http.StatusGatewayTimeout,
					"provider_timeout", "provider request timed out")
				return
			}

			writeError(w, http.StatusServiceUnavailable,
				"provider_unavailable", "provider is temporarily unavailable")
			return
		}

		// 将内部响应转换为 HTTP 响应
		output := ChatCompletionResponse{
			Provider: resp.Provider,
			Model:    resp.Model,
			Content:  resp.Content,
			Usage: HTTPUsage{
				InputTokens:  resp.Usage.InputTokens,
				OutputTokens: resp.Usage.OutputTokens,
			},
		}

		writeJSON(w, http.StatusOK, output)
	})
}
