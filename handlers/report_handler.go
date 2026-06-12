package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/istvzsig/panaszlada/models"
	"github.com/istvzsig/panaszlada/storage"
)

type ReportHandler struct {
	store  storage.ReportStore
	apiKey string
}

type CreateReportRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// =====================
// VALIDATION
// =====================
var validCategories = map[string]models.ReportCategory{
	"homeless":       models.CategoryHomeless,
	"safety":         models.CategorySafety,
	"infrastructure": models.CategoryInfra,
	"other":          models.CategoryOther,
}

var validStatuses = map[string]models.ReportStatus{
	"open":        models.StatusOpen,
	"in_progress": models.StatusInProgress,
	"resolved":    models.StatusResolved,
	"invalid":     models.StatusInvalid,
}

// =====================
// CONSTRUCTOR
// =====================
func NewReportHandler(store storage.ReportStore, apiKey string) *ReportHandler {
	return &ReportHandler{store: store, apiKey: apiKey}
}

// =====================
// CREATE REPORT
// =====================
func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	var req CreateReportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}

	if req.Title == "" || req.Description == "" {
		writeJSON(w, 400, map[string]string{"error": "missing fields"})
		return
	}

	cat, ok := validCategories[req.Category]
	if !ok {
		writeJSON(w, 400, map[string]string{"error": "invalid category"})
		return
	}

	if req.Latitude < 45 || req.Latitude > 49 {
		writeJSON(w, 400, map[string]string{"error": "out of bounds (lat)"})
		return
	}

	if req.Longitude < 16 || req.Longitude > 23 {
		writeJSON(w, 400, map[string]string{"error": "out of bounds (lng)"})
		return
	}

	report := models.Report{
		ID:           uuid.New().String(),
		TrackingCode: generateTrackingCode(),
		Title:        req.Title,
		Description:  req.Description,
		Category:     cat,
		Severity:     models.SeverityLow,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Status:       models.StatusOpen,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := h.store.Create(r.Context(), report); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 201, map[string]any{
		"id":            report.ID,
		"tracking_code": report.TrackingCode,
		"status":        report.Status,
	})
}

// =====================
// GET
// =====================
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	tc := mux.Vars(r)["tracking_code"]

	report, err := h.store.GetByCode(r.Context(), tc)
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "not found"})
		return
	}

	writeJSON(w, 200, report)
}

// =====================
// LIST
// =====================
func (h *ReportHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.store.List(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, reports)
}

// =====================
// UPDATE STATUS
// =====================
func (h *ReportHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}

	status, ok := validStatuses[req.Status]
	if !ok {
		writeJSON(w, 400, map[string]string{"error": "invalid status"})
		return
	}

	if err := h.store.UpdateStatus(r.Context(), id, string(status)); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]string{"ok": "true"})
}

// =====================
// DELETE
// =====================
func (h *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.store.Delete(r.Context(), id); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]string{"deleted": id})
}

// =====================
// TRACKING CODE
// =====================
func generateTrackingCode() string {
	id := uuid.New().String()
	return "KH-" + strings.ToUpper(id[:6])
}

// =====================
// JSON helper
// =====================
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
