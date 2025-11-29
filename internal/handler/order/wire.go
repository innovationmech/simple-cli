//go:build wireinject
// +build wireinject

package order

import (
	"github.com/google/wire"
	"github.com/innovationmech/simple-cli/internal/repository"
	orderSrv "github.com/innovationmech/simple-cli/internal/service/order"
	"gorm.io/gorm"
)

// OrderProviderSet 是 Order 模块的依赖提供者集合
// 包含了构建 OrderHandler 所需的所有依赖
var OrderProviderSet = wire.NewSet(
	repository.NewOrderRepository,
	repository.NewProductRepository,
	orderSrv.NewOrderService,
	NewOrderHandler,
)

// InitializeOrderHandler 使用 Wire 初始化 OrderHandler
// Wire 会根据 ProviderSet 自动生成依赖注入代码
func InitializeOrderHandler(db *gorm.DB) (*OrderHandler, error) {
	wire.Build(OrderProviderSet)
	return nil, nil
}
