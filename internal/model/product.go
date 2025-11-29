package model

import "time"

// Product 商品数据模型
type Product struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

// CreateProductResponse 创建商品响应
type CreateProductResponse struct {
	ID string `json:"id"`
}

// GetProductRequest 获取商品请求
type GetProductRequest struct {
	ID string `uri:"id" binding:"required"`
}

// GetProductResponse 获取商品响应
type GetProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	ID          string  `uri:"id" binding:"required"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

// UpdateProductResponse 更新商品响应
type UpdateProductResponse struct {
	ID string `json:"id"`
}

// DeleteProductRequest 删除商品请求
type DeleteProductRequest struct {
	ID string `uri:"id" binding:"required"`
}

// DeleteProductResponse 删除商品响应
type DeleteProductResponse struct {
	ID string `json:"id"`
}

// ListProductsRequest 商品列表请求
type ListProductsRequest struct {
	Page     int `form:"page" binding:"gte=0"`
	PageSize int `form:"page_size" binding:"gte=0,lte=100"`
}

// ListProductsResponse 商品列表响应
type ListProductsResponse struct {
	Products []GetProductResponse `json:"products"`
	Total    int64                `json:"total"`
}
