package service

import "eventbus/internal/model"

// StatsOverview 事件总线概览统计。
type StatsOverview struct {
	TotalTopics      int            `json:"total_topics"`
	TotalPublishers  int            `json:"total_publishers"`
	TotalSubscribers int            `json:"total_subscribers"`
	TotalEvents      int            `json:"total_events"`
	EventsByStatus   map[string]int `json:"events_by_status"`
	TotalDeadLetters int            `json:"total_dead_letters"`
}

// Overview 计算整体概览统计。
func (s *Service) Overview() *StatsOverview {
	o := &StatsOverview{
		EventsByStatus: make(map[string]int),
	}

	o.TotalTopics = len(s.store.ListTopics())
	o.TotalPublishers = len(s.store.ListPublishers())
	o.TotalSubscribers = len(s.store.ListSubscribers())
	o.TotalDeadLetters = len(s.store.ListDeadLetters())

	for _, e := range s.store.ListEvents() {
		o.TotalEvents++
		o.EventsByStatus[e.Status]++
	}
	// 确保四种状态都有键，便于前端展示。
	for _, st := range []string{model.EventPending, model.EventDelivered, model.EventFailed, model.EventDead} {
		if _, ok := o.EventsByStatus[st]; !ok {
			o.EventsByStatus[st] = 0
		}
	}

	return o
}
