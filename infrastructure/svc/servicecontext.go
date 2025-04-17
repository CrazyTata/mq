package svc

import (
	"context"
	"mq/infrastructure/config"
	"mq/interfaces/middleware"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	config *config.Config
	db     sqlx.SqlConn
	cache  cache.ClusterConf
	redis  *redis.Redis
	Login  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	logger := logx.WithContext(context.Background())
	redisConf := redis.RedisConf{
		Host: c.Redis.Host,
		Type: c.Redis.Type,
		Pass: c.Redis.Pass,
		Tls:  c.Redis.Tls,
	}
	cacheConf := cache.CacheConf{
		{
			RedisConf: redisConf,
			Weight:    100,
		},
	}
	conn := sqlx.NewMysql(c.DB.DataSource)

	redisClient, err := redis.NewRedis(redisConf)
	if err != nil {
		logger.Errorf("Failed to create Redis client: %v", err)
	}

	return &ServiceContext{
		config: &c,
		db:     conn,
		cache:  cacheConf,
		redis:  redisClient,
		Login:  middleware.NewLoginMiddleware(c.FrontendAuth.AccessSecret).Handle,
	}
}

func (s *ServiceContext) GetConfig() *config.Config {
	return s.config
}

func (s *ServiceContext) GetDB() sqlx.SqlConn {
	return s.db
}

func (s *ServiceContext) GetCache() cache.ClusterConf {
	return s.cache
}

func (s *ServiceContext) GetRedis() *redis.Redis {
	return s.redis
}
