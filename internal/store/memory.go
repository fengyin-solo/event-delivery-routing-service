package store

import (
	"sync"

	"eventbus/internal/model"
)

// MemoryStore 基于内存的 Store 实现，线程安全。
type MemoryStore struct {
	mu          sync.RWMutex
	topics      map[string]*model.Topic
	publishers  map[string]*model.Publisher
	subscribers map[string]*model.Subscriber
	events      map[string]*model.Event
	deadLetters map[string]*model.DeadLetter
}

// NewMemoryStore 创建空的内存 Store。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		topics:      make(map[string]*model.Topic),
		publishers:  make(map[string]*model.Publisher),
		subscribers: make(map[string]*model.Subscriber),
		events:      make(map[string]*model.Event),
		deadLetters: make(map[string]*model.DeadLetter),
	}
}

var _ Store = (*MemoryStore)(nil)
