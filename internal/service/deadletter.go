package service

import (
	"sort"

	"eventbus/internal/model"
)

// GetDeadLetter 按 ID 查询死信。
func (s *Service) GetDeadLetter(id string) (*model.DeadLetter, error) {
	return s.store.GetDeadLetter(id)
}

// ListDeadLetters 按筛选条件查询死信列表，支持分页。
func (s *Service) ListDeadLetters(filter model.DeadLetterFilter, page, size int) ([]*model.DeadLetter, int, error) {
	all := s.store.ListDeadLetters()
	matched := make([]*model.DeadLetter, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].FailedAt.After(matched[j].FailedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DeadLetter{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteDeadLetter 删除死信。
func (s *Service) DeleteDeadLetter(id string) error {
	return s.store.DeleteDeadLetter(id)
}
