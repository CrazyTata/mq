package dto

type CreatePatientResponse struct {
	ID int64 `json:"id"`
}

type Patient struct {
	Id          int64  `json:"id,optional"`
	FriendlyId  string `json:"friendly_id,optional"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	Attachments string `json:"attachments"`
	Details     string `json:"details,optional"`
	Status      string `json:"status"`
	History     string `json:"history"`
	Allergies   string `json:"allergies"`
	Note        string `json:"note"`
	LastVisit   string `json:"last_visit,optional"`
	CreatedAt   string `json:"created_at,optional"`
	UpdatedAt   string `json:"updated_at,optional"`
}

type PatientListResponse struct {
	Total int64      `json:"total"`
	List  []*Patient `json:"list"`
}

type GetPatientByFriendlyIDRequest struct {
	FriendlyID string `form:"friendly_id"`
}

type GetPatientListRequest struct {
	Page     int64  `form:"page"`
	PageSize int64  `form:"page_size"`
	Search   string `form:"search,optional"`
	Status   string `form:"status,optional"`
	Order    int64  `form:"order,optional"` // 1: 降序, 0: 升序
}
