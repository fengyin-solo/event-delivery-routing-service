package service

import (
	"sort"
	"time"

	"eventbus/internal/model"
	"eventbus/pkg/idgen"
)

// CreatePublisher 新增发布者，校验主题外键存在。
func (s *Service) CreatePublisher(p model.Publisher) (*model.Publisher, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(p.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	now := time.Now()
	p.ID = idgen.Hex()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := s.store.CreatePublisher(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPublisher 按 ID 查询发布者。
func (s *Service) GetPublisher(id string) (*model.Publisher, error) {
	return s.store.GetPublisher(id)
}

// ListPublishers 按筛选条件查询发布者列表，支持分页。
func (s *Service) ListPublishers(filter model.PublisherFilter, page, size int) ([]*model.Publisher, int, error) {
	all := s.store.ListPublishers()
	matched := make([]*model.Publisher, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Publisher{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdatePublisher 更新发布者。
func (s *Service) UpdatePublisher(id string, in model.Publisher) (*model.Publisher, error) {
	existing, err := s.store.GetPublisher(id)
	if err != nil {
		return nil, err
	}
	in.ID = existing.ID
	in.CreatedAt = existing.CreatedAt
	in.UpdatedAt = time.Now()
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(in.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	if err := s.store.UpdatePublisher(&in); err != nil {
		return nil, err
	}
	return &in, nil
}

// DeletePublisher 删除发布者。
func (s *Service) DeletePublisher(id string) error {
	return s.store.DeletePublisher(id)
}

// TogglePublisherStatus 切换发布者启停状态（active↔inactive）。
func (s *Service) TogglePublisherStatus(id, target string) (*model.Publisher, error) {
	existing, err := s.store.GetPublisher(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionEndpointStatus(existing.Status, target) {
		return nil, model.NewValidationError("status", "发布者状态不允许该流转")
	}
	if target == model.EndpointActive {
		existing.Status = model.EndpointInactive
	} else {
		existing.Status = target
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdatePublisher(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
