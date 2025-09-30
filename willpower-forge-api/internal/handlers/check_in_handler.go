package handlers

import (
	"net/http"
	"strconv"
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

	checkIn := models.CheckIn{
		GoalID:      req.GoalID,
		UserID:      userID,
		Date:        time.Now().Format("2006-01-02"),
		Status:      req.Status,
		ReviewNotes: req.ReviewNotes,
	}

	if err := h.db.Create(&checkIn).Error; err != nil {
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	respondSuccess(c, http.StatusCreated, "Check-in recorded", checkIn)
}

func (h *CheckInHandler) ListCheckIns(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, 40102, "Unauthorized")
		return
	}

	goalIDParam := c.Query("goal_id")
	if goalIDParam == "" {
		respondError(c, http.StatusBadRequest, 40001, "goal_id is required")
		return
	}

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

	var checkIns []models.CheckIn
	if err := h.db.Where("goal_id = ? AND user_id = ?", goalIDUint, userID).
		Order("date DESC").Find(&checkIns).Error; err != nil {
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	respondSuccess(c, http.StatusOK, "Success", checkIns)
}

type goalSummaryRow struct {
	GoalID uint
	Status string
	Count  int64
}

type GoalSummary struct {
	GoalID    uint   `json:"goal_id"`
	Title     string `json:"title"`
	Completed int64  `json:"completed"`
	Partial   int64  `json:"partial"`
	Failed    int64  `json:"failed"`
}

func (h *CheckInHandler) GoalSummaries(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, 40102, "Unauthorized")
		return
	}

	dateFilter := c.Query("date")
	if dateFilter != "" {
		if _, err := time.Parse("2006-01-02", dateFilter); err != nil {
			respondError(c, http.StatusBadRequest, 40001, "Invalid date format")
			return
		}
	}

	var goals []models.Goal
	if err := h.db.Where("user_id = ?", userID).Order("created_at ASC").Find(&goals).Error; err != nil {
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	summaries := make([]GoalSummary, 0, len(goals))
	summaryMap := make(map[uint]int)

	for idx, goal := range goals {
		summaries = append(summaries, GoalSummary{
			GoalID: goal.ID,
			Title:  goal.Title,
		})
		summaryMap[goal.ID] = idx
	}

	if len(goals) == 0 {
		respondSuccess(c, http.StatusOK, "Success", summaries)
		return
	}

	var rows []goalSummaryRow
	query := h.db.Model(&models.CheckIn{}).
		Select("goal_id, status, COUNT(*) AS count").
		Where("user_id = ?", userID)

	if dateFilter != "" {
		query = query.Where("date = ?", dateFilter)
	}

	if err := query.Group("goal_id, status").Find(&rows).Error; err != nil {
		respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		return
	}

	for _, row := range rows {
		idx, exists := summaryMap[row.GoalID]
		if !exists {
			continue
		}
		switch row.Status {
		case "completed":
			summaries[idx].Completed = row.Count
		case "partial":
			summaries[idx].Partial = row.Count
		case "failed":
			summaries[idx].Failed = row.Count
		}
	}

	respondSuccess(c, http.StatusOK, "Success", summaries)
}
