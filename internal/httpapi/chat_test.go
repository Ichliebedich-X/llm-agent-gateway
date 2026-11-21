package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"llm-agent-gateway/internal/llm"
)

// testChatGateway 是 Chat Handler 测试专用的 Gateway
type testChatGateway struct {
	response llm.ChatResponse
	err      error
	called   bool
	lastReq  llm.ChatRequest
}

func (g *testChatGateway) Chat(
	ctx context.Context,
	req llm.ChatRequest,
) (llm.ChatResponse, error) {
	g.called = true
	g.lastReq = req
	return g.response, g.err
}

// 创建合法的测试请求
func validChatRequest() string {
	return `{
		"model": "mock",
		"messages": [
			{"role": "user", "content": "你好"}
		],
		"max_tokens": 128
	}`
}

// 测试正常请求
func TestChatHandlerSuccess(t *testing.T) {
	gw := &testChatGateway{
		response: llm.ChatResponse{
			Provider: "mock",
			Model:    "mock-v1",
			Content:  "你好，我是 Mock。",
			Usage: llm.Usage{
				InputTokens:  5,
				OutputTokens: 8,
			},
		},
	}

	handler := NewChatHandler(gw)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(validChatRequest()),
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s",
			recorder.Code, recorder.Body.String())
	}

	if !gw.called {
		t.Fatal("expected gateway to be called")
	}

	if gw.lastReq.Model != "mock" {
		t.Errorf("expected model mock, got %q", gw.lastReq.Model)
	}

	if len(gw.lastReq.Messages) != 1 ||
		gw.lastReq.Messages[0].Content != "你好" {
		t.Errorf("unexpected messages: %+v", gw.lastReq.Messages)
	}

	if recorder.Header().Get("Content-Type") !=
		"application/json; charset=utf-8" {
		t.Errorf("unexpected Content-Type: %q",
			recorder.Header().Get("Content-Type"))
	}

	if !strings.Contains(recorder.Body.String(), "mock-v1") {
		t.Errorf("response does not contain model: %s",
			recorder.Body.String())
	}
}

// 测试非法 JSON
func TestChatHandlerInvalidJSON(t *testing.T) {
	gw := &testChatGateway{}
	handler := NewChatHandler(gw)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(`{"model":`),
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	if gw.called {
		t.Fatal("gateway should not be called for invalid JSON")
	}
}

// 测试 messages 为空
func TestChatHandlerEmptyMessages(t *testing.T) {
	gw := &testChatGateway{}
	handler := NewChatHandler(gw)

	body := `{
		"model": "mock",
		"messages": [],
		"max_tokens": 128
	}`

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(body),
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	if gw.called {
		t.Fatal("gateway should not be called when messages are empty")
	}
}

// 测试非法消息角色
func TestChatHandlerInvalidRole(t *testing.T) {
	gw := &testChatGateway{}
	handler := NewChatHandler(gw)

	body := `{
		"model": "mock",
		"messages": [
			{"role": "unknown", "content": "你好"}
		],
		"max_tokens": 128
	}`

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(body),
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	if gw.called {
		t.Fatal("gateway should not be called for invalid role")
	}
}

// 测试超出最大请求 Body 限制
func TestChatHandlerRequestTooLarge(t *testing.T) {
	gw := &testChatGateway{}
	handler := NewChatHandler(gw)

	body := `{"model":"mock","messages":[{"role":"user","content":"` +
		strings.Repeat("a", maxRequestBodySize) +
		`"}],"max_tokens":128}`

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(body),
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", recorder.Code)
	}

	if gw.called {
		t.Fatal("gateway should not be called for oversized body")
	}
}

// 测试 Gateway 返回普通错误
func TestChatHandlerGatewayError(t *testing.T) {
	gw := &testChatGateway{
		err: errors.New("provider connection failed"),
	}
	handler := NewChatHandler(gw)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(validChatRequest()),
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", recorder.Code)
	}

	if !gw.called {
		t.Fatal("expected gateway to be called")
	}

	// 不应该把内部错误原文暴露给客户端
	if strings.Contains(recorder.Body.String(), "provider connection failed") {
		t.Fatal("internal error should not be exposed to client")
	}
}

// 测试 Gateway 超时
func TestChatHandlerGatewayTimeout(t *testing.T) {
	gw := &testChatGateway{
		err: context.DeadlineExceeded,
	}
	handler := NewChatHandler(gw)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(validChatRequest()),
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected status 504, got %d", recorder.Code)
	}
}

// 测试非 POST 请求
func TestChatHandlerMethodNotAllowed(t *testing.T) {
	gw := &testChatGateway{}
	handler := NewChatHandler(gw)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/v1/chat/completions",
		nil,
	)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}

	if recorder.Header().Get("Allow") != http.MethodPost {
		t.Errorf("expected Allow: POST, got %q",
			recorder.Header().Get("Allow"))
	}

	if gw.called {
		t.Fatal("gateway should not be called for GET request")
	}
}
