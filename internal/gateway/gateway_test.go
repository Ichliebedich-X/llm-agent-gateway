package gateway

import (
	"context"
	"errors"
	"testing"

	"llm-agent-gateway/internal/llm"
)

// fakeProvider 是测试专用的 Provider
/*
MockProvider 是项目中的模拟模型实现，可以在开发和演示时使用
fakeProvider 是测试专用的假实现，方便人为控制返回结果和错误
*/
type fakeProvider struct {
	response llm.ChatResponse
	err      error
	called   bool
}

func (f *fakeProvider) Name() string {
	return "fake"
}

func (f *fakeProvider) Chat(
	ctx context.Context,
	req llm.ChatRequest,
) (llm.ChatResponse, error) {
	f.called = true
	return f.response, f.err
}

func TestGatewayChat(t *testing.T) {
	// 准备假 Provider
	fake := &fakeProvider{
		response: llm.ChatResponse{
			Content:  "测试回答",
			Provider: "fake",
			Model:    "test-model",
		},
	}

	// 创建 Gateway
	gw := New(fake)

	// 发起请求
	req := llm.ChatRequest{
		Model: "test-model",
		Messages: []llm.Message{
			{
				Role:    llm.RoleUser,
				Content: "你好",
			},
		},
	}

	resp, err := gw.Chat(context.Background(), req)

	// 检查错误
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 检查是否调用了 Provider
	if !fake.called {
		t.Fatal("expected provider to be called")
	}

	// 检查响应
	if resp.Content != "测试回答" {
		t.Errorf("expected 测试回答, got %s", resp.Content)
	}

	if resp.Provider != "fake" {
		t.Errorf("expected fake, got %s", resp.Provider)
	}
}

func TestGatewayProviderError(t *testing.T) {
	expectedErr := errors.New("provider failed")

	fake := &fakeProvider{
		err: expectedErr,
	}

	gw := New(fake)

	_, err := gw.Chat(
		context.Background(),
		llm.ChatRequest{},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected provider error, got %v", err)
	}
}
