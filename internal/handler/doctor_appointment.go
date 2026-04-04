package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"doctorgo/internal/service"
)

type DoctorAppointmentHandler struct {
	service *service.AppointmentService
}

func NewDoctorAppointmentHandler(s *service.AppointmentService) *DoctorAppointmentHandler {
	return &DoctorAppointmentHandler{service: s}
}

func (h *DoctorAppointmentHandler) List(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	items, err := h.service.ListDoctor(doctorID)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	success(c, http.StatusOK, items)
}

func (h *DoctorAppointmentHandler) Get(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.GetDoctor(id, doctorID)
	if err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}

func (h *DoctorAppointmentHandler) SetStatus(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.SetStatus(id, doctorID, req.Status)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}
