package statistics

import (
	"net/http"

	"mq/common/response"
	"mq/infrastructure/provider"
	"mq/infrastructure/svc"
)

// GetNowStatisticsHandler handles requests to get current statistics
func GetNowStatisticsHandler(svc *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statisticsService := provider.InitializeStatisticsService(svc)
		resp, err := statisticsService.GetNowStatistics(r.Context())
		response.Response(r, w, resp, err)
	}
}
