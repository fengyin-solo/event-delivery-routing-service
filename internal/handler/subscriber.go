package handler

import (
	"net/http"

	"eventbus/internal/model"
	"eventbus/pkg/httpx"
)

func (s *Server) registerSubscriberRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/subscribers", s.createSubscriber)
	mux.HandleFunc("GET /api/subscribers", s.listSubscribers)
	mux.HandleFunc("GET /api/subscribers/{id}", s.getSubscriber)
	mux.HandleFunc("PUT /api/subscribers/{id}", s.updateSubscriber)
	mux.HandleFunc("DELETE /api/subscribers/{id}", s.deleteSubscriber)
	mux.HandleFunc("POST /api/subscribers/{id}/activate", s.activateSubscriber)
	mux.HandleFunc("POST /api/subscribers/{id}/disable", s.disableSubscriber)
	mux.HandleFunc("POST /api/subscribers/batch-disable", s.batchDisableSubscribers)
}

type subscriberRequest struct {
	Name     string `json:"name"`
	TopicID  string `json:"topic_id"`
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
}

func (s *Server) createSubscriber(w http.ResponseWriter, r *http.Request) {
	var req subscriberRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.CreateSubscriber(model.Subscriber{
		Name:     req.Name,
		TopicID:  req.TopicID,
		Endpoint: req.Endpoint,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sub)
}

func (s *Server) listSubscribers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SubscriberFilter{
		TopicID: r.URL.Query().Get("topic_id"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSubscribers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSubscriber(w http.ResponseWriter, r *http.Request) {
	sub, err := s.svc.GetSubscriber(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) updateSubscriber(w http.ResponseWriter, r *http.Request) {
	var req subscriberRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.UpdateSubscriber(r.PathValue("id"), model.Subscriber{
		Name:     req.Name,
		TopicID:  req.TopicID,
		Endpoint: req.Endpoint,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) deleteSubscriber(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSubscriber(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) activateSubscriber(w http.ResponseWriter, r *http.Request) {
	sub, err := s.svc.ToggleSubscriberStatus(r.PathValue("id"), model.EndpointActive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) disableSubscriber(w http.ResponseWriter, r *http.Request) {
	sub, err := s.svc.ToggleSubscriberStatus(r.PathValue("id"), model.EndpointInactive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

type batchDisableSubscribersRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDisableSubscribers(w http.ResponseWriter, r *http.Request) {
	var req batchDisableSubscribersRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchDisableSubscribers(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"disabled": count})
}
