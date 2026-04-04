package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"doctor_go/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials), errors.Is(err, service.ErrInactiveUser):
			fail(c, http.StatusUnauthorized, err.Error())
		default:
			fail(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	success(c, http.StatusOK, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt64("user_id")
	user, err := h.authService.Me(userID)
	if err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	success(c, http.StatusOK, user)
}
