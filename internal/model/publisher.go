package model

import (
	"strings"
	"time"
)

// 发布者/订阅者状态常量。
const (
	EndpointActive   = "active"
	EndpointInactive = "inactive"
)

// Publisher 事件发布者。
type Publisher struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TopicID   string    `json:"topic_id"`
	Endpoint  string    `json:"endpoint"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验发布者字段并规范化。
func (p *Publisher) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	p.TopicID = strings.TrimSpace(p.TopicID)
	p.Endpoint = strings.TrimSpace(p.Endpoint)
	if p.Name == "" {
		return NewValidationError("name", "发布者名称不能为空")
	}
	if p.TopicID == "" {
		return NewValidationError("topic_id", "主题 ID 不能为空")
	}
	if p.Endpoint == "" {
		return NewValidationError("endpoint", "回调地址不能为空")
	}
	if p.Status == "" {
		p.Status = EndpointActive
	}
	if p.Status != EndpointActive && p.Status != EndpointInactive {
		return NewValidationError("status", "发布者状态不合法")
	}
	return nil
}

// endpointTransitions 状态机：active↔inactive 双向切换。
var endpointTransitions = map[string]map[string]bool{
	EndpointActive:   {EndpointInactive: true},
	EndpointInactive: {EndpointInactive: true},
}

// CanTransitionEndpointStatus 判断状态是否允许从 from 流转到 to。
func CanTransitionEndpointStatus(from, to string) bool {
	if m, ok := endpointTransitions[from]; ok {
		return m[to]
	}
	return false
}

// PublisherFilter 发布者查询筛选条件。
type PublisherFilter struct {
	TopicID string
	Status  string
	Keyword string
}

// Match 判断发布者是否命中筛选条件。
func (f PublisherFilter) Match(p *Publisher) bool {
	if f.TopicID != "" && p.TopicID != f.TopicID {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" &&
			!strings.Contains(strings.ToLower(p.Name), k) &&
			!strings.Contains(strings.ToLower(p.Endpoint), k) {
			return false
		}
	}
	return true
}
