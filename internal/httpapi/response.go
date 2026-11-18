package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON 负责把 Go 数据结构编码成 JSON，并写入 HTTP 响应
func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("failed to marshal JSON response: %v", err)

		status = http.StatusInternalServerError
		data = []byte(`{"error":{"code":"internal_error","message":"internal server error"}}`)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		log.Printf("failed to write HTTP response: %v", err)
	}
}

// writeError 负责组装统一的错误结构，再交给 writeJSON
func writeError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	response := ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	}

	writeJSON(w, status, response)
}
