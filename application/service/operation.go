package service

import (
	"context"
	"mq/application/assembler"
	"mq/application/dto"
	"mq/common/util"
	"mq/domain/operation"
)

// OperationService 操作记录应用服务
type OperationService struct {
	repo operation.OperationRepository
}

// NewHealthAppService 创建健康记录应用服务
func NewOperationService(repo operation.OperationRepository) *OperationService {
	return &OperationService{
		repo: repo,
	}
}

// Create 创建操作记录
func (s *OperationService) Create(ctx context.Context, req *dto.CreateOperationRequest) (*dto.Response, error) {
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	username, err := util.GetUsernameFromContext(ctx)
	if err != nil {
		return nil, err
	}
	operation := operation.Create(req.Action, req.Target, req.Details, username, userID)
	_, err = s.repo.Create(ctx, operation)
	if err != nil {
		return nil, err
	}

	return assembler.Return(err)
}

// GetList 获取健康记录列表
func (s *OperationService) GetList(ctx context.Context, req *dto.GetOperationListRequest) (*dto.GetOperationResponse, error) {
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	operations, total, err := s.repo.GetList(ctx, req.Search, userID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOOperationList(operations, total), nil
}

// CountRecentByUserID 根据用户ID获取操作记录数量
func (s *OperationService) CountRecentByUserID(ctx context.Context, userID string, days int) (int64, error) {
	return s.repo.CountRecentByUserID(ctx, userID, days)
}
