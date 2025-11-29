package user

import (
	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/app"
)

// UserModule 用户模块，实现 server.Module 接口
type UserModule struct {
	handler *UserHandler
}

// Init 从 Container 获取依赖并初始化用户模块
// Service 已在 Container 中创建为单例，可被多个模块共享
func (m *UserModule) Init(container *app.Container) error {
	m.handler = NewUserHandler(container.UserService)
	return nil
}

// RegisterRoutes 注册用户模块的所有路由
func (m *UserModule) RegisterRoutes(router *gin.Engine) {
	m.handler.RegisterRoutes(router)
}
