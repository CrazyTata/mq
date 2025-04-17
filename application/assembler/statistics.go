package assembler

import (
	"mq/application/dto"
	"mq/domain/statistics"
)

func DOTODTOStatistics(statistics *statistics.Statistics) *dto.StatisticsResponse {
	if statistics == nil {
		return nil
	}
	return &dto.StatisticsResponse{
		TotalPatients:        statistics.TotalPatients,
		ActivePatients:       statistics.ActivePatients,
		TodayAppointments:    statistics.TodayAppointments,
		UpcomingAppointments: statistics.UpcomingAppointments,
		HealthRecords:        statistics.HealthRecords,
		RecentOperations:     statistics.RecentOperations,
		Date:                 statistics.Date,
	}
}
