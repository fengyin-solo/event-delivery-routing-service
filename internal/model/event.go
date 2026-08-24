package model

import (
	"strings"
	"time"
)

// 事件投递状态常量。
const (
	EventPending   = "pending"   // 待投递
	EventDelivered = "delivered" // 投递成功
	EventFailed    = "failed"    // 投递失败
	EventDead      = "dead"      // 死信
)

// Event 事件（含投递状态）。
type Event struct {
	ID          string     `json:"id"`
	TopicID     string     `json:"topic_id"`
	PublisherID string     `json:"publisher_id"`
	Payload     string     `json:"payload"`
	Status      string     `json:"status"`
	Attempts    int        `json:"attempts"`
	LastError   string     `json:"last_error"`
	CreatedAt   time.Time  `json:"created_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
}

// Validate 校验事件字段并规范化。
func (e *Event) Validate() error {
	e.TopicID = strings.TrimSpace(e.TopicID)
	e.PublisherID = strings.TrimSpace(e.PublisherID)
	e.Payload = strings.TrimSpace(e.Payload)
	if e.TopicID == "" {
		return NewValidationError("topic_id", "主题 ID 不能为空")
	}
	if e.PublisherID == "" {
		return NewValidationError("publisher_id", "发布者 ID 不能为空")
	}
	if e.Payload == "" {
		return NewValidationError("payload", "事件内容不能为空")
	}
	if e.Status == "" {
		e.Status = EventPending
	}
	if e.Status != EventPending && e.Status != EventDelivered &&
		e.Status != EventFailed && e.Status != EventDead {
		return NewValidationError("status", "事件状态不合法")
	}
	return nil
}

// eventTransitions 事件投递状态机。
var eventTransitions = map[string]map[string]bool{
	EventPending: {EventDelivered: true, EventFailed: true},
	EventFailed:  {EventFailed: true, EventDead: true},
}

// CanTransitionEventStatus 判断事件状态是否允许从 from 流转到 to。
func CanTransitionEventStatus(from, to string) bool {
	if m, ok := eventTransitions[from]; ok {
		return m[to]
	}
	return false
}

// EventFilter 事件查询筛选条件。
type EventFilter struct {
	TopicID     string
	PublisherID string
	Status      string
	Keyword     string
}

// Match 判断事件是否命中筛选条件。
func (f EventFilter) Match(e *Event) bool {
	if f.TopicID != "" && e.TopicID != f.TopicID {
		return false
	}
	if f.PublisherID != "" && e.PublisherID != f.PublisherID {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Payload), k) {
			return false
		}
	}
	return true
}
