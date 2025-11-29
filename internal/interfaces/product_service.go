package interfaces

import (
	"context"

	"github.com/innovationmech/simple-cli/internal/model"
)

// ProductService 商品服务接口
// 定义商品相关的业务操作
type ProductService interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, page, pageSize int) ([]*model.Product, int64, error)
}

