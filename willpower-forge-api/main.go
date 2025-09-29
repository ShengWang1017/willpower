package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"willpower-forge-api/internal/database"
	"willpower-forge-api/internal/handlers"
	"willpower-forge-api/internal/routes"
	"willpower-forge-api/internal/services"
)

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

	router := gin.Default()
	router.Use(cors.Default())

	routes.SetupRoutes(router, authHandler, goalHandler, checkInHandler)

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
