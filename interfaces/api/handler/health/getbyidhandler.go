package health

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetByIDHandler handles requests to get health record by ID
func GetByIDHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.GetIdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		healthService := provider.InitializeHealthService(svc)
		resp, err := healthService.GetByID(r.Context(), req.ID)
		response.Response(r, w, resp, err)
	}
}
