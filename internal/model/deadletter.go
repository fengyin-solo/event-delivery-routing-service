package model

import (
	"strings"
	"time"
)

// DeadLetter 死信记录，保存投递超限失败的事件。
type DeadLetter struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	TopicID   string    `json:"topic_id"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验死信字段并规范化。
func (d *DeadLetter) Validate() error {
	d.EventID = strings.TrimSpace(d.EventID)
	d.TopicID = strings.TrimSpace(d.TopicID)
	d.Reason = strings.TrimSpace(d.Reason)
	if d.EventID == "" {
		return NewValidationError("event_id", "事件 ID 不能为空")
	}
	if d.TopicID == "" {
		return NewValidationError("topic_id", "主题 ID 不能为空")
	}
	return nil
}

// DeadLetterFilter 死信查询筛选条件。
type DeadLetterFilter struct {
	TopicID string
}

// Match 判断死信是否命中筛选条件。
func (f DeadLetterFilter) Match(d *DeadLetter) bool {
	if f.TopicID != "" && d.TopicID != f.TopicID {
		return false
	}
	return true
}
