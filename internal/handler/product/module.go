package product

import (
	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/app"
)

// ProductModule 商品模块，实现 server.Module 接口
type ProductModule struct {
	handler *ProductHandler
}

// Init 从 Container 获取依赖并初始化商品模块
// Service 已在 Container 中创建为单例，可被多个模块共享
func (m *ProductModule) Init(container *app.Container) error {
	m.handler = NewProductHandler(container.ProductService)
	return nil
}

// RegisterRoutes 注册商品模块的所有路由
func (m *ProductModule) RegisterRoutes(router *gin.Engine) {
	m.handler.RegisterRoutes(router)
}
