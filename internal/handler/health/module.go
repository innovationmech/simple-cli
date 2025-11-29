package health

import (
	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/app"
)

// HealthModule 健康检查模块，实现 server.Module 接口
type HealthModule struct{}

// Init 初始化健康检查模块（该模块无需额外依赖）
func (m *HealthModule) Init(container *app.Container) error {
	return nil
}

// RegisterRoutes 注册健康检查路由
func (m *HealthModule) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", HealthCheck)
}
