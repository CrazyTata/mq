package user

import (
	"encoding/json"
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
)

// UpdateUserHandler handles user information update requests
func UpdateUserHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userService := provider.InitializeUserService(svc)
		resp, err := userService.UpdateUser(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
