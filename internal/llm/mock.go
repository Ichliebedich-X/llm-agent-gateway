package llm

import (
	"context"
	"fmt"
)

// MockProvider 是一个不访问真实模型的模拟实现
type MockProvider struct{}

// Name 返回Provider名称
func (m MockProvider) Name() string {
	return "mock"
}

// Chat 模拟模型对话调用
func (m MockProvider) Chat(
	ctx context.Context,
	req ChatRequest,
) (ChatResponse, error) { //同时返回结果和错误：成功时返回响应和 nil；失败时返回空响应和具体错误

	//检查请求是否已经被取消，如果上游取消了请求，mock就不再继续返回成功结果
	if err := ctx.Err(); err != nil {
		return ChatResponse{}, err
	}

	//获取用户最后一条消息
	var userMessage string

	for i := len(req.Message) - 1; i >= 0; i-- {
		if req.Message[i].Role == RoleUser {
			userMessage = req.Message[i].Content
			break
		}
	}

	return ChatResponse{
		Content:  fmt.Sprintf("Mock reply：收到的消息：%s", userMessage),
		Provider: m.Name(),
		Model:    req.Model,
		Usage:    Usage{},
	}, nil
}

// 编译期检查：确保 MockProvider实现了 Provider
var _ Provider = MockProvider{}
