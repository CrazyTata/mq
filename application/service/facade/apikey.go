package facade

import (
	"context"
	"sync"

	"mq/application/service"
	"mq/common/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

// ApiKeyFacade API密钥外观
type ApiKeyFacade struct {
	keyManager     *auth.ApiKeyManager
	subjectService *service.SubjectService
	mu             sync.RWMutex
}

// NewApiKeyFacade 创建新的API密钥外观
func NewApiKeyFacade(subjectService *service.SubjectService) *ApiKeyFacade {
	facade := &ApiKeyFacade{
		keyManager:     auth.NewApiKeyManager(),
		subjectService: subjectService,
	}

	// 服务启动时加载一次配置
	if err := facade.RefreshCache(context.Background()); err != nil {
		logx.Errorf("初始化API密钥缓存失败: %v", err)
	}

	return facade
}

// Init 初始化外观，注册事件监听
func (f *ApiKeyFacade) Init() {
	// 注册事件监听
	f.subjectService.OnSubjectUpdated(f.HandleSubjectUpdate)
}

// HandleSubjectUpdate 处理主题更新事件
func (f *ApiKeyFacade) HandleSubjectUpdate(ctx context.Context) error {
	return f.RefreshCache(ctx)
}

// GetKeyManager 获取密钥管理器
func (f *ApiKeyFacade) GetKeyManager() *auth.ApiKeyManager {
	return f.keyManager
}

// RefreshCache 刷新缓存
func (f *ApiKeyFacade) RefreshCache(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 从应用服务层获取所有激活的API密钥
	subjects, err := f.subjectService.GetAllActiveSubjects(ctx)
	if err != nil {
		return err
	}

	// 构建新的缓存
	keys := make(map[string]*auth.ApiKeyStruct)
	for _, s := range subjects {
		key := &auth.ApiKeyStruct{
			AppId:        s.AppId,
			AppSecret:    s.AppSecret,
			MerchantID:   s.MerchantId,
			MerchantName: s.MerchantName,
			IsActive:     s.IsActive == 1,
			ExpireAt:     s.ExpireAt.Time,
			AllowedIPs:   s.AllowedIps.String,
			AllowedPaths: s.AllowedPaths.String,
		}
		keys[key.AppId] = key
	}

	// 更新缓存
	f.keyManager.RefreshCache(keys)
	return nil
}
