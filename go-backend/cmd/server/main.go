package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"user-api/config"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/middleware"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	zapLogger := logger.NewLogger()
	defer zapLogger.Sync()

	// Initialize database connection
	db, err := config.NewDBConnection()
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	zapLogger.Info("Database connection established")

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, zapLogger)
	userHandler := handler.NewUserHandler(userService, zapLogger)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.LoggerMiddleware(zapLogger))

	// Setup routes
	routes.SetupRoutes(app, userHandler)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	zapLogger.Info("Starting server on port " + port)

	// Start server
	if err := app.Listen(":" + port); err != nil {
		zapLogger.Fatal("Failed to start server", err)
	}
}
