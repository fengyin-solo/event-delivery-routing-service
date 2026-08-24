package service_test

import (
	"encoding/json"
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

type subscriberListEnvelope struct {
	Code int `json:"code"`
	Data struct {
		Items []model.Subscriber `json:"items"`
		Pagination struct {
			Page int `json:"page"`
			Size int `json:"size"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"data"`
}

func TestSubscriberHTTPFiltersInactiveEndpointAndKeepsFirstPage(t *testing.T) {
	cfg := &config.Config{MaxPageSize: 10, MaxAttempts: 3}
	log := logger.NewLevel(logger.LevelError)
	svc := service.New(store.NewMemoryStore(), log, cfg)
	topic, err := svc.CreateTopic(model.Topic{Name: "signals"})
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	first, err := svc.CreateSubscriber(model.Subscriber{Name: "alpha", TopicID: topic.ID, Endpoint: "https://edge.example/alerts"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err = svc.CreateSubscriber(model.Subscriber{Name: "beta", TopicID: topic.ID, Endpoint: "https://edge.example/metrics"}); err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, err = svc.ToggleSubscriberStatus(first.ID, model.EndpointInactive); err != nil {
		t.Fatalf("disable first: %v", err)
	}
	server := handler.NewServer(svc, log, cfg)
	req := httptest.NewRequest(http.MethodGet, "/api/subscribers?status=inactive&keyword=alerts&page=1&size=1", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", rec.Code, rec.Body.String())
	}
	var body subscriberListEnvelope
	if err = json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != 0 || body.Data.Pagination.Total != 1 || len(body.Data.Items) != 1 {
		t.Fatalf("filtered first page mismatch: %+v", body)
	}
	if body.Data.Items[0].ID != first.ID || body.Data.Items[0].Status != model.EndpointInactive {
		t.Fatalf("wrong subscriber returned: %+v", body.Data.Items)
	}
}
