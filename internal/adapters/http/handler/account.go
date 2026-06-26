package handler

import (
	"account-service/internal/domain/entity"
	"account-service/internal/domain/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service ports.AccountService
}

func NewAccountHandler(service ports.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/accounts", h.CreateAccount)
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var acc entity.Account
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Save(c.Request.Context(), acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, acc)
}
