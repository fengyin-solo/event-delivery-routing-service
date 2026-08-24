// Package service 实现业务逻辑层。
package service

import (
	"eventbus/internal/config"
	"eventbus/internal/store"
	"eventbus/pkg/logger"
)

// Service 业务服务，聚合各实体业务方法。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建业务服务。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}

// maxAttempts 返回事件投递最大重试次数。
func (s *Service) maxAttempts() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 1
}
