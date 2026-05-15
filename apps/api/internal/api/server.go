package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type Server struct {
	service markerCreator
	store   markerReader
	mux     *http.ServeMux
}

type markerCreator interface {
	CreateMarker(ctx context.Context, input stream.CreateMarkerInput) (stream.Marker, stream.Job, error)
}

type markerReader interface {
	GetMarker(ctx context.Context, markerID string) (stream.Marker, error)
	GetJob(ctx context.Context, jobID string) (stream.Job, error)
}

func NewServer(service markerCreator, store markerReader) *Server {
	s := &Server{
		service: service,
		store:   store,
		mux:     http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
	s.mux.HandleFunc("POST /markers", s.handleCreateMarker)
	s.mux.HandleFunc("GET /markers/", s.handleGetMarker)
	s.mux.HandleFunc("GET /jobs/", s.handleGetJob)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCreateMarker(w http.ResponseWriter, r *http.Request) {
	var input stream.CreateMarkerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	marker, job, err := s.service.CreateMarker(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"marker_id": marker.ID,
		"job_id":    job.ID,
		"status":    string(job.Status),
	})
}

func (s *Server) handleGetMarker(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/markers/")
	marker, err := s.store.GetMarker(r.Context(), id)
	if err != nil {
		if errors.Is(err, stream.ErrNotFound) {
			writeError(w, http.StatusNotFound, "marker not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load marker")
		return
	}

	writeJSON(w, http.StatusOK, marker)
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/jobs/")
	job, err := s.store.GetJob(r.Context(), id)
	if err != nil {
		if errors.Is(err, stream.ErrNotFound) {
			writeError(w, http.StatusNotFound, "job not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load job")
		return
	}

	writeJSON(w, http.StatusOK, job)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
