package handler

import (
	"net/http"

	"eventbus/internal/model"
	"eventbus/pkg/httpx"
)

func (s *Server) registerPublisherRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/publishers", s.createPublisher)
	mux.HandleFunc("GET /api/publishers", s.listPublishers)
	mux.HandleFunc("GET /api/publishers/{id}", s.getPublisher)
	mux.HandleFunc("PUT /api/publishers/{id}", s.updatePublisher)
	mux.HandleFunc("DELETE /api/publishers/{id}", s.deletePublisher)
	mux.HandleFunc("POST /api/publishers/{id}/activate", s.activatePublisher)
	mux.HandleFunc("POST /api/publishers/{id}/disable", s.disablePublisher)
}

type publisherRequest struct {
	Name     string `json:"name"`
	TopicID  string `json:"topic_id"`
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
}

func (s *Server) createPublisher(w http.ResponseWriter, r *http.Request) {
	var req publisherRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreatePublisher(model.Publisher{
		Name:     req.Name,
		TopicID:  req.TopicID,
		Endpoint: req.Endpoint,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listPublishers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.PublisherFilter{
		TopicID: r.URL.Query().Get("topic_id"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListPublishers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getPublisher(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetPublisher(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updatePublisher(w http.ResponseWriter, r *http.Request) {
	var req publisherRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdatePublisher(r.PathValue("id"), model.Publisher{
		Name:     req.Name,
		TopicID:  req.TopicID,
		Endpoint: req.Endpoint,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deletePublisher(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeletePublisher(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) activatePublisher(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.TogglePublisherStatus(r.PathValue("id"), model.EndpointActive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) disablePublisher(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.TogglePublisherStatus(r.PathValue("id"), model.EndpointInactive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}
