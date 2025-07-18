package main

import (
	"log"
	"rinha-backend-2025/internal/config"
	"rinha-backend-2025/internal/database"
	"rinha-backend-2025/internal/handlers"
	"rinha-backend-2025/internal/repository"
	"rinha-backend-2025/internal/services"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    // Load config
    cfg := config.Load()

    // Connect to database
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Run migrations
    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Initialize repositories
    paymentRepo := repository.NewPaymentRepository(db)

    // Initialize services
    processorClient := services.NewProcessorClient(cfg.ProcessorDefaultURL, cfg.ProcessorFallbackURL)
    paymentService := services.NewPaymentService(paymentRepo, processorClient)

    // Initialize handlers
    paymentHandler := handlers.NewPaymentHandler(paymentService)

    // Configure Fiber with Sonic JSON
    app := fiber.New(fiber.Config{
        JSONEncoder: sonic.Marshal,
        JSONDecoder: sonic.Unmarshal,
        Prefork:     false,
    })

    // Middlewares
    app.Use(logger.New())
    app.Use(recover.New())

    // Routes
    app.Post("/payments", paymentHandler.ProcessPayment)
    app.Get("/payments-summary", paymentHandler.GetPaymentsSummary)
    app.Get("/health", handlers.HealthCheck)

    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(app.Listen(":" + cfg.Port))
}