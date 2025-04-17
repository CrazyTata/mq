package qiniu

import "context"

type UploadInterface interface {
	GetUploadToken(ctx context.Context) (string, error)
}
