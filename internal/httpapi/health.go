package httpapi

import (
	"encoding/json"
	"net/http"
)

// 检查服务进程是否存活（存活检查）
// 比如：「这个 Go 进程还活着吗？」
// 即使 MySQL 暂时连接失败，只要 Go 进程还能够正常工作，通常也应该返回 200。
func healthHandler(w http.ResponseWriter, r *http.Request) {
	//w：负责向客户端写响应。比如状态码、响应头、JSON 数据；r：负责接收客户端发来的请求。比如请求方法、URL、请求头、请求体
	//设置响应类型为JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//设置状态码：200
	w.WriteHeader(http.StatusOK)

	//创建 map类型的响应数据
	//把 map转换成：
	//{"status":"ok"}
	response := map[string]string{
		"status": "ok",
	}

	//将 Go数据转换成JSON并写入响应
	json.NewEncoder(w).Encode(response)
}

// 检查服务是否已经准备好接受业务请求（就绪检查）
// 比如：「这个服务现在是否具备接收业务请求的条件？」
// 以后接入 MySQL、Redis、模型供应商等依赖后，就可以根据服务的实际就绪条件进行判断。
func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(response)
}
