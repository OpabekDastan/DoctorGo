package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"doctor_go/internal/service"
)

type AdminAppointmentHandler struct {
	service *service.AppointmentService
}

func NewAdminAppointmentHandler(s *service.AppointmentService) *AdminAppointmentHandler {
	return &AdminAppointmentHandler{service: s}
}

func (h *AdminAppointmentHandler) List(c *gin.Context) {
	items, err := h.service.ListAdmin()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	success(c, http.StatusOK, items)
}

func (h *AdminAppointmentHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.GetAdmin(id)
	if err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}

func (h *AdminAppointmentHandler) Create(c *gin.Context) {
	var input service.UpsertAppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	createdBy := c.GetInt64("user_id")
	item, err := h.service.Create(createdBy, input)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	success(c, http.StatusCreated, item)
}

func (h *AdminAppointmentHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input service.UpsertAppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.Update(id, input)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}

func (h *AdminAppointmentHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
