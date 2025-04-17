package patient

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetPatientListHandler handles requests to get patient list
func GetPatientListHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.GetPatientListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		patientService := provider.InitializePatientService(svc)
		resp, err := patientService.GetPatientList(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
