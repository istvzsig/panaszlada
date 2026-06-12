package models

import "time"

type ReportStatus string
type ReportCategory string
type ReportSeverity string

// =====================
// STATUS
// =====================
const (
	StatusOpen       ReportStatus = "open"
	StatusInProgress ReportStatus = "in_progress"
	StatusResolved   ReportStatus = "resolved"
	StatusInvalid    ReportStatus = "invalid"
)

// =====================
// CATEGORY
// =====================
const (
	CategoryHomeless ReportCategory = "homeless"
	CategorySafety   ReportCategory = "safety"
	CategoryInfra    ReportCategory = "infrastructure"
	CategoryOther    ReportCategory = "other"
)

// =====================
// SEVERITY
// =====================
const (
	SeverityLow    ReportSeverity = "low"
	SeverityMedium ReportSeverity = "medium"
	SeverityHigh   ReportSeverity = "high"
)

// =====================
// REPORT
// =====================
type Report struct {
	ID           string `json:"id"`
	TrackingCode string `json:"tracking_code"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Category ReportCategory `json:"category"`
	Severity ReportSeverity `json:"severity"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Status ReportStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
