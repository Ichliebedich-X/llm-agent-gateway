package gateway

import (
	"context"

	"llm-agent-gateway/internal/llm"
)

// Gateway 负责模型调用的业务编排
type Gateway struct {
	//provider 并不是某个具体的 MockProvider，而是一个接口类型
	//Gateway不需要知道底层用的是哪个模型
	provider llm.Provider
}

// New 创建 Gateway：Gateway 不负责自己创建 Mock，而是接收外部传入的 Provider
/*
依赖注入（Dependency Injection，DI）：
简单来说，就是一个模块需要什么依赖，由外部提供，而不是自己在内部创建
这样以后想换成 OpenAI Provider，只需要传入另一个实现接口的对象，不必修改 Gateway 的代码
*/
func New(provider llm.Provider) *Gateway {
	return &Gateway{
		provider: provider,
	}
}

// Chat 调用 Provider，完成一次模型对话
func (g *Gateway) Chat(
	ctx context.Context,
	req llm.ChatRequest,
) (llm.ChatResponse, error) {
	return g.provider.Chat(ctx, req)
}
