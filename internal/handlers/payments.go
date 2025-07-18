package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
    paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
    return &PaymentHandler{
        paymentService: paymentService,
    }
}

func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
    var req models.PaymentRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // Validate UUID
    if req.CorrelationID == uuid.Nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid correlationId"})
    }

    // Validate amount
    if req.Amount <= 0 {
        return c.Status(400).JSON(fiber.Map{"error": "Amount must be positive"})
    }

    err := h.paymentService.ProcessPayment(c.Context(), req.CorrelationID, req.Amount)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to process payment"})
    }

    return c.Status(200).JSON(fiber.Map{"message": "Payment processed successfully"})
}

func (h *PaymentHandler) GetPaymentsSummary(c *fiber.Ctx) error {
    // Parse query parameters
    fromStr := c.Query("from")
    toStr := c.Query("to")

    var from, to *time.Time

    if fromStr != "" {
        if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
            from = &parsed
        }
    }

    if toStr != "" {
        if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
            to = &parsed
        }
    }

    summary, err := h.paymentService.GetPaymentsSummary(c.Context(), from, to)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get payments summary"})
    }

    return c.JSON(summary)
}