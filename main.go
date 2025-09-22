package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"booking-system/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create the booking application
	app := handlers.NewBookingApp()
	
	// Create Gin router
	router := gin.Default()
	
	// Middleware for logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS middleware for frontend integration
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})
	
	// API Routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", app.HealthCheck)
		
		// Conferences
		api.GET("/conferences", app.GetConferences)
		
		// Users
		api.POST("/users", app.CreateUser)
		api.GET("/users/:userID/bookings", app.GetUserBookings)
		api.GET("/users/:userID/reservations", app.GetUserReservations)
		
		// Bookings (direct booking - old way)
		api.POST("/bookings", app.CreateBooking)
		api.GET("/bookings", app.GetAllBookings)  // Get all bookings for testing
		api.GET("/bookings/:id", app.GetBooking)
		
		// Reservations (new payment queue system)
		api.POST("/reservations", app.CreateReservation)
		api.GET("/reservations/:id", app.GetReservation)
		api.POST("/reservations/:id/confirm", app.ConfirmReservation)
		api.DELETE("/reservations/:id", app.CancelReservation)

		// Wait queue
		api.POST("/queue/enqueue", app.EnqueueWait)
		api.GET("/queue/:conferenceID/position", app.GetQueuePosition)
		api.POST("/queue/claim", app.ClaimNext)
	}
	
	// Serve static files and frontend
	router.Static("/static", "./")
	router.StaticFile("/", "./index.html")
	
    
	
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Support both local development and cloud deployment
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
		if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("RENDER") != "" || os.Getenv("DOCKER_ENV") == "true" {
			host = "0.0.0.0" // Listen on all interfaces for cloud deployment or Docker
		}
	}
	
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("🚀 Booking System Server starting on %s", addr)
	log.Printf("🌐 Frontend: http://%s", addr)
	log.Printf("🔌 API: http://%s/api/v1/", addr)
	log.Printf("🧪 Ready for multiplayer concurrency testing!")
	
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}