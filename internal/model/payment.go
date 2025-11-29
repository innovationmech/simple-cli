package model

import "time"

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PaymentMethodAlipay     PaymentMethod = "alipay"
	PaymentMethodWechat     PaymentMethod = "wechat"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBalance    PaymentMethod = "balance"
)

// Payment 支付记录数据模型
type Payment struct {
	ID            string        `json:"id" gorm:"primaryKey"`
	OrderID       string        `json:"order_id" gorm:"index"`
	UserID        string        `json:"user_id" gorm:"index"`
	Amount        float64       `json:"amount"`
	Method        PaymentMethod `json:"method"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transaction_id"` // 第三方交易号
	PaidAt        *time.Time    `json:"paid_at"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID string        `json:"order_id" binding:"required"`
	UserID  string        `json:"user_id" binding:"required"`
	Amount  float64       `json:"amount" binding:"required,gt=0"`
	Method  PaymentMethod `json:"method" binding:"required"`
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	ID         string `json:"id"`
	PaymentURL string `json:"payment_url"` // 支付跳转链接
}

// GetPaymentRequest 获取支付详情请求
type GetPaymentRequest struct {
	ID string `uri:"id" binding:"required"`
}

// GetPaymentResponse 获取支付详情响应
type GetPaymentResponse struct {
	ID            string        `json:"id"`
	OrderID       string        `json:"order_id"`
	UserID        string        `json:"user_id"`
	Amount        float64       `json:"amount"`
	Method        PaymentMethod `json:"method"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transaction_id"`
	PaidAt        *time.Time    `json:"paid_at"`
	CreatedAt     time.Time     `json:"created_at"`
}

// PaymentCallbackRequest 支付回调请求（模拟）
type PaymentCallbackRequest struct {
	PaymentID     string `json:"payment_id" binding:"required"`
	TransactionID string `json:"transaction_id" binding:"required"`
	Success       bool   `json:"success"`
}

// PaymentCallbackResponse 支付回调响应
type PaymentCallbackResponse struct {
	ID     string        `json:"id"`
	Status PaymentStatus `json:"status"`
}

// RefundRequest 退款请求
type RefundRequest struct {
	ID     string `uri:"id" binding:"required"`
	Reason string `json:"reason"`
}

// RefundResponse 退款响应
type RefundResponse struct {
	ID     string        `json:"id"`
	Status PaymentStatus `json:"status"`
}

// ListPaymentsRequest 支付记录列表请求
type ListPaymentsRequest struct {
	UserID   string `form:"user_id"`
	OrderID  string `form:"order_id"`
	Page     int    `form:"page" binding:"gte=0"`
	PageSize int    `form:"page_size" binding:"gte=0,lte=100"`
}

// ListPaymentsResponse 支付记录列表响应
type ListPaymentsResponse struct {
	Payments []GetPaymentResponse `json:"payments"`
	Total    int64                `json:"total"`
}
