package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	
	"github.com/nabiilNajm26/go-bank/internal/delivery/http"
	"github.com/nabiilNajm26/go-bank/internal/delivery/http/middleware"
	"github.com/nabiilNajm26/go-bank/internal/infrastructure/database"
	"github.com/nabiilNajm26/go-bank/internal/repository/postgres"
	"github.com/nabiilNajm26/go-bank/internal/usecase"
	"github.com/nabiilNajm26/go-bank/pkg/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database configuration
	dbConfig := &database.Config{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         5432,
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "gobank"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxLifetime:  5 * time.Minute,
	}

	// Connect to database
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(
		getEnv("JWT_ACCESS_SECRET", "access-secret-key"),
		getEnv("JWT_REFRESH_SECRET", "refresh-secret-key"),
		time.Hour,
		24*time.Hour*7,
	)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtManager)
	accountUseCase := usecase.NewAccountUseCase(accountRepo, userRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, accountRepo, db)

	// Initialize handlers
	authHandler := http.NewAuthHandler(authUseCase)
	accountHandler := http.NewAccountHandler(accountUseCase)
	transactionHandler := http.NewTransactionHandler(transactionUseCase)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))

	// Routes
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Use(middleware.AuthMiddleware(jwtManager))

	// Account routes
	accounts := protected.Group("/accounts")
	accounts.Post("/", accountHandler.CreateAccount)
	accounts.Get("/", accountHandler.GetUserAccounts)
	accounts.Get("/:id", accountHandler.GetAccount)

	// Transaction routes
	transactions := protected.Group("/transactions")
	transactions.Post("/transfer", transactionHandler.Transfer)
	transactions.Get("/", transactionHandler.GetTransactionHistory)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}