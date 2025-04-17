package user

import (
	"encoding/json"
	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
	"net/http"
)

// LoginByAppleHandler handles user login requests
func LoginByAppleHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginByAppleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userService := provider.InitializeUserService(svc)
		resp, err := userService.LoginByApple(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
