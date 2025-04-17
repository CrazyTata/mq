package appointment

import (
	"net/http"

	"mq/application/dto"
	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetByIDHandler handles requests to get appointment by ID
func GetByIDHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.GetAppointmentByIDRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		appointmentService := provider.InitializeAppointmentService(svc)
		resp, err := appointmentService.GetByID(r.Context(), req.ID)
		response.Response(r, w, resp, err)
	}
}
