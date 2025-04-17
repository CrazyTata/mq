package dto

// AppointmentRequest 预约请求
type AppointmentRequest struct {
	Id          int64  `json:"id,optional"`
	PatientID   int64  `json:"patient_id"`
	PatientName string `json:"patient_name"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Duration    string `json:"duration"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

// CreateAppointmentResponse 创建预约响应
type CreateAppointmentResponse struct {
	ID int64 `json:"id"`
}

// AppointmentResponse 预约响应
type AppointmentResponse struct {
	ID          int64  `json:"id"`
	PatientID   int64  `json:"patient_id"`
	PatientName string `json:"patient_name"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Duration    string `json:"duration"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// AppointmentListResponse 预约列表响应
type AppointmentListResponse struct {
	List     []*AppointmentResponse `json:"list"`
	Total    int64                  `json:"total"`
	Page     int64                  `json:"page"`
	PageSize int64                  `json:"page_size"`
}

// DeleteAppointmentRequest 删除预约请求
type DeleteAppointmentRequest struct {
	ID int64 `json:"id"`
}

// GetAppointmentByIDRequest 获取预约请求
type GetAppointmentByIDRequest struct {
	ID int64 `form:"id"`
}

// GetAppointmentListRequest 获取预约列表请求
type GetAppointmentListRequest struct {
	Page       int64  `form:"page"`
	PageSize   int64  `form:"page_size"`
	PatientId  int64  `form:"patient_id,optional"`
	Search     string `form:"search,optional"`
	SearchType int64  `form:"search_type,optional"`
	Order      int64  `form:"order,optional"` // 1: 降序, 0: 升序
}
