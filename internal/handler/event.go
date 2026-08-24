package handler

import (
	"net/http"

	"eventbus/internal/model"
	"eventbus/pkg/httpx"
)

func (s *Server) registerEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/events", s.createEvent)
	mux.HandleFunc("GET /api/events", s.listEvents)
	mux.HandleFunc("GET /api/events/{id}", s.getEvent)
	mux.HandleFunc("DELETE /api/events/{id}", s.deleteEvent)
	mux.HandleFunc("POST /api/events/{id}/deliver", s.deliverEvent)
	mux.HandleFunc("POST /api/events/{id}/fail", s.failEvent)
	mux.HandleFunc("POST /api/events/{id}/retry", s.retryEvent)
}

type createEventRequest struct {
	TopicID     string `json:"topic_id"`
	PublisherID string `json:"publisher_id"`
	Payload     string `json:"payload"`
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateEvent(model.Event{
		TopicID:     req.TopicID,
		PublisherID: req.PublisherID,
		Payload:     req.Payload,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EventFilter{
		TopicID:     r.URL.Query().Get("topic_id"),
		PublisherID: r.URL.Query().Get("publisher_id"),
		Status:      r.URL.Query().Get("status"),
		Keyword:     r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	e, err := s.svc.GetEvent(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteEvent(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) deliverEvent(w http.ResponseWriter, r *http.Request) {
	e, err := s.svc.DeliverEvent(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type failEventRequest struct {
	Reason string `json:"reason"`
}

func (s *Server) failEvent(w http.ResponseWriter, r *http.Request) {
	var req failEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.FailEvent(r.PathValue("id"), req.Reason)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type retryEventRequest struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func (s *Server) retryEvent(w http.ResponseWriter, r *http.Request) {
	var req retryEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.RetryEvent(r.PathValue("id"), req.Success, req.Reason)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}
