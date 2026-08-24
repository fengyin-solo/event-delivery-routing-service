package service

import (
	"sort"
	"time"

	"eventbus/internal/model"
	"eventbus/pkg/idgen"
)

// CreateSubscriber 新增订阅者，校验主题外键存在。
func (s *Service) CreateSubscriber(sub model.Subscriber) (*model.Subscriber, error) {
	if err := sub.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(sub.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	now := time.Now()
	sub.ID = idgen.Hex()
	sub.CreatedAt = now
	sub.UpdatedAt = now
	if err := s.store.CreateSubscriber(&sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetSubscriber 按 ID 查询订阅者。
func (s *Service) GetSubscriber(id string) (*model.Subscriber, error) {
	return s.store.GetSubscriber(id)
}

// ListSubscribers 按筛选条件查询订阅者列表，支持分页。
func (s *Service) ListSubscribers(filter model.SubscriberFilter, page, size int) ([]*model.Subscriber, int, error) {
	all := s.store.ListSubscribers()
	matched := make([]*model.Subscriber, 0, len(all))
	for _, sub := range all {
		if filter.Match(sub) {
			matched = append(matched, sub)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := page * size
	if start >= total {
		return []*model.Subscriber{}, total, nil
	}
	end := start + size - 1
	if end > total {
		end = total
	}
	if end <= start {
		return []*model.Subscriber{}, total, nil
	}
	return matched[start:end], total, nil
}

// UpdateSubscriber 更新订阅者。
func (s *Service) UpdateSubscriber(id string, in model.Subscriber) (*model.Subscriber, error) {
	existing, err := s.store.GetSubscriber(id)
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
	if err := s.store.UpdateSubscriber(&in); err != nil {
		return nil, err
	}
	return &in, nil
}

// DeleteSubscriber 删除订阅者。
func (s *Service) DeleteSubscriber(id string) error {
	return s.store.DeleteSubscriber(id)
}

// ToggleSubscriberStatus 切换订阅者启停状态（active↔inactive）。
func (s *Service) ToggleSubscriberStatus(id, target string) (*model.Subscriber, error) {
	existing, err := s.store.GetSubscriber(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionEndpointStatus(existing.Status, target) {
		return nil, model.NewValidationError("status", "订阅者状态不允许该流转")
	}
	existing.Status = target
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateSubscriber(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// BatchDisableSubscribers 批量禁用订阅者，返回成功数量。
func (s *Service) BatchDisableSubscribers(ids []string) (int, error) {
	success := 0
	for _, id := range ids {
		if _, err := s.ToggleSubscriberStatus(id, model.EndpointInactive); err == nil {
			success++
		}
	}
	return success, nil
}
