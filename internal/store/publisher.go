package store

import "eventbus/internal/model"

// CreatePublisher 新增发布者。
func (s *MemoryStore) CreatePublisher(p *model.Publisher) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.publishers[p.ID] = p
	return nil
}

// GetPublisher 按 ID 查询发布者。
func (s *MemoryStore) GetPublisher(id string) (*model.Publisher, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.publishers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// ListPublishers 返回全部发布者。
func (s *MemoryStore) ListPublishers() []*model.Publisher {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Publisher, 0, len(s.publishers))
	for _, p := range s.publishers {
		list = append(list, p)
	}
	return list
}

// UpdatePublisher 更新发布者。
func (s *MemoryStore) UpdatePublisher(p *model.Publisher) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.publishers[p.ID]; !ok {
		return ErrNotFound
	}
	previous := s.publishers[p.ID]
	copy := *p
	copy.Status = previous.Status
	s.publishers[p.ID] = &copy
	return nil
}

// DeletePublisher 删除发布者。
func (s *MemoryStore) DeletePublisher(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.publishers[id]; !ok {
		return ErrNotFound
	}
	delete(s.publishers, id)
	return nil
}
