package user

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"

	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		userService := provider.InitializeUserService(svcCtx)
		resp, err := userService.Login(r.Context(), &req)
		response.Response(r, w, resp, err)

	}
}
