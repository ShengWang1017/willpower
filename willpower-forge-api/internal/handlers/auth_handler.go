package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"willpower-forge-api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, 40001, "Invalid input")
		return
	}

	if err := h.authService.RegisterUser(req.Username, req.Password); err != nil {
		switch err {
		case services.ErrUserExists:
			respondError(c, http.StatusConflict, 40901, "Username already exists")
		default:
			respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		}
		return
	}

	respondSuccess(c, http.StatusCreated, "User registered successfully", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, 40001, "Invalid input")
		return
	}

	token, userID, err := h.authService.LoginUser(req.Username, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredential:
			respondError(c, http.StatusUnauthorized, 40101, "Invalid username or password")
		default:
			respondError(c, http.StatusInternalServerError, 50001, "Internal server error")
		}
		return
	}

	respondSuccess(c, http.StatusOK, "Login successful", gin.H{
		"token":   token,
		"user_id": userID,
	})
}
