package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
    ID             int       `json:"id" db:"id"`
    CorrelationID  uuid.UUID `json:"correlationId" db:"correlation_id"`
    Amount         float64   `json:"amount" db:"amount"`
    ProcessedBy    string    `json:"processedBy" db:"processed_by"` // "default" or "fallback"
    RequestedAt    time.Time `json:"requestedAt" db:"requested_at"`
    CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}

type PaymentRequest struct {
    CorrelationID uuid.UUID `json:"correlationId"`
    Amount        float64   `json:"amount"`
}

type PaymentsSummaryResponse struct {
    Default  ProcessorSummary `json:"default"`
    Fallback ProcessorSummary `json:"fallback"`
}

type ProcessorSummary struct {
    TotalRequests int     `json:"totalRequests"`
    TotalAmount   float64 `json:"totalAmount"`
}