package user

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"

	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 重置密码
func ResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.ResetPasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		userService := provider.InitializeUserService(svcCtx)
		resp, err := userService.ResetPassword(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
