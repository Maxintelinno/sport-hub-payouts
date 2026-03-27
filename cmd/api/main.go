package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sport-hub/sport-hub-payouts/config"
	"github.com/sport-hub/sport-hub-payouts/database"
	"github.com/sport-hub/sport-hub-payouts/internal/handler"
	customMiddleware "github.com/sport-hub/sport-hub-payouts/internal/middleware"
	"github.com/sport-hub/sport-hub-payouts/internal/repository"
	"github.com/sport-hub/sport-hub-payouts/internal/service"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Database connection
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Printf("Warning: Could not connect to database: %v. Running in mock mode if needed.", err)
	}

	// Initialize Echo
	e := echo.New()

	// Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize Layers
	var walletRepo repository.SettlementRepository
	if db != nil {
		walletRepo = repository.NewSettlementRepository(db)
	} else {
		// Fallback or Mock (not implemented here for brevity)
		log.Println("Note: Repository initialized without database connection.")
	}

	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Sport Hub Payouts API")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "UP",
		})
	})

	// V1 Routes
	v1 := e.Group("/v1")

	// Owner Routes
	owner := v1.Group("/owner", customMiddleware.JWTMiddleware(cfg.JwtSecret, cfg.SkipJwtVerify))
	owner.GET("/wallet/summary", walletHandler.GetSummary)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
