package service

import (
	"testing"

	"eventbus/internal/config"
	"eventbus/internal/model"
	"eventbus/internal/store"
	"eventbus/pkg/logger"
)

func TestConfiguredRetryLimitMovesOnlyFinalFailureToDeadLetter(t *testing.T) {
	if !model.CanTransitionEventStatus(model.EventFailed, model.EventDead) {
		t.Fatal("failed event must permit transition to dead at the configured limit")
	}
	t.Setenv("MAX_ATTEMPTS", "3")
	t.Setenv("MAX_PAGE_SIZE", "25")
	cfg := config.Load()
	if cfg.MaxAttempts != 3 {
		t.Fatalf("loaded max attempts=%d, want 3", cfg.MaxAttempts)
	}
	s := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), cfg)
	topic, err := s.CreateTopic(model.Topic{Name: "delivery"})
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	publisher, err := s.CreatePublisher(model.Publisher{Name: "relay", TopicID: topic.ID, Endpoint: "http://relay"})
	if err != nil {
		t.Fatalf("create publisher: %v", err)
	}
	event, err := s.CreateEvent(model.Event{TopicID: topic.ID, PublisherID: publisher.ID, Payload: "payload"})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	if _, err = s.FailEvent(event.ID, "attempt one"); err != nil {
		t.Fatalf("attempt one: %v", err)
	}
	second, err := s.RetryEvent(event.ID, false, "attempt two")
	if err != nil || second.Status != model.EventFailed || second.Attempts != 2 {
		t.Fatalf("second failure should remain retryable: event=%+v err=%v", second, err)
	}
	final, err := s.RetryEvent(event.ID, false, "attempt three")
	if err != nil || final.Status != model.EventDead || final.Attempts != 3 {
		t.Fatalf("third failure should become dead: event=%+v err=%v", final, err)
	}
	letters, total, err := s.ListDeadLetters(model.DeadLetterFilter{TopicID: topic.ID}, 1, 10)
	if err != nil || total != 1 || len(letters) != 1 || letters[0].Reason != "attempt three" {
		t.Fatalf("final failure receipt mismatch: total=%d letters=%+v err=%v", total, letters, err)
	}
}
