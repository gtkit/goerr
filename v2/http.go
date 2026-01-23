package goerr

import (
	"encoding/json"
	"net/http"
)

// HTTPResponse HTTP 错误响应结构
type HTTPResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	TraceID string                 `json:"trace_id,omitempty"`
}

// WriteHTTPError 将错误写入 HTTP 响应
func WriteHTTPError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	var resp HTTPResponse
	var status int

	// 检查是否是我们的错误类型
	if e, ok := err.(Error); ok {
		resp = HTTPResponse{
			Code:    e.Code(),
			Message: e.Message(),
			Details: e.Details(),
		}
		status = e.HTTPStatus()
	} else {
		// 标准错误
		resp = HTTPResponse{
			Code:    ErrInternal,
			Message: err.Error(),
		}
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// 编码 JSON
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(resp)
}

// HTTPMiddleware HTTP 中间件,统一处理 panic 和错误
func HTTPMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch x := rec.(type) {
				case string:
					err = New(ErrInternal, x)
				case error:
					err = Wrap(x, ErrInternal, "panic recovered")
				default:
					err = New(ErrInternal, "unknown panic")
				}
				WriteHTTPError(w, err)
			}
		}()

		next(w, r)
	}
}

// ResponseWriter 包装的响应写入器,支持错误捕获
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// NewResponseWriter 创建响应写入器
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader 记录状态码
func (w *ResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.statusCode = code
		w.ResponseWriter.WriteHeader(code)
		w.written = true
	}
}

// Write 写入数据
func (w *ResponseWriter) Write(data []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

// StatusCode 获取状态码
func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}
