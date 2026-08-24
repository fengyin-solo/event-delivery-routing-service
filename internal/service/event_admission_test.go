package service

import (
	"testing"

	"eventbus/internal/config"
	"eventbus/internal/model"
	"eventbus/internal/store"
	"eventbus/pkg/logger"
)

func TestAdmissionRejectsMissingReferences(t *testing.T) {
	cfg := &config.Config{MaxPageSize: 100, MaxAttempts: 3}
	s := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), cfg)
	topic, _ := s.CreateTopic(model.Topic{Name: "admission-topic"})
	publisher, _ := s.CreatePublisher(model.Publisher{
		Name: "admission-publisher", TopicID: topic.ID, Endpoint: "http://publisher/callback",
	})
	cases := []model.Event{
		{TopicID: topic.ID, PublisherID: publisher.ID, Payload: "   "},
		{TopicID: "missing-topic", PublisherID: publisher.ID, Payload: "payload"},
		{TopicID: topic.ID, PublisherID: "missing-publisher", Payload: "payload"},
	}
	for i, candidate := range cases {
		if created, err := s.CreateEvent(candidate); err == nil || created != nil {
			t.Fatalf("case %d should be rejected without event, created=%+v err=%v", i, created, err)
		}
	}
	items, total, err := s.ListEvents(model.EventFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("list after rejection: %v", err)
	}
	if total != 0 || len(items) != 0 {
		t.Fatalf("rejected inputs leaked events: total=%d items=%+v", total, items)
	}
	created, err := s.CreateEvent(model.Event{TopicID: topic.ID, PublisherID: publisher.ID, Payload: " accepted "})
	if err != nil {
		t.Fatalf("valid event rejected: %v", err)
	}
	if created.Payload != "accepted" || created.Status != model.EventPending {
		t.Fatalf("valid event was not normalized: %+v", created)
	}
}
