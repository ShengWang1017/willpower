package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"willpower-forge-api/internal/models"
)

type GoalHandler struct {
	db *gorm.DB
}

type CreateGoalRequest struct {
	Type  string `json:"type" binding:"required,oneof=I_WILL I_WONT I_WANT"`
	Title string `json:"title" binding:"required,min=1,max=255"`
}

func NewGoalHandler(db *gorm.DB) *GoalHandler {
	return &GoalHandler{db: db}
}

func (h *GoalHandler) CreateGoal(c *gin.Context) {
	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, 40001, "Invalid input")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, 40102, "Unauthorized")
		return
	}

	goal := models.Goal{
		UserID: userID,
		Type:   req.Type,
		Title:  req.Title,
		Status: "active",
	}

	if err := h.db.Create(&goal).Error; err != nil {
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	respondSuccess(c, http.StatusCreated, "Goal created", goal)
}

func (h *GoalHandler) GetGoals(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, 40102, "Unauthorized")
		return
	}

	var goals []models.Goal
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error; err != nil {
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	respondSuccess(c, http.StatusOK, "Success", goals)
}

func (h *GoalHandler) GetGoalByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, 40102, "Unauthorized")
		return
	}

	goalIDParam := c.Param("id")
	goalID, err := strconv.ParseUint(goalIDParam, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, 40001, "Invalid goal id")
		return
	}

	goalIDUint := uint(goalID)

	var goal models.Goal
	if err := h.db.Where("id = ? AND user_id = ?", goalIDUint, userID).First(&goal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondError(c, http.StatusNotFound, 40401, "Goal not found")
			return
		}
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	respondSuccess(c, http.StatusOK, "Success", goal)
}

func getUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := val.(uint)
	return userID, ok
}
