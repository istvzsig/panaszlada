package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"github.com/istvzsig/panaszlada/models"
	"github.com/istvzsig/panaszlada/storage"
)

type ReportHandler struct {
	store     storage.ReportStore
	apiKey    string
	sanitizer *bluemonday.Policy
}

// =====================
// REQUESTS
// =====================

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
// VALIDATION MAPS
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
	return &ReportHandler{
		store:     store,
		apiKey:    apiKey,
		sanitizer: bluemonday.UGCPolicy(),
	}
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

	// sanitize input (XSS protection)
	title := strings.TrimSpace(h.sanitizer.Sanitize(req.Title))
	desc := strings.TrimSpace(h.sanitizer.Sanitize(req.Description))

	if title == "" || desc == "" {
		writeJSON(w, 400, map[string]string{"error": "missing fields"})
		return
	}

	// default category
	if req.Category == "" {
		req.Category = "other"
	}

	cat, ok := validCategories[req.Category]
	if !ok {
		writeJSON(w, 400, map[string]string{"error": "invalid category"})
		return
	}

	// bounds check
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
		Title:        title,
		Description:  desc,
		Category:     cat,
		Severity:     models.SeverityLow,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Status:       models.StatusOpen,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := h.store.Create(r.Context(), report); err != nil {
		writeJSON(w, 500, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, 201, map[string]any{
		"id":            report.ID,
		"tracking_code": report.TrackingCode,
		"status":        report.Status,
	})
}

// =====================
// GET REPORT
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
		writeJSON(w, 500, map[string]string{"error": "internal server error"})
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
		writeJSON(w, 500, map[string]string{"error": "internal server error"})
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
		writeJSON(w, 500, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, 200, map[string]string{"deleted": id})
}

// =====================
// TRACKING CODE (FIXED)
// =====================

func generateTrackingCode() string {
	id := uuid.New().ID()
	return fmt.Sprintf("KH-%06X", id&0xffffff)
}

// =====================
// JSON HELPER
// =====================

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
