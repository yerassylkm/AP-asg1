package http

import (
	"net/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	uc *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req struct {
		OrderID       string `json:"order_id" binding:"required"`
		Amount        int64  `json:"amount" binding:"required"`
		CustomerEmail string `json:"customer_email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.uc.ProcessPayment(c.Request.Context(), req.OrderID, req.Amount, req.CustomerEmail)
	if err != nil && payment == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	orderID := c.Param("order_id")
	payment, err := h.uc.GetByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}
