package main

import (
	"its-complicated/internal/handlers"
	"its-complicated/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize dependencies
	userStorage := storage.NewUserStorage()
	userHandler := handlers.NewUserHandler(userStorage)

	// 2. Set up router
	router := gin.Default()

	// 3. Define routes
	api := router.Group("/api")
	{
		api.POST("/register", userHandler.RegisterUser)
	}

	// 4. Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
