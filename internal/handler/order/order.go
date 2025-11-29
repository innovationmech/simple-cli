package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/types"
)

// OrderHandler 订单 HTTP 处理器
type OrderHandler struct {
	orderService interfaces.OrderService
}

// NewOrderHandler 创建订单处理器实例
// 此函数将作为 Wire Provider 使用
func NewOrderHandler(orderService interfaces.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request model.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	order := &model.Order{
		ID:        uuid.New().String(),
		UserID:    request.UserID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}

	if err := h.orderService.CreateOrder(c.Request.Context(), order); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create order: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusCreated,
			Message: "Order created successfully",
		},
		Data: model.CreateOrderResponse{
			ID:          order.ID,
			TotalAmount: order.TotalAmount,
		},
	})
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	var request model.GetOrderRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusNotFound,
				Message: "Order not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Order retrieved successfully",
		},
		Data: model.GetOrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			ProductID:   order.ProductID,
			Quantity:    order.Quantity,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
		},
	})
}

// UpdateOrderStatus 更新订单状态
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var request model.UpdateOrderStatusRequest
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
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	if err := h.orderService.UpdateOrderStatus(c.Request.Context(), request.ID, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update order status: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Order status updated successfully",
		},
		Data: model.UpdateOrderStatusResponse{
			ID:     request.ID,
			Status: request.Status,
		},
	})
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	var request model.CancelOrderRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	if err := h.orderService.CancelOrder(c.Request.Context(), request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to cancel order: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Order cancelled successfully",
		},
		Data: model.CancelOrderResponse{
			ID:     request.ID,
			Status: model.OrderStatusCancelled,
		},
	})
}

// ListOrders 获取订单列表
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var request model.ListOrdersRequest
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

	orders, total, err := h.orderService.ListOrdersByUser(c.Request.Context(), request.UserID, request.Page, request.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to list orders",
			},
		})
		return
	}

	// 转换响应
	var orderResponses []model.GetOrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, model.GetOrderResponse{
			ID:          o.ID,
			UserID:      o.UserID,
			ProductID:   o.ProductID,
			Quantity:    o.Quantity,
			TotalAmount: o.TotalAmount,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Orders retrieved successfully",
		},
		Data: model.ListOrdersResponse{
			Orders: orderResponses,
			Total:  total,
		},
	})
}

// RegisterRoutes 注册订单相关路由
func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	orders := router.Group("/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("", h.ListOrders)
		orders.GET("/:id", h.GetOrder)
		orders.PUT("/:id/status", h.UpdateOrderStatus)
		orders.POST("/:id/cancel", h.CancelOrder)
	}
}

