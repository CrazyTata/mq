package tool

import (
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
	"net/http"
)

// UploadTokenHandler 处理获取用户课程的请求
func UploadTokenHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		toolService := provider.InitializeToolService(svc)
		token, err := toolService.UploadToken(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Response(r, w, token, nil)
	}
}
