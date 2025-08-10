package main

import (
	"log"
	"os"

	"library-management-backend/internal/database"
	"library-management-backend/internal/handlers"
	"library-management-backend/internal/middleware"
	"library-management-backend/internal/services"
	"library-management-backend/pkg/config"

	_ "library-management-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Library Management API
// @version 1.0
// @description A simple library management system API with URL processing service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	if os.Getenv("GIN_MODE") == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Initialize validator
	validate := validator.New()

	// Initialize services
	bookService := services.NewBookService(db.DB, logger)
	urlService := services.NewURLService(logger)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookService, validate, logger)
	urlHandler := handlers.NewURLHandler(urlService, validate, logger)

	// Initialize Gin router
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler(logger))
	router.Use(gin.Recovery())

	// API routes
	api := router.Group("/api")
	{
		// Book routes
		books := api.Group("/books")
		{
			books.GET("", bookHandler.GetBooks)
			books.POST("", bookHandler.CreateBook)
			books.GET("/:id", bookHandler.GetBook)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
		}

		// URL processing route
		api.POST("/url-process", urlHandler.ProcessURL)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	logger.WithField("port", cfg.Server.Port).Info("Starting server")
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
