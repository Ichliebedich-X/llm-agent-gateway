package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"llm-agent-gateway/internal/llm"
)

// fakeChatGateway 是专门用于测试的 Gateway 替身
// 它实现了 ChatGateway 接口，不需要调用真实模型
type fakeChatGateway struct{}

// Chat 模拟一次模型对话，返回固定响应
func (f fakeChatGateway) Chat(
	ctx context.Context,
	req llm.ChatRequest,
) (llm.ChatResponse, error) {
	return llm.ChatResponse{
		Provider: "test",
		Model:    "test-model",
		Content:  "test response",
	}, nil
}

// newTestRouter 创建测试使用的 Router
func newTestRouter() *http.ServeMux {
	return NewRouter(fakeChatGateway{})
}

// 测试 health 接口
func TestHealthHandler(t *testing.T) {
	// 1. 创建路由
	router := newTestRouter()

	// 2. 创建模拟响应记录器
	recorder := httptest.NewRecorder()

	// 3. 创建模拟 GET 请求
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	// 4. 将请求交给路由处理
	router.ServeHTTP(recorder, request)

	// 5. 检查响应状态码
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, recorder.Code)
	}

	// 6. 检查响应 Body 是否为合法 JSON
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// 7. 检查 JSON 中的 status 是否为 ok
	if response["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", response["status"])
	}

	// 8. 检查响应头 Content-Type
	contentType := recorder.Header().Get("Content-Type")

	if contentType != "application/json; charset=utf-8" {
		t.Fatalf(
			"expected Content-Type application/json; charset=utf-8, got %q",
			contentType,
		)
	}
}

// 测试 ready 接口
func TestReadyHandler(t *testing.T) {
	// 1. 创建路由
	router := newTestRouter()

	// 2. 创建模拟响应记录器
	recorder := httptest.NewRecorder()

	// 3. 创建模拟 GET 请求
	request := httptest.NewRequest(http.MethodGet, "/ready", nil)

	// 4. 将请求交给路由处理
	router.ServeHTTP(recorder, request)

	// 5. 检查响应状态码
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, recorder.Code)
	}

	// 6. 检查响应 Body
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", response["status"])
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Fatalf(
			"expected Content-Type application/json; charset=utf-8, got %q",
			contentType,
		)
	}
}

// 测试 404
func TestNotFound(t *testing.T) {
	// 1. 创建路由
	router := newTestRouter()

	// 2. 创建模拟响应记录器
	recorder := httptest.NewRecorder()

	// 3. 创建一个不存在的接口请求
	request := httptest.NewRequest(http.MethodGet, "/not-found", nil)

	// 4. 交给路由处理
	router.ServeHTTP(recorder, request)

	// 5. 检查状态码是否为 404
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d",
			http.StatusNotFound, recorder.Code)
	}
}

// 测试 405
func TestMethodNotAllowed(t *testing.T) {
	// 1. 创建路由
	router := newTestRouter()

	// 2. 创建模拟 POST 请求
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/health", nil)

	// 3. 交给路由处理
	router.ServeHTTP(recorder, request)

	// 4. 检查状态码
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d",
			http.StatusMethodNotAllowed, recorder.Code)
	}
}
