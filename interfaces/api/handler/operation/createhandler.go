package operation

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// CreateHandler handles requests to create a new operation record
func CreateHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateOperationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		operationService := provider.InitializeOperationService(svc)
		resp, err := operationService.Create(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
