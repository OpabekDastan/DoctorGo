package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"doctorgo/internal/service"
)

type AdminDoctorHandler struct {
	service *service.DoctorService
}

func NewAdminDoctorHandler(s *service.DoctorService) *AdminDoctorHandler {
	return &AdminDoctorHandler{service: s}
}

func (h *AdminDoctorHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	success(c, http.StatusOK, items)
}

func (h *AdminDoctorHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.Get(id)
	if err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}

func (h *AdminDoctorHandler) Create(c *gin.Context) {
	var input service.CreateDoctorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.Create(input)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	success(c, http.StatusCreated, item)
}

func (h *AdminDoctorHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input service.UpdateDoctorInput
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

func (h *AdminDoctorHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
