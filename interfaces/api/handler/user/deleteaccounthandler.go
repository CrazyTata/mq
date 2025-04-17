package user

import (
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
	"net/http"
)

// DeleteAccountHandler handles user account deletion requests
func DeleteAccountHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userService := provider.InitializeUserService(svc)
		resp, err := userService.DeleteAccount(r.Context())
		response.Response(r, w, resp, err)
	}
}
