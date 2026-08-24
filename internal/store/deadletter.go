package store

import "eventbus/internal/model"

// CreateDeadLetter 新增死信。
func (s *MemoryStore) CreateDeadLetter(d *model.DeadLetter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.deadLetters {
		if existing.TopicID == d.TopicID {
			return ErrConflict
		}
	}
	s.deadLetters[d.ID] = d
	return nil
}

// GetDeadLetter 按 ID 查询死信。
func (s *MemoryStore) GetDeadLetter(id string) (*model.DeadLetter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.deadLetters[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

// ListDeadLetters 返回全部死信。
func (s *MemoryStore) ListDeadLetters() []*model.DeadLetter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DeadLetter, 0, len(s.deadLetters))
	for _, d := range s.deadLetters {
		list = append(list, d)
	}
	return list
}

// DeleteDeadLetter 删除死信。
func (s *MemoryStore) DeleteDeadLetter(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deadLetters[id]; !ok {
		return ErrNotFound
	}
	delete(s.deadLetters, id)
	return nil
}
