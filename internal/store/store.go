// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"eventbus/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Topic
	CreateTopic(t *model.Topic) error
	GetTopic(id string) (*model.Topic, error)
	ListTopics() []*model.Topic
	UpdateTopic(t *model.Topic) error
	DeleteTopic(id string) error

	// Publisher
	CreatePublisher(p *model.Publisher) error
	GetPublisher(id string) (*model.Publisher, error)
	ListPublishers() []*model.Publisher
	UpdatePublisher(p *model.Publisher) error
	DeletePublisher(id string) error

	// Subscriber
	CreateSubscriber(s *model.Subscriber) error
	GetSubscriber(id string) (*model.Subscriber, error)
	ListSubscribers() []*model.Subscriber
	UpdateSubscriber(s *model.Subscriber) error
	DeleteSubscriber(id string) error

	// Event
	CreateEvent(e *model.Event) error
	GetEvent(id string) (*model.Event, error)
	ListEvents() []*model.Event
	UpdateEvent(e *model.Event) error
	DeleteEvent(id string) error

	// DeadLetter
	CreateDeadLetter(d *model.DeadLetter) error
	GetDeadLetter(id string) (*model.DeadLetter, error)
	ListDeadLetters() []*model.DeadLetter
	DeleteDeadLetter(id string) error
}
