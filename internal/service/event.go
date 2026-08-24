package service

import (
	"sort"
	"time"

	"eventbus/internal/model"
	"eventbus/pkg/idgen"
)

// CreateEvent 新增事件，校验主题与发布者外键存在，初始为待投递。
func (s *Service) CreateEvent(e model.Event) (*model.Event, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(e.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	if _, err := s.store.GetPublisher(e.PublisherID); err != nil {
		return nil, model.NewValidationError("publisher_id", "发布者不存在")
	}
	e.ID = idgen.Hex()
	e.CreatedAt = time.Now()
	e.Status = model.EventPending
	e.Attempts = 0
	if err := s.store.CreateEvent(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

// GetEvent 按 ID 查询事件。
func (s *Service) GetEvent(id string) (*model.Event, error) {
	return s.store.GetEvent(id)
}

// ListEvents 按筛选条件查询事件列表，支持分页。
func (s *Service) ListEvents(filter model.EventFilter, page, size int) ([]*model.Event, int, error) {
	all := s.store.ListEvents()
	matched := make([]*model.Event, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Event{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteEvent 删除事件。
func (s *Service) DeleteEvent(id string) error {
	return s.store.DeleteEvent(id)
}

// DeliverEvent 模拟投递成功（pending/failed→delivered），写入投递时间。
func (s *Service) DeliverEvent(id string) (*model.Event, error) {
	existing, err := s.store.GetEvent(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionEventStatus(existing.Status, model.EventDelivered) {
		return nil, model.NewValidationError("status", "当前事件状态不允许标记投递成功")
	}
	now := time.Now()
	existing.Status = model.EventDelivered
	existing.DeliveredAt = &now
	existing.LastError = ""
	if err := s.store.UpdateEvent(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// FailEvent 模拟投递失败（pending→failed），记录失败原因。
func (s *Service) FailEvent(id, reason string) (*model.Event, error) {
	existing, err := s.store.GetEvent(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionEventStatus(existing.Status, model.EventFailed) {
		return nil, model.NewValidationError("status", "当前事件状态不允许标记投递失败")
	}
	existing.Status = model.EventFailed
	existing.Attempts = 1
	existing.LastError = reason
	if err := s.store.UpdateEvent(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// RetryEvent 重试失败事件：attempts+1；成功则投递，超限则进入死信。
func (s *Service) RetryEvent(id string, success bool, reason string) (*model.Event, error) {
	existing, err := s.store.GetEvent(id)
	if err != nil {
		return nil, err
	}
	if existing.Status != model.EventFailed {
		return nil, model.NewValidationError("status", "仅失败状态的事件可重试")
	}
	if success && !model.CanTransitionEventStatus(existing.Status, model.EventDelivered) {
		return nil, model.NewValidationError("status", "失败事件不能转为投递成功")
	}
	existing.Attempts++

	if success {
		now := time.Now()
		existing.Status = model.EventDelivered
		existing.DeliveredAt = &now
		existing.LastError = reason
		if err := s.store.UpdateEvent(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	existing.LastError = reason
	if existing.Attempts >= s.maxAttempts() {
		// 达到最大重试次数，进入死信。
		existing.Status = model.EventDead
		if err := s.store.UpdateEvent(existing); err != nil {
			return nil, err
		}
		if err := s.createDeadLetter(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	if err := s.store.UpdateEvent(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// createDeadLetter 为进入死信的事件生成死信记录。
func (s *Service) createDeadLetter(e *model.Event) error {
	now := time.Now()
	d := &model.DeadLetter{
		ID:        idgen.Hex(),
		EventID:   e.ID,
		TopicID:   e.TopicID,
		Reason:    e.LastError,
		FailedAt:  now,
		CreatedAt: now,
	}
	return s.store.CreateDeadLetter(d)
}
