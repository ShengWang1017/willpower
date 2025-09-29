package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"willpower-forge-api/internal/models"
)

type CheckInHandler struct {
	db *gorm.DB
}

type CreateCheckInRequest struct {
	GoalID      uint   `json:"goal_id" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=completed failed partial"`
	ReviewNotes string `json:"review_notes"`
}

func NewCheckInHandler(db *gorm.DB) *CheckInHandler {
	return &CheckInHandler{db: db}
}

func (h *CheckInHandler) CreateOrUpdateCheckIn(c *gin.Context) {
	var req CreateCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, 40001, "Invalid input")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, 40102, "Unauthorized")
		return
	}

	var goal models.Goal
	if err := h.db.Where("id = ? AND user_id = ?", req.GoalID, userID).First(&goal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondError(c, http.StatusNotFound, 40401, "Goal not found")
			return
		}
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	today := time.Now().Format("2006-01-02")

	var checkIn models.CheckIn
	err := h.db.Where("goal_id = ? AND date = ?", req.GoalID, today).First(&checkIn).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			checkIn = models.CheckIn{
				GoalID:      req.GoalID,
				UserID:      userID,
				Date:        today,
				Status:      req.Status,
				ReviewNotes: req.ReviewNotes,
			}
			if err := h.db.Create(&checkIn).Error; err != nil {
				respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
				return
			}
			respondSuccess(c, http.StatusCreated, "Check-in recorded", checkIn)
			return
		}
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	checkIn.Status = req.Status
	checkIn.ReviewNotes = req.ReviewNotes
	if err := h.db.Save(&checkIn).Error; err != nil {
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	respondSuccess(c, http.StatusOK, "Check-in recorded", checkIn)
}
