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
	service              markerCreator
	store                markerReader
	runtime              *RuntimeService
	internalHeaderName   string
	internalServiceToken string
	mux                  *http.ServeMux
}

type markerCreator interface {
	CreateMarker(ctx context.Context, input stream.CreateMarkerInput) (stream.Marker, stream.Job, error)
}

type markerReader interface {
	GetMarker(ctx context.Context, markerID string) (stream.Marker, error)
	GetJob(ctx context.Context, jobID string) (stream.Job, error)
}

func NewServer(service markerCreator, store markerReader, runtime *RuntimeService, internalHeaderName, internalServiceToken string) *Server {
	s := &Server{
		service:              service,
		store:                store,
		runtime:              runtime,
		internalHeaderName:   internalHeaderName,
		internalServiceToken: internalServiceToken,
		mux:                  http.NewServeMux(),
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
	s.mux.HandleFunc("/internal/channels/", s.handleInternalChannels)
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

func (s *Server) requireInternalAuth(w http.ResponseWriter, r *http.Request) bool {
	if s.internalServiceToken == "" {
		return true
	}
	if r.Header.Get(s.internalHeaderName) != s.internalServiceToken {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}

func (s *Server) handleInternalChannels(w http.ResponseWriter, r *http.Request) {
	if s.runtime == nil {
		writeError(w, http.StatusNotFound, "runtime service unavailable")
		return
	}
	if !s.requireInternalAuth(w, r) {
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/internal/channels/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}
	channelLogin := parts[0]

	switch {
	case len(parts) == 2 && parts[1] == "runtime" && r.Method == http.MethodGet:
		state, err := s.runtime.GetRuntime(r.Context(), channelLogin)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 3 && parts[1] == "sessions" && parts[2] == "sync" && r.Method == http.MethodPost:
		state, err := s.runtime.SyncSession(r.Context(), channelLogin)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 4 && parts[1] == "segments" && parts[3] == "start" && r.Method == http.MethodPost:
		state, err := s.runtime.StartSegment(r.Context(), channelLogin, parts[2])
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 3 && parts[1] == "segments" && parts[2] == "advance" && r.Method == http.MethodPost:
		state, err := s.runtime.AdvanceSegment(r.Context(), channelLogin)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 3 && parts[1] == "title" && parts[2] == "apply-current" && r.Method == http.MethodPost:
		state, err := s.runtime.ApplyCurrentTitle(r.Context(), channelLogin)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 3 && parts[1] == "title" && parts[2] == "restore" && r.Method == http.MethodPost:
		state, err := s.runtime.RestoreTitle(r.Context(), channelLogin)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 3 && parts[1] == "title" && parts[2] == "toggle" && r.Method == http.MethodPost:
		state, err := s.runtime.ToggleTitles(r.Context(), channelLogin)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 3 && parts[1] == "title" && parts[2] == "format" && r.Method == http.MethodGet:
		format, err := s.runtime.GetTitleFormat(r.Context(), channelLogin)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"format": format})
	case len(parts) == 3 && parts[1] == "title" && parts[2] == "format" && r.Method == http.MethodPut:
		var body struct {
			Format string `json:"format"`
			Reset  bool   `json:"reset"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		var (
			state stream.RuntimeState
			err   error
		)
		if body.Reset {
			state, err = s.runtime.ResetTitleFormat(r.Context(), channelLogin)
		} else {
			state, err = s.runtime.SetTitleFormat(r.Context(), channelLogin, body.Format)
		}
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	case len(parts) == 2 && parts[1] == "markers" && r.Method == http.MethodPost:
		var body MarkerRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		state, err := s.runtime.AddMarker(r.Context(), channelLogin, body)
		if err != nil {
			s.writeRuntimeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	default:
		writeError(w, http.StatusNotFound, "route not found")
	}
}

func (s *Server) writeRuntimeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, stream.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	default:
		writeError(w, http.StatusBadRequest, err.Error())
	}
}
