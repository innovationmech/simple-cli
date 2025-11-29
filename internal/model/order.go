package model

import "time"

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order 订单数据模型
type Order struct {
	ID          string      `json:"id" gorm:"primaryKey"`
	UserID      string      `json:"user_id" gorm:"index"`
	ProductID   string      `json:"product_id" gorm:"index"`
	Quantity    int         `json:"quantity"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	ID          string  `json:"id"`
	TotalAmount float64 `json:"total_amount"`
}

// GetOrderRequest 获取订单请求
type GetOrderRequest struct {
	ID string `uri:"id" binding:"required"`
}

// GetOrderResponse 获取订单响应
type GetOrderResponse struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	ProductID   string      `json:"product_id"`
	Quantity    int         `json:"quantity"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	ID     string      `uri:"id" binding:"required"`
	Status OrderStatus `json:"status" binding:"required"`
}

// UpdateOrderStatusResponse 更新订单状态响应
type UpdateOrderStatusResponse struct {
	ID     string      `json:"id"`
	Status OrderStatus `json:"status"`
}

// ListOrdersRequest 订单列表请求
type ListOrdersRequest struct {
	UserID   string `form:"user_id"`
	Page     int    `form:"page" binding:"gte=0"`
	PageSize int    `form:"page_size" binding:"gte=0,lte=100"`
}

// ListOrdersResponse 订单列表响应
type ListOrdersResponse struct {
	Orders []GetOrderResponse `json:"orders"`
	Total  int64              `json:"total"`
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	ID string `uri:"id" binding:"required"`
}

// CancelOrderResponse 取消订单响应
type CancelOrderResponse struct {
	ID     string      `json:"id"`
	Status OrderStatus `json:"status"`
}

