package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/types"
)

// PaymentHandler 支付 HTTP 处理器
type PaymentHandler struct {
	paymentService interfaces.PaymentService
}

// NewPaymentHandler 创建支付处理器实例
// 此函数将作为 fx Provider 使用
func NewPaymentHandler(paymentService interfaces.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment 创建支付
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var request model.CreatePaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	payment := &model.Payment{
		ID:      uuid.New().String(),
		OrderID: request.OrderID,
		UserID:  request.UserID,
		Amount:  request.Amount,
		Method:  request.Method,
	}

	paymentURL, err := h.paymentService.CreatePayment(c.Request.Context(), payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create payment: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusCreated,
			Message: "Payment created successfully",
		},
		Data: model.CreatePaymentResponse{
			ID:         payment.ID,
			PaymentURL: paymentURL,
		},
	})
}

// GetPayment 获取支付详情
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	var request model.GetPaymentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	payment, err := h.paymentService.GetPayment(c.Request.Context(), request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusNotFound,
				Message: "Payment not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Payment retrieved successfully",
		},
		Data: model.GetPaymentResponse{
			ID:            payment.ID,
			OrderID:       payment.OrderID,
			UserID:        payment.UserID,
			Amount:        payment.Amount,
			Method:        payment.Method,
			Status:        payment.Status,
			TransactionID: payment.TransactionID,
			PaidAt:        payment.PaidAt,
			CreatedAt:     payment.CreatedAt,
		},
	})
}

// PaymentCallback 支付回调（模拟第三方回调）
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	var request model.PaymentCallbackRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	if err := h.paymentService.ProcessCallback(c.Request.Context(), request.PaymentID, request.TransactionID, request.Success); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to process callback: " + err.Error(),
			},
		})
		return
	}

	status := model.PaymentStatusSuccess
	if !request.Success {
		status = model.PaymentStatusFailed
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Callback processed successfully",
		},
		Data: model.PaymentCallbackResponse{
			ID:     request.PaymentID,
			Status: status,
		},
	})
}

// RefundPayment 退款
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	var request model.RefundRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		// JSON body is optional for refund
	}

	if err := h.paymentService.RefundPayment(c.Request.Context(), request.ID, request.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to refund payment: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Payment refunded successfully",
		},
		Data: model.RefundResponse{
			ID:     request.ID,
			Status: model.PaymentStatusRefunded,
		},
	})
}

// ListPayments 获取支付记录列表
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	var request model.ListPaymentsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	// 设置默认值
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	payments, total, err := h.paymentService.ListPayments(c.Request.Context(), request.UserID, request.OrderID, request.Page, request.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to list payments",
			},
		})
		return
	}

	// 转换响应
	var paymentResponses []model.GetPaymentResponse
	for _, p := range payments {
		paymentResponses = append(paymentResponses, model.GetPaymentResponse{
			ID:            p.ID,
			OrderID:       p.OrderID,
			UserID:        p.UserID,
			Amount:        p.Amount,
			Method:        p.Method,
			Status:        p.Status,
			TransactionID: p.TransactionID,
			PaidAt:        p.PaidAt,
			CreatedAt:     p.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Payments retrieved successfully",
		},
		Data: model.ListPaymentsResponse{
			Payments: paymentResponses,
			Total:    total,
		},
	})
}

// RegisterRoutes 注册支付相关路由
func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	payments := router.Group("/payments")
	{
		payments.POST("", h.CreatePayment)
		payments.GET("", h.ListPayments)
		payments.GET("/:id", h.GetPayment)
		payments.POST("/callback", h.PaymentCallback)
		payments.POST("/:id/refund", h.RefundPayment)
	}
}
