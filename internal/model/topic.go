package model

import (
	"strings"
	"time"
)

// Topic 事件主题。
type Topic struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate 校验主题字段并规范化。
func (t *Topic) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.Description = strings.TrimSpace(t.Description)
	if t.Name == "" {
		return NewValidationError("name", "主题名称不能为空")
	}
	return nil
}

// TopicFilter 主题查询筛选条件。
type TopicFilter struct {
	Keyword string
}

// Match 判断主题是否命中筛选条件。
func (f TopicFilter) Match(t *Topic) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" &&
			!strings.Contains(strings.ToLower(t.Name), k) &&
			!strings.Contains(strings.ToLower(t.Description), k) {
			return false
		}
	}
	return true
}
