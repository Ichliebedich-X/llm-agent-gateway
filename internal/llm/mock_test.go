package llm

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestMockProviderChat(t *testing.T) {
	//1.通过 Provider接口使用 MockProvider
	var provider Provider = MockProvider{}

	//2.构造请求
	req := ChatRequest{
		Model: "mock-model",
		Message: []Message{
			{
				Role:    RoleUser,
				Content: "你好",
			},
		},
		MaxTokens: 100,
	}

	//3.调用 Chat
	resp, err := provider.Chat(context.Background(), req)

	//4.检查错误
	if err != nil {
		t.Fatalf("Chat returned error:%v", err)
	}

	//5.检查响应
	if resp.Provider != "mock" {
		t.Errorf("excepted provider mock got %s", resp.Provider)
	}

	if resp.Model != "mock-model" {
		t.Errorf("excepted model mock-model,got %s", resp.Model)
	}

	if !strings.Contains(resp.Content, "你好") {
		t.Errorf("response dose not contain user message:%s", resp.Content)
	}
}

func TestMockProviderCanceledContext(t *testing.T) {
	provider := MockProvider{}

	//创建一个已经超时的 context
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Nanosecond,
	)
	defer cancel()

	// 等待 context 取消或超时，然后再调用 Mock
	//这样测试就不依赖机器运行速度，不需要使用 time.Sleep 猜测取消是否已经发生
	<-ctx.Done()

	req := ChatRequest{
		Model: "mock-model",
	}

	_, err := provider.Chat(ctx, req)

	if err == nil {
		t.Fatal("expected error when context is canceled")
	}
}
