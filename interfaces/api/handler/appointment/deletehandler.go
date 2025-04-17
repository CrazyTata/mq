package appointment

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// DeleteHandler handles requests to delete an appointment
func DeleteHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.DeleteAppointmentRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		appointmentService := provider.InitializeAppointmentService(svc)
		err := appointmentService.Delete(r.Context(), req.ID)
		response.Response(r, w, nil, err)
	}
}
