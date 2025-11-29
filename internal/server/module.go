package server

import (
	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/app"
)

// Module 定义了业务模块的标准接口
// 每个业务模块需要实现此接口以支持模块化注册
type Module interface {
	// Init 从 Container 获取依赖并初始化模块
	// Container 集中管理所有单例依赖，确保组件被多处使用时共享同一实例
	Init(container *app.Container) error
	// RegisterRoutes 注册该模块的所有路由
	RegisterRoutes(router *gin.Engine)
}
