package order

import (
	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/app"
)

// OrderModule 订单模块，实现 server.Module 接口
// 使用 Google Wire 进行依赖注入，而非手动初始化
type OrderModule struct {
	handler *OrderHandler
}

// Init 使用 Wire 自动注入依赖并初始化订单模块
// 与其他模块不同，此模块使用 Wire 框架进行依赖注入
// Wire 会在编译时自动生成依赖注入代码（见 wire_gen.go）
func (m *OrderModule) Init(container *app.Container) error {
	// 使用 Wire 生成的 InitializeOrderHandler 函数
	// Wire 会自动解析依赖链：DB → Repository → Service → Handler
	handler, err := InitializeOrderHandler(container.DB)
	if err != nil {
		return err
	}
	m.handler = handler
	return nil
}

// RegisterRoutes 注册订单模块的所有路由
func (m *OrderModule) RegisterRoutes(router *gin.Engine) {
	m.handler.RegisterRoutes(router)
}
