package main

import (
	"log"
	"os"

	"github.com/conor/wiki-notes-backend/database"
	"github.com/conor/wiki-notes-backend/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Create Gin router
	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Initialize handlers
	pageHandler := handlers.NewPageHandler(database.DB)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/pages", pageHandler.GetAllPages)
		api.GET("/page/:id", pageHandler.GetPageByID)
		api.GET("/page/by-name/:name", pageHandler.GetPageByName)
		api.POST("/page", pageHandler.CreatePage)
		api.PATCH("/page/:id", pageHandler.UpdatePage)
		api.DELETE("/page/:id", pageHandler.DeletePage)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
