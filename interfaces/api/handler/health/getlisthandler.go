package health

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetListHandler handles requests to get health record list
func GetListHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.GetHealthListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		healthService := provider.InitializeHealthService(svc)
		resp, err := healthService.GetList(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
