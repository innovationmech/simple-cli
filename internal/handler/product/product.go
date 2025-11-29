package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/types"
)

// ProductHandler 商品 HTTP 处理器
type ProductHandler struct {
	productService interfaces.ProductService
}

// NewProductHandler 创建商品处理器实例
func NewProductHandler(productService interfaces.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// CreateProduct 创建商品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request model.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	product := &model.Product{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := h.productService.CreateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create product",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusCreated,
			Message: "Product created successfully",
		},
		Data: model.CreateProductResponse{
			ID: product.ID,
		},
	})
}

// GetProduct 获取商品详情
func (h *ProductHandler) GetProduct(c *gin.Context) {
	var request model.GetProductRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusNotFound,
				Message: "Product not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Product retrieved successfully",
		},
		Data: model.GetProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		},
	})
}

// UpdateProduct 更新商品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var request model.UpdateProductRequest
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

	// 先获取现有商品
	product, err := h.productService.GetProduct(c.Request.Context(), request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusNotFound,
				Message: "Product not found",
			},
		})
		return
	}

	// 更新字段
	if request.Name != "" {
		product.Name = request.Name
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price > 0 {
		product.Price = request.Price
	}
	if request.Stock >= 0 {
		product.Stock = request.Stock
	}

	if err := h.productService.UpdateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update product",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Product updated successfully",
		},
		Data: model.UpdateProductResponse{
			ID: product.ID,
		},
	})
}

// DeleteProduct 删除商品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	var request model.DeleteProductRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request: " + err.Error(),
			},
		})
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete product",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Product deleted successfully",
		},
		Data: model.DeleteProductResponse{
			ID: request.ID,
		},
	})
}

// ListProducts 获取商品列表
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var request model.ListProductsRequest
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

	products, total, err := h.productService.ListProducts(c.Request.Context(), request.Page, request.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to list products",
			},
		})
		return
	}

	// 转换响应
	var productResponses []model.GetProductResponse
	for _, p := range products {
		productResponses = append(productResponses, model.GetProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
		})
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "Products retrieved successfully",
		},
		Data: model.ListProductsResponse{
			Products: productResponses,
			Total:    total,
		},
	})
}

// RegisterRoutes 注册商品相关路由
func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	products := router.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("", h.ListProducts)
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}

