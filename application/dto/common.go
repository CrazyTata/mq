package dto

type Response struct {
	Message string `json:"message"`
}

type Page struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

type GetIdRequest struct {
	ID int64 `form:"id"`
}

type PostIdRequest struct {
	ID int64 `json:"id"`
}
