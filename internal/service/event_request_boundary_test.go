package service_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"eventbus/internal/config"
	"eventbus/internal/handler"
	"eventbus/internal/model"
	"eventbus/internal/service"
	"eventbus/internal/store"
	"eventbus/pkg/logger"
)

func TestRequestBoundaryLeavesNoRejectedEvents(t *testing.T) {
	cfg := &config.Config{MaxPageSize: 10, MaxAttempts: 3}
	log := logger.NewLevel(logger.LevelError)
	svc := service.New(store.NewMemoryStore(), log, cfg)
	topic, err := svc.CreateTopic(model.Topic{Name: "ingress"})
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	publisher, err := svc.CreatePublisher(model.Publisher{Name: "gateway", TopicID: topic.ID, Endpoint: "https://gateway.example/push"})
	if err != nil {
		t.Fatalf("create publisher: %v", err)
	}
	server := handler.NewServer(svc, log, cfg)
	validPrefix := fmt.Sprintf(`{"topic_id":%q,"publisher_id":%q,"payload":""}`, topic.ID, publisher.ID)
	cases := []string{
		validPrefix,
		fmt.Sprintf(`{"topic_id":%q,"publisher_id":%q,"payload":"one"}{"payload":"two"}`, topic.ID, publisher.ID),
	}
	for i, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		server.Routes().ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("case %d status=%d body=%s", i, rec.Code, rec.Body.String())
		}
	}
	validBody := fmt.Sprintf(`{"topic_id":%q,"publisher_id":%q,"payload":"plain text"}`, topic.ID, publisher.ID)
	validReq := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewBufferString(validBody))
	validReq.Header.Set("Content-Type", "application/json")
	validRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(validRec, validReq)
	if validRec.Code != http.StatusCreated {
		t.Fatalf("plain text event status=%d body=%s", validRec.Code, validRec.Body.String())
	}
	items, total, err := svc.ListEvents(model.EventFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Payload != "plain text" {
		t.Fatalf("rejected requests changed valid events: total=%d items=%+v", total, items)
	}
}
