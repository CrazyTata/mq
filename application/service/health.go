package service

import (
	"context"
	"encoding/json"
	"mq/application/assembler"
	"mq/application/dto"
	"mq/common/util"
	"mq/domain/health"
	"mq/infrastructure/queue"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// HealthService 健康记录应用服务
type HealthService struct {
	repo  health.HealthRepository
	queue queue.QueueManager
}

// NewHealthAppService 创建健康记录应用服务
func NewHealthService(repo health.HealthRepository, queue queue.QueueManager) *HealthService {
	return &HealthService{
		repo:  repo,
		queue: queue,
	}
}

// Create 创建健康记录
func (s *HealthService) Save(ctx context.Context, req *dto.HealthRequest) (*dto.CreateHealthResponse, error) {
	logger := logx.WithContext(ctx)
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}
	var healthID int64
	var healthRecord *health.HealthRecords
	if req.Id == 0 {
		healthRecord = health.Create(req.PatientID, req.PatientName, date, req.RecordType, req.Diagnosis, req.Treatment, req.Notes, req.VitalSigns, req.Medications, req.Attachments, userID)
		healthID, err = s.repo.Create(ctx, healthRecord)
		if err != nil {
			return nil, err
		}
	} else {
		healthRecord, err = s.repo.GetByID(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		healthRecord.Update(req.PatientName, date, req.RecordType, req.Diagnosis, req.Treatment, req.Notes, req.VitalSigns, req.Medications, req.Attachments, req.PatientID)
		err = s.repo.UpdateHealth(ctx, healthRecord)
		if err != nil {
			return nil, err
		}
		healthID = healthRecord.Id
	}

	// 发送消息到队列
	messageData, err := json.Marshal(healthRecord)
	if err != nil {
		return nil, err
	}

	message := &queue.Message{
		ID:   util.NewSnowflake().String(),
		Type: queue.HealthRecordSaved,
		Body: messageData,
	}

	if err := s.queue.Publish(ctx, message); err != nil {
		// 记录错误但不影响主流程
		logger.Errorf("Failed to publish health record message: %v", err)
	}

	return &dto.CreateHealthResponse{
		ID: healthID,
	}, nil
}

// GetByID 根据ID获取健康记录
func (s *HealthService) GetByID(ctx context.Context, id int64) (*dto.HealthResponse, error) {
	health, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOHealth(health), nil
}

// Delete 删除健康记录
func (s *HealthService) Delete(ctx context.Context, id int64) error {
	logger := logx.WithContext(ctx)
	healthRecord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	healthRecord.MarkAsDeleted()

	// 发送删除消息到队列
	messageData, err := json.Marshal(healthRecord)
	if err != nil {
		return err
	}

	message := &queue.Message{
		ID:   util.NewSnowflake().String(),
		Type: queue.HealthRecordSaved,
		Body: messageData,
	}

	if err := s.queue.Publish(ctx, message); err != nil {
		// 记录错误但不影响主流程
		logger.Errorf("Failed to publish health record delete message: %v", err)
	}

	return s.repo.UpdateHealth(ctx, healthRecord)
}

// GetList 获取健康记录列表
func (s *HealthService) GetList(ctx context.Context, req *dto.GetHealthListRequest) (*dto.HealthListResponse, error) {
	userID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	healths, total, err := s.repo.GetList(ctx, userID, req.Search, req.Status, req.Order, req.PatientId, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return assembler.DOTODTOHealthList(healths, total), nil
}

// CountByUserID 根据用户ID获取健康记录数量
func (s *HealthService) CountByUserID(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountByUserID(ctx, userID)
}
