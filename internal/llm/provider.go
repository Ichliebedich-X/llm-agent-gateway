package llm

import "context"

// Provider 定义模型供应商需要实现的统一能力
type Provider interface {
	//Name 返回 Provider的名称
	Name() string
	//Chat 发送对话请求并返回模型响应
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
}
