package payment

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/app"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// PaymentModule 支付模块，实现 server.Module 接口
// 使用 Uber fx 进行依赖注入，而非手动初始化或 Wire
type PaymentModule struct {
	handler *PaymentHandler
}

// Init 使用 fx 自动注入依赖并初始化支付模块
// 与其他模块不同，此模块使用 Uber fx 框架进行运行时依赖注入
//
// fx 的特点：
// - 运行时依赖注入（基于反射）
// - 支持生命周期管理（OnStart/OnStop）
// - 支持模块化组织依赖
// - 自动解析依赖图
func (m *PaymentModule) Init(container *app.Container) error {
	var handler *PaymentHandler

	// 创建 fx App 来解析依赖
	// fx.NopLogger 用于禁用 fx 的默认日志输出
	fxApp := fx.New(
		fx.NopLogger,

		// 提供数据库连接（从 Container 获取）
		fx.Provide(func() *gorm.DB {
			return container.DB
		}),

		// 加载 Payment 模块的所有 Provider
		FxModule,

		// 提取构建好的 Handler
		fx.Populate(&handler),
	)

	// 启动 fx App（触发依赖构建）
	if err := fxApp.Start(context.Background()); err != nil {
		return err
	}

	m.handler = handler
	return nil
}

// RegisterRoutes 注册支付模块的所有路由
func (m *PaymentModule) RegisterRoutes(router *gin.Engine) {
	m.handler.RegisterRoutes(router)
}
