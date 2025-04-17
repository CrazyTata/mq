package dto

// CreateOperationRequest 创建操作记录请求
type CreateOperationRequest struct {
	Action  string `json:"action"`
	Target  string `json:"target"`
	Details string `json:"details"`
}

// GetOperationListRequest 获取操作记录列表请求
type GetOperationListRequest struct {
	Search   string `form:"search,optional"`
	Page     int64  `form:"page"`
	PageSize int64  `form:"page_size"`
}

// GetOperationResponse 获取操作记录列表响应
type GetOperationResponse struct {
	Count int64                `json:"count"`
	Data  []*OperationResponse `json:"data"`
}

// OperationResponse 操作记录响应
type OperationResponse struct {
	ID        int64  `json:"id"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	Details   string `json:"details"`
	Username  string `json:"username"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
