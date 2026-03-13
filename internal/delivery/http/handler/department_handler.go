package handler

import (
	"net/http"
	"strconv"

	depart_usecases "github.com/genss333/go-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	uc *depart_usecases.DepartmentUsecase
}

func NewDepartmentHandler(uc *depart_usecases.DepartmentUsecase) *DepartmentHandler {
	return &DepartmentHandler{uc: uc}
}

func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dept, err := h.uc.GetDepartmentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dept)
}
