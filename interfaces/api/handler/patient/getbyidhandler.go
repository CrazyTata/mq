package patient

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetByIDHandler handles requests to get patient by ID
func GetByIDHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.GetIdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		patientService := provider.InitializePatientService(svc)
		resp, err := patientService.GetPatientByID(r.Context(), req.ID)
		response.Response(r, w, resp, err)
	}
}
