package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type fakeService struct{}

func (fakeService) CreateMarker(_ context.Context, input stream.CreateMarkerInput) (stream.Marker, stream.Job, error) {
	return stream.Marker{ID: "marker-1"}, stream.Job{ID: "job-1", Status: stream.JobStatusPending}, nil
}

type fakeStore struct{}

func (fakeStore) GetMarker(_ context.Context, markerID string) (stream.Marker, error) {
	return stream.Marker{ID: markerID, StreamID: "stream", Label: "clip", Timestamp: "00:00:10"}, nil
}

func (fakeStore) GetJob(_ context.Context, jobID string) (stream.Job, error) {
	return stream.Job{ID: jobID, MarkerID: "marker-1", Status: stream.JobStatusCompleted}, nil
}

func TestCreateMarker(t *testing.T) {
	server := NewServer(fakeService{}, fakeStore{})
	body, _ := json.Marshal(map[string]string{
		"stream_id": "stream",
		"label":     "clip",
		"timestamp": "00:00:10",
	})

	req := httptest.NewRequest(http.MethodPost, "/markers", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestHealth(t *testing.T) {
	server := NewServer(fakeService{}, fakeStore{})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}
