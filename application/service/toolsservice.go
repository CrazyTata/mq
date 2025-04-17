package service

import (
	"context"

	"mq/application/dto"
	"mq/infrastructure/integration/qiniu"
)

// ToolService 工具应用服务
type ToolService struct {
	repo qiniu.UploadInterface
}

// NewToolService 创建工具应用服务
func NewToolService(repo qiniu.UploadInterface) *ToolService {
	return &ToolService{
		repo: repo,
	}
}

// Search 根据ID获取课程
func (s *ToolService) UploadToken(ctx context.Context) (*dto.UploadTokenResponse, error) {

	token, err := s.repo.GetUploadToken(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.UploadTokenResponse{
		Token: token,
	}, nil
}
