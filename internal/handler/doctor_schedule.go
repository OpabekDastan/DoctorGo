package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"doctor_go/internal/service"
)

type DoctorScheduleHandler struct {
	service *service.ScheduleService
}

func NewDoctorScheduleHandler(s *service.ScheduleService) *DoctorScheduleHandler {
	return &DoctorScheduleHandler{service: s}
}

func (h *DoctorScheduleHandler) List(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	items, err := h.service.List(doctorID)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	success(c, http.StatusOK, items)
}

func (h *DoctorScheduleHandler) Get(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.Get(id, doctorID)
	if err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}

func (h *DoctorScheduleHandler) Create(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	var input service.UpsertScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.Create(doctorID, input)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	success(c, http.StatusCreated, item)
}

func (h *DoctorScheduleHandler) Update(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input service.UpsertScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.Update(id, doctorID, input)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}

func (h *DoctorScheduleHandler) Delete(c *gin.Context) {
	doctorID := c.GetInt64("user_id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id, doctorID); err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
