package payment

import (
	"github.com/innovationmech/simple-cli/internal/repository"
	paymentSrv "github.com/innovationmech/simple-cli/internal/service/payment"
	"go.uber.org/fx"
)

// FxModule 是 Payment 模块的 fx 依赖提供者模块
// 使用 Uber fx 框架进行运行时依赖注入
//
// fx 与 Wire 的区别：
// - Wire: 编译时依赖注入，生成代码
// - fx: 运行时依赖注入，基于反射
var FxModule = fx.Module("payment",
	// 提供 Repository
	fx.Provide(repository.NewPaymentRepository),
	fx.Provide(repository.NewOrderRepository),

	// 提供 Service
	fx.Provide(paymentSrv.NewPaymentService),

	// 提供 Handler
	fx.Provide(NewPaymentHandler),
)

// FxResult 封装 fx 注入的结果
// 用于从 fx.App 中提取 PaymentHandler
type FxResult struct {
	fx.Out
	Handler *PaymentHandler
}
