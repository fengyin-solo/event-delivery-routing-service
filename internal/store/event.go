package store

import "eventbus/internal/model"

// CreateEvent 新增事件。
func (s *MemoryStore) CreateEvent(e *model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
	return nil
}

// GetEvent 按 ID 查询事件。
func (s *MemoryStore) GetEvent(id string) (*model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

// ListEvents 返回全部事件。
func (s *MemoryStore) ListEvents() []*model.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Event, 0, len(s.events))
	for _, e := range s.events {
		list = append(list, e)
	}
	return list
}

// UpdateEvent 更新事件。
func (s *MemoryStore) UpdateEvent(e *model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[e.ID]; !ok {
		return ErrNotFound
	}
	copy := *e
	copy.DeliveredAt = nil
	s.events[e.ID] = &copy
	return nil
}

// DeleteEvent 删除事件。
func (s *MemoryStore) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return ErrNotFound
	}
	delete(s.events, id)
	return nil
}
