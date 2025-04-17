package service

import (
	"context"
	"sync"

	"mq/domain/subject"

	"github.com/zeromicro/go-zero/core/logx"
)

// SubjectUpdateHandler 主题更新处理器
type SubjectUpdateHandler func(ctx context.Context) error

// SubjectService 主题服务
type SubjectService struct {
	subjectRepo subject.SubjectRepository
	handlers    []SubjectUpdateHandler
	mu          sync.RWMutex
}

// NewSubjectService 创建新的主题服务
func NewSubjectService(repo subject.SubjectRepository) *SubjectService {
	return &SubjectService{
		subjectRepo: repo,
		handlers:    make([]SubjectUpdateHandler, 0),
	}
}

// OnSubjectUpdated 注册主题更新处理器
func (s *SubjectService) OnSubjectUpdated(handler SubjectUpdateHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, handler)
}

// notifyHandlers 通知所有处理器
func (s *SubjectService) notifyHandlers(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, handler := range s.handlers {
		if err := handler(ctx); err != nil {
			logx.Errorf("处理主题更新事件失败: %v", err)
		}
	}
}

// GetAllActiveSubjects 获取所有激活的API密钥
func (s *SubjectService) GetAllActiveSubjects(ctx context.Context) ([]*subject.Subject, error) {
	// 从数据库获取所有激活的API密钥
	subjects, err := s.subjectRepo.FindAllActive(ctx)
	if err != nil {
		logx.Errorf("获取激活的API密钥失败: %v", err)
		return nil, err
	}
	return subjects, nil
}

// GetSubjectByAppId 根据AppId获取API密钥信息
func (s *SubjectService) GetSubjectByAppId(ctx context.Context, appId string) (*subject.Subject, error) {
	subject, err := s.subjectRepo.FindByAppId(ctx, appId)
	if err != nil {
		logx.Errorf("获取API密钥信息失败: %v", err)
		return nil, err
	}
	return subject, nil
}

// UpdateSubject 更新主题信息
func (s *SubjectService) UpdateSubject(ctx context.Context, subject *subject.Subject) error {
	// 更新主题信息
	if err := s.subjectRepo.UpdateSubject(ctx, subject); err != nil {
		return err
	}

	// 通知所有处理器
	s.notifyHandlers(ctx)

	return nil
}
