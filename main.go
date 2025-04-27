package main

import (
	"fmt"
	"log"
	"strconv"

	"packify/internal/config"
	"packify/internal/handlers"
	"packify/internal/models"
	"packify/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting Packify API...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup database
	if err := models.SetupDatabase(db); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Initialize services
	packService := services.NewPackService(db)

	// Initialize template renderer
	renderer, err := handlers.NewTemplateRenderer()
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}

	// Initialize handlers
	handler := handlers.NewHandler(packService, renderer)

	// Create Echo instance
	e := echo.New()
	e.Renderer = renderer

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Register routes
	handler.RegisterRoutes(e)

	// Add a simple health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Start server
	serverAddr := ":" + strconv.Itoa(cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
