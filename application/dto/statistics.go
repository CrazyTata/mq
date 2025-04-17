package dto

type StatisticsResponse struct {
	TotalPatients        int64  `json:"total_patients"`
	ActivePatients       int64  `json:"active_patients"`
	TodayAppointments    int64  `json:"today_appointments"`
	UpcomingAppointments int64  `json:"upcoming_appointments"`
	HealthRecords        int64  `json:"health_records"`
	RecentOperations     int64  `json:"recent_operations"`
	Date                 string `json:"date"`
}
