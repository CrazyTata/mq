package user

import (
	"encoding/json"
	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
	"net/http"
)

// LoginByMessageHandler handles user login requests
func LoginByMessageHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginByMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userService := provider.InitializeUserService(svc)
		resp, err := userService.LoginByMessage(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
