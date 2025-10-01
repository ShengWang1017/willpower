package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"willpower-forge-api/internal/database"
	"willpower-forge-api/internal/handlers"
	"willpower-forge-api/internal/routes"
	"willpower-forge-api/internal/services"
)

//go:embed web/dist
var webFS embed.FS

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	database.AutoMigrateModels(db)

	authService := services.NewAuthService(db)
	authHandler := handlers.NewAuthHandler(authService)
	goalHandler := handlers.NewGoalHandler(db)
	checkInHandler := handlers.NewCheckInHandler(db)

	// Start scheduled cleanup service
	cleanupService := services.NewCleanupService(db)
	cleanupService.StartScheduledCleanup()

	router := gin.Default()
	router.Use(cors.Default())

	routes.SetupRoutes(router, authHandler, goalHandler, checkInHandler)

	// Serve embedded static files
	staticFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatalf("failed to load static files: %v", err)
	}

	// Serve static files for all non-API routes
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Check if file exists in embedded FS
		if _, err := staticFS.Open(path[1:]); err == nil {
			c.FileFromFS(path, http.FS(staticFS))
			return
		}
		// Serve index.html for SPA routes
		c.FileFromFS("/", http.FS(staticFS))
	})

	// Get port from environment variable, default to 5173
	port := os.Getenv("PORT")
	if port == "" {
		port = "5173"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
