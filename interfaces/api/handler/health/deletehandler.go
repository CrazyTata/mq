package health

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// DeleteHandler handles requests to delete a health record
func DeleteHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.PostIdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		healthService := provider.InitializeHealthService(svc)
		err := healthService.Delete(r.Context(), req.ID)
		response.Response(r, w, nil, err)
	}
}
