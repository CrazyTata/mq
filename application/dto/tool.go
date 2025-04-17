package dto

type UploadTokenResponse struct {
	Token string `json:"token"`
}

type AddNotificationRequest struct {
	Ids []int64 `json:"ids"`
}
