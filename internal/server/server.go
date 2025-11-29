package server

import (
	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/app"
	"github.com/innovationmech/simple-cli/internal/config"
	"github.com/innovationmech/simple-cli/internal/handler/health"
	"github.com/innovationmech/simple-cli/internal/handler/order"
	"github.com/innovationmech/simple-cli/internal/handler/payment"
	"github.com/innovationmech/simple-cli/internal/handler/product"
	"github.com/innovationmech/simple-cli/internal/handler/user"
)

func NewServer() *gin.Engine {
	server := gin.Default()

	// 创建依赖容器，集中管理所有单例依赖
	container, err := app.NewContainer(config.GetDB())
	if err != nil {
		panic(err)
	}

	// 注册所有业务模块
	// 新增模块只需在此切片中追加即可
	// - User/Product: 使用手动依赖注入（通过 Container）
	// - Order: 使用 Google Wire 框架（编译时依赖注入）
	// - Payment: 使用 Uber fx 框架（运行时依赖注入）
	modules := []Module{
		&health.HealthModule{},
		&user.UserModule{},
		&product.ProductModule{},
		&order.OrderModule{},     // 使用 Wire 依赖注入
		&payment.PaymentModule{}, // 使用 fx 依赖注入
	}

	for _, m := range modules {
		if err := m.Init(container); err != nil {
			panic(err)
		}
		m.RegisterRoutes(server)
	}

	return server
}
