package patient

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// DeleteHandler handles requests to delete patient
func DeleteHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.PostIdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		patientService := provider.InitializePatientService(svc)
		err := patientService.DeletePatient(r.Context(), req.ID)
		response.Response(r, w, nil, err)
	}
}
