package user

import (
	"mq/common/response"

	"encoding/json"
	"mq/application/dto"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
	"net/http"
)

// LoginByTokenHandler handles user login requests
func LoginByTokenHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginByTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userService := provider.InitializeUserService(svc)
		resp, err := userService.LoginByToken(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
