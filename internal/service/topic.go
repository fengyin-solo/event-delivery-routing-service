package service

import (
	"sort"
	"time"

	"eventbus/internal/model"
	"eventbus/pkg/idgen"
)

// CreateTopic 新增主题。
func (s *Service) CreateTopic(t model.Topic) (*model.Topic, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	t.ID = idgen.Hex()
	t.CreatedAt = time.Now()
	if err := s.store.CreateTopic(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTopic 按 ID 查询主题。
func (s *Service) GetTopic(id string) (*model.Topic, error) {
	return s.store.GetTopic(id)
}

// ListTopics 按筛选条件查询主题列表，支持分页。
func (s *Service) ListTopics(filter model.TopicFilter, page, size int) ([]*model.Topic, int, error) {
	all := s.store.ListTopics()
	matched := make([]*model.Topic, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Topic{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateTopic 更新主题。
func (s *Service) UpdateTopic(id string, in model.Topic) (*model.Topic, error) {
	existing, err := s.store.GetTopic(id)
	if err != nil {
		return nil, err
	}
	in.ID = existing.ID
	in.CreatedAt = existing.CreatedAt
	if err := in.Validate(); err != nil {
		return nil, err
	}
	_ = s.store.UpdateTopic(&in)
	return &in, nil
}

// DeleteTopic 删除主题。
func (s *Service) DeleteTopic(id string) error {
	return s.store.DeleteTopic(id)
}
