package handler

import (
	"net/http"

	"eventbus/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /healthz", s.healthz)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.Overview())
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]string{"status": "ok"})
}
