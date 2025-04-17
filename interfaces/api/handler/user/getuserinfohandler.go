package user

import (
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
	"net/http"
)

// GetUserInfoHandler handles requests to get user information
func GetUserInfoHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userService := provider.InitializeUserService(svc)
		resp, err := userService.GetUserInfo(r.Context())
		response.Response(r, w, resp, err)
	}
}
