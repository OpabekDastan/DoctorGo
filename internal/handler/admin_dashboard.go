package handler

import (
	"net/http"

	"doctor_go/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminDashboardHandler struct {
	svc *service.DashboardService
}

func NewAdminDashboardHandler(svc *service.DashboardService) *AdminDashboardHandler {
	return &AdminDashboardHandler{svc: svc}
}

func (h *AdminDashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		fail(c, http.StatusInternalServerError, "could not fetch stats")
		return
	}
	success(c, http.StatusOK, stats)
}
