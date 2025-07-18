package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
    repo            *repository.PaymentRepository
    processorClient *ProcessorClient
}

func NewPaymentService(repo *repository.PaymentRepository, processorClient *ProcessorClient) *PaymentService {
    return &PaymentService{
        repo:            repo,
        processorClient: processorClient,
    }
}

func (s *PaymentService) ProcessPayment(ctx context.Context, correlationID uuid.UUID, amount float64) error {
    requestedAt := time.Now().UTC()
    
    // Try default processor first
    processor := "default"
    err := s.processorClient.ProcessPayment(ctx, correlationID, amount, requestedAt, true)
    
    if err != nil {
        // Fallback to secondary processor
        processor = "fallback"
        err = s.processorClient.ProcessPayment(ctx, correlationID, amount, requestedAt, false)
        if err != nil {
            return err
        }
    }

    // Store in database
    return s.repo.CreatePayment(ctx, correlationID, amount, processor, requestedAt)
}

func (s *PaymentService) GetPaymentsSummary(ctx context.Context, from, to *time.Time) (*models.PaymentsSummaryResponse, error) {
    return s.repo.GetPaymentsSummary(ctx, from, to)
}