package user

import (
	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// BackendUserHandler 处理获取后端用户列表的请求
func BackendUserHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.BackendUserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 调用服务
		userService := provider.InitializeUserService(svc)
		resp, err := userService.GetBackendUsers(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
