package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 模拟 health 接口测试
func TestHealthHandler(t *testing.T) {
	//1.创建路由
	router := NewRouter()

	//2.创建模拟响应记录器
	//虚拟客户端收件箱，记录 Handler 返回的状态码、Header 和 Body
	recorder := httptest.NewRecorder()

	//3.创建模拟GET请求
	//模拟客户端发送 health请求
	request := httptest.NewRequest(http.MethodGet, "/health", nil) //nil表示没带请求体 body

	//4.将请求交给路由处理
	//把请求交给路由表，让它根据路径找到对应的 Handler，然后执行
	router.ServeHTTP(recorder, request)

	//5.检查响应状态码
	//预期返回值200，否则测试失败并打印预期值和实际值
	//t是go测试框架传进来的测试对象，Fatalf表示报告错误并立即结束当前测试
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d,got %d",
			http.StatusOK, recorder.Code)
	}

	//6.检查响应 body是否为合法的JSON
	var response map[string]string
	//json.Unmarshal()负责把 JSON 数据解析成 Go 变量
	//recorder.Body 记录的是 Handler 返回的响应正文；bytes()将正文转换成字节切片方便后续处理
	//解析之后 &response 取变量的地址，因为 Unmarshal 需要通过指针把解析结果写入变量
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response body:%v", err)
	}

	//7.检查JSON中的 status是否为ok
	if response["status"] != "ok" {
		t.Fatalf("expected status ok,got %q", response["status"])
	}

	//8.检查响应头 Content-Type
	//recorder.Header() 获取 Handler 返回的 HTTP 响应头
	//.Get("Content-Type") 则是从响应头中读取 Content-Type 的值
	//确认返回客户端的是JSON格式
	contentType := recorder.Header().Get("Content-Type")

	if contentType != "application/json; charset=utf-8" {
		t.Fatalf("expected Content-Type application/json; charset=utf-8, got %q", contentType)
	}
}

// 模拟 ready接口测试
func TestReadyHandler(t *testing.T) {
	// 1. 创建路由
	router := NewRouter()

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
}

// 404测试
func TestNotFound(t *testing.T) {
	// 1. 创建路由
	router := NewRouter()

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

// 405测试
func TestMethodNotAllowed(t *testing.T) {
	// 1. 创建路由
	router := NewRouter()

	// 2. 创建模拟响应记录器
	recorder := httptest.NewRecorder()

	// 3. 创建 POST 请求
	request := httptest.NewRequest(http.MethodPost, "/health", nil)

	// 4. 交给路由处理
	router.ServeHTTP(recorder, request)

	// 5. 检查状态码是否为 405
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d",
			http.StatusMethodNotAllowed, recorder.Code)
	}
}
