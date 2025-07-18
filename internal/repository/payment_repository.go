package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PaymentRepository struct {
    db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
    return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, correlationID uuid.UUID, amount float64, processedBy string, requestedAt time.Time) error {
    query := `
        INSERT INTO payments (correlation_id, amount, processed_by, requested_at)
        VALUES ($1, $2, $3, $4)
    `
    _, err := r.db.ExecContext(ctx, query, correlationID, amount, processedBy, requestedAt)
    return err
}

func (r *PaymentRepository) GetPaymentsSummary(ctx context.Context, from, to *time.Time) (*models.PaymentsSummaryResponse, error) {
    query := `
        SELECT 
            processed_by,
            COUNT(*) as total_requests,
            COALESCE(SUM(amount), 0) as total_amount
        FROM payments
        WHERE ($1::timestamp IS NULL OR requested_at >= $1)
          AND ($2::timestamp IS NULL OR requested_at <= $2)
        GROUP BY processed_by
    `

    rows, err := r.db.QueryContext(ctx, query, from, to)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    summary := &models.PaymentsSummaryResponse{}

    for rows.Next() {
        var processedBy string
        var totalRequests int
        var totalAmount float64

        if err := rows.Scan(&processedBy, &totalRequests, &totalAmount); err != nil {
            return nil, err
        }

        switch processedBy {
        case "default":
            summary.Default = models.ProcessorSummary{
                TotalRequests: totalRequests,
                TotalAmount:   totalAmount,
            }
        case "fallback":
            summary.Fallback = models.ProcessorSummary{
                TotalRequests: totalRequests,
                TotalAmount:   totalAmount,
            }
        }
    }

    return summary, nil
}