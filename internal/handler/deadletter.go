package handler

import (
	"net/http"

	"eventbus/internal/model"
	"eventbus/pkg/httpx"
)

func (s *Server) registerDeadLetterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/dead-letters", s.listDeadLetters)
	mux.HandleFunc("GET /api/dead-letters/{id}", s.getDeadLetter)
	mux.HandleFunc("DELETE /api/dead-letters/{id}", s.deleteDeadLetter)
}

func (s *Server) listDeadLetters(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DeadLetterFilter{TopicID: r.URL.Query().Get("topic_id")}
	items, total, err := s.svc.ListDeadLetters(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDeadLetter(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.GetDeadLetter(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDeadLetter(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteDeadLetter(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
