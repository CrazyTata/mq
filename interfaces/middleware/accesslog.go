package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// AccessLog 是记录API访问日志的中间件
var AccessLog = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		raw := r.URL.RawQuery

		// 如果有查询参数，添加到路径中
		if raw != "" {
			path = path + "?" + raw
		}

		// 获取客户端真实IP
		clientIP := r.Header.Get("X-Real-IP")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Forwarded-For")
		}
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		// 保存请求体以备日志记录，对于敏感接口可以选择不记录请求体
		var reqBody []byte
		if r.Method == "POST" || r.Method == "PUT" {
			// 判断是否是敏感接口，如果是则不记录请求体
			sensitiveAPI := false
			sensitiveAPIs := []string{"/api/payment", "/api/user/login", "/api/user/register"}
			for _, api := range sensitiveAPIs {
				if path == api {
					sensitiveAPI = true
					break
				}
			}

			if !sensitiveAPI && r.Body != nil {
				reqBody, _ = io.ReadAll(r.Body)
				// 重新设置请求体，因为读取后会被消费
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}
		}

		// 记录请求信息
		logx.Infof("访问: %s | %s | %s | %s", r.Method, path, clientIP, r.UserAgent())

		// 调用下一个处理器
		next(w, r)

		// 计算请求处理时间
		latency := time.Since(start)

		// 记录响应信息，也还是要记录相应体
		logx.Infof("响应: %s | %s | %s | %v", r.Method, path, clientIP, latency)
	}
}
