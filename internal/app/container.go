package app

import (
	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/repository"
	productSrv "github.com/innovationmech/simple-cli/internal/service/product"
	userSrv "github.com/innovationmech/simple-cli/internal/service/user"
	"gorm.io/gorm"
)

// Container 集中管理所有依赖，确保单例
// 当组件被多处使用时，通过 Container 共享同一实例
type Container struct {
	DB *gorm.DB

	// Repositories
	UserRepo    repository.UserRepository
	ProductRepo repository.ProductRepository

	// Services
	UserService    interfaces.UserService
	ProductService interfaces.ProductService
}

// NewContainer 创建并初始化依赖容器
// 按照依赖顺序初始化：DB → Repositories → Services
func NewContainer(db *gorm.DB) (*Container, error) {
	c := &Container{DB: db}

	// 初始化 Repositories
	c.UserRepo = repository.NewUserRepository(db)
	c.ProductRepo = repository.NewProductRepository(db)

	// 初始化 Services
	var err error
	c.UserService, err = userSrv.NewUserService(userSrv.WithUserRepository(c.UserRepo))
	if err != nil {
		return nil, err
	}

	c.ProductService, err = productSrv.NewProductService(productSrv.WithProductRepository(c.ProductRepo))
	if err != nil {
		return nil, err
	}

	return c, nil
}
