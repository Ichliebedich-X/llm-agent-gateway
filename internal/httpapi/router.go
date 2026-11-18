package httpapi

import "net/http"

/*
路由（Router）
可以把它理解成一个“请求分发员”：
客户端访问 GET /health，就交给 healthHandler
客户端访问 GET /ready，就交给 readyHandler
客户端访问不存在的路径，就返回 404
Go 标准库里负责这件事的组件叫 http.ServeMux。它本质上就是一张路由表
*/

// 创建路由
func NewRouter() *http.ServeMux {
	//NewServeMUx()创建路由表，mux是表的变量名
	mux := http.NewServeMux()

	//注册路由
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /ready", readyHandler)
	/*
		HandleFunc：注册一个请求处理规则
		"GET /health"：规定请求方法必须是 GET，路径必须是 /health；healthHandler：匹配成功后执行的函数
		所以这行代码相当于告诉 Go：“有人通过 GET 访问 /health，就调用 healthHandler”
	*/

	//返回路由表
	return mux
}
