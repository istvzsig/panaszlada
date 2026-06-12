package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/istvzsig/panaszlada/storage"
)

func TestCreateReport_Success(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewReportHandler(store, "test-key")

	body := `{
		"title": "test",
		"description": "test desc",
		"category": "safety",
		"latitude": 47.5,
		"longitude": 19.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateReport(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
}

func TestCreateReport_InvalidJSON(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewReportHandler(store, "test-key")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports", strings.NewReader("{bad json"))
	w := httptest.NewRecorder()

	h.CreateReport(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateReport_OutOfBounds(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewReportHandler(store, "test-key")

	body := `{
		"title": "test",
		"description": "test",
		"category": "safety",
		"latitude": 10,
		"longitude": 10
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateReport(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetReport_NotFound(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewReportHandler(store, "test-key")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/KH-123", nil)
	req = mux.SetURLVars(req, map[string]string{
		"tracking_code": "KH-123",
	})

	w := httptest.NewRecorder()

	h.GetReport(w, req)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewReportHandler(store, "test-key")

	body := `{"status":"invalid_status"}`

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/reports/1", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	h.UpdateStatus(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
