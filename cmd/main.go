package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/config"
	"github.com/anshiq/bookmyshow-go/internal/handler"
	"github.com/anshiq/bookmyshow-go/internal/middleware"
	"github.com/anshiq/bookmyshow-go/internal/payment"
	"github.com/anshiq/bookmyshow-go/internal/repository"
	"github.com/anshiq/bookmyshow-go/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB (primary)
	mongoClient, err := config.ConnectMongo(cfg.MongoDB_URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Initialize MongoDB DB1 (for config storage)
	mongoClientDB1, err := config.ConnectMongoDB1(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to MongoDB DB1: %v", err)
	} else {
		defer mongoClientDB1.Disconnect(context.Background())
	}

	// Initialize MongoDB DB2
	mongoClientDB2, err := config.ConnectMongoDB2(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to MongoDB DB2: %v", err)
	} else {
		defer mongoClientDB2.Disconnect(context.Background())
	}

	// Initialize Redis
	redisClient := config.ConnectRedis(cfg.RedisAddr)
	defer redisClient.Close()

	// Initialize Elasticsearch
	esClient := config.ConnectElasticsearch(cfg.ElasticsearchURL)

	// Initialize repositories
	movieRepo := repository.NewMovieRepository(mongoClient, cfg.DatabaseName)
	theatreRepo := repository.NewTheatreRepository(mongoClient, cfg.DatabaseName)
	screenRepo := repository.NewScreenRepository(mongoClient, cfg.DatabaseName)
	showRepo := repository.NewShowRepository(mongoClient, cfg.DatabaseName)
	ticketRepo := repository.NewTicketRepository(mongoClient, cfg.DatabaseName)
	userRepo := repository.NewUserRepository(mongoClient, cfg.DatabaseName)
	seatRepo := repository.NewSeatRepository(mongoClient, cfg.DatabaseName)
	movieElasticRepo := repository.NewMovieElasticRepository(esClient)
	configRepo := repository.NewConfigRepository(mongoClientDB1, cfg.DatabaseName_DB1)

	// Initialize payment service
	paymentFactory := payment.NewPaymentProviderFactory()
	paymentService := paymentFactory.GetPaymentService()

	// Initialize services
	movieService := service.NewMovieService(movieRepo, movieElasticRepo, redisClient)
	theatreService := service.NewTheatreService(theatreRepo)
	screenService := service.NewScreenService(screenRepo)
	showService := service.NewShowService(showRepo, movieElasticRepo, theatreRepo, screenRepo, seatRepo)
	ticketService := service.NewTicketService(ticketRepo, seatRepo, showRepo, theatreRepo, screenRepo, movieElasticRepo, redisClient, paymentService)
	userService := service.NewUserService(userRepo)
	configService := service.NewConfigService(configRepo)

	// Initialize handlers
	movieHandler := handler.NewMovieHandler(movieService)
	theatreHandler := handler.NewTheatreHandler(theatreService)
	screenHandler := handler.NewScreenHandler(screenService)
	showHandler := handler.NewShowHandler(showService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	userHandler := handler.NewUserHandler(userService)
	configHandler := handler.NewConfigHandler(configService)

	// Setup router
	router := gin.Default()
	router.Use(middleware.UserFilter())
	router.Use(middleware.CORS())

	// Setup routes
	setupRoutes(router, movieHandler, theatreHandler, screenHandler, showHandler, ticketHandler, userHandler, configHandler)

	// Admin routes
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AdminAuth())
	{
		adminRoutes.POST("/movie", movieHandler.AddMovie)
		adminRoutes.POST("/theatre", theatreHandler.AddTheatre)
		adminRoutes.POST("/screen", screenHandler.AddScreen)
		adminRoutes.POST("/show", showHandler.AddShow)
	}

	// User routes
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/signup", userHandler.Signup)
		userRoutes.POST("/login", userHandler.Login)
		userRoutes.GET("/movies", movieHandler.GetAllMovies)
		userRoutes.GET("/movie", movieHandler.SearchMovies)
		userRoutes.GET("/shows", showHandler.GetMovieShows)
		userRoutes.POST("/book-tickets", ticketHandler.BookTicket)
		userRoutes.GET("/tickets", ticketHandler.GetAllTickets)
		userRoutes.DELETE("/cancel-tickets", ticketHandler.CancelTicket)
		userRoutes.POST("/db-config", configHandler.SetConfig)
	}

	// Config routes
	configRoutes := router.Group("/config")
	{
		configRoutes.GET("/db-type", configHandler.GetByDBType)
		configRoutes.GET("/get-by-type-and-hashId", configHandler.GetByDBTypeAndHashID)
	}

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRoutes(router *gin.Engine,
	movieHandler *handler.MovieHandler,
	theatreHandler *handler.TheatreHandler,
	screenHandler *handler.ScreenHandler,
	showHandler *handler.ShowHandler,
	ticketHandler *handler.TicketHandler,
	userHandler *handler.UserHandler,
	configHandler *handler.ConfigHandler) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
