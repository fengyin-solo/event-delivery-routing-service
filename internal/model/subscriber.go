package model

import (
	"strings"
	"time"
)

// Subscriber 事件订阅者。
type Subscriber struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TopicID   string    `json:"topic_id"`
	Endpoint  string    `json:"endpoint"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验订阅者字段并规范化。
func (s *Subscriber) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.TopicID = strings.TrimSpace(s.TopicID)
	s.Endpoint = strings.TrimSpace(s.Endpoint)
	if s.Name == "" {
		return NewValidationError("name", "订阅者名称不能为空")
	}
	if s.TopicID == "" {
		return NewValidationError("topic_id", "主题 ID 不能为空")
	}
	if s.Endpoint == "" {
		return NewValidationError("endpoint", "订阅地址不能为空")
	}
	if s.Status == "" {
		s.Status = EndpointActive
	}
	if s.Status != EndpointActive && s.Status != EndpointInactive {
		return NewValidationError("status", "订阅者状态不合法")
	}
	return nil
}

// SubscriberFilter 订阅者查询筛选条件。
type SubscriberFilter struct {
	TopicID string
	Status  string
	Keyword string
}

// Match 判断订阅者是否命中筛选条件。
func (f SubscriberFilter) Match(s *Subscriber) bool {
	if f.TopicID != "" && s.TopicID != f.TopicID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" &&
			!strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Endpoint), k) {
			return false
		}
	}
	return true
}
