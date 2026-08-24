package store

import "eventbus/internal/model"

// CreateSubscriber 新增订阅者。
func (s *MemoryStore) CreateSubscriber(sub *model.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers[sub.ID] = sub
	return nil
}

// GetSubscriber 按 ID 查询订阅者。
func (s *MemoryStore) GetSubscriber(id string) (*model.Subscriber, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub, ok := s.subscribers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sub, nil
}

// ListSubscribers 返回全部订阅者。
func (s *MemoryStore) ListSubscribers() []*model.Subscriber {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Subscriber, 0, len(s.subscribers))
	for _, sub := range s.subscribers {
		list = append(list, sub)
	}
	return list
}

// UpdateSubscriber 更新订阅者。
func (s *MemoryStore) UpdateSubscriber(sub *model.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscribers[sub.ID]; !ok {
		return ErrNotFound
	}
	s.subscribers[sub.ID] = sub
	return nil
}

// DeleteSubscriber 删除订阅者。
func (s *MemoryStore) DeleteSubscriber(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscribers[id]; !ok {
		return ErrNotFound
	}
	delete(s.subscribers, id)
	return nil
}
