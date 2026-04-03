package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"doctorgo/internal/service"
)

type AdminPatientHandler struct {
	service *service.PatientService
}

func NewAdminPatientHandler(s *service.PatientService) *AdminPatientHandler {
	return &AdminPatientHandler{service: s}
}

func (h *AdminPatientHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	success(c, http.StatusOK, items)
}

func (h *AdminPatientHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.service.Get(id)
	if err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	success(c, http.StatusOK, item)
}

func (h *AdminPatientHandler) Create(c *gin.Context) {
	var input service.UpsertPatientInput
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

func (h *AdminPatientHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input service.UpsertPatientInput
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

func (h *AdminPatientHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		fail(c, http.StatusNotFound, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
