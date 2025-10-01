package routes

import (
	"github.com/gin-gonic/gin"

	"willpower-forge-api/internal/handlers"
	"willpower-forge-api/internal/middleware"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, goalHandler *handlers.GoalHandler, checkInHandler *handlers.CheckInHandler) {
	api := router.Group("/api/v1")

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware())

	authenticated.POST("/goals", goalHandler.CreateGoal)
	authenticated.GET("/goals", goalHandler.GetGoals)
	authenticated.GET("/goals/:id", goalHandler.GetGoalByID)
	authenticated.PUT("/goals/:id", goalHandler.UpdateGoal)
	authenticated.PATCH("/goals/:id/status", goalHandler.UpdateGoalStatus)
	authenticated.DELETE("/goals/:id", goalHandler.DeleteGoal)
	authenticated.GET("/goals/recycle-bin", goalHandler.GetDeletedGoals)
	authenticated.POST("/goals/:id/restore", goalHandler.RestoreGoal)
	authenticated.DELETE("/goals/:id/permanent", goalHandler.PermanentDeleteGoal)

	authenticated.POST("/checkins", checkInHandler.CreateOrUpdateCheckIn)
	authenticated.GET("/checkins", checkInHandler.ListCheckIns)
	authenticated.GET("/checkins/summary", checkInHandler.GoalSummaries)
}
