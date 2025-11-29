package product

import (
	"context"
	"errors"

	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/repository"
)

// ProductSrv 是 ProductService 接口的别名，方便外部引用
type ProductSrv = interfaces.ProductService

// ProductServiceConfig 商品服务配置
type ProductServiceConfig struct {
	ProductRepository repository.ProductRepository
}

// ProductServiceOption 函数式选项模式
type ProductServiceOption func(*ProductServiceConfig)

type productService struct {
	config *ProductServiceConfig
}

// WithProductRepository 注入商品仓储依赖
func WithProductRepository(repo repository.ProductRepository) ProductServiceOption {
	return func(config *ProductServiceConfig) {
		config.ProductRepository = repo
	}
}

// NewProductService 创建商品服务实例
// 使用函数式选项模式注入依赖
func NewProductService(opts ...ProductServiceOption) (ProductSrv, error) {
	config := &ProductServiceConfig{}
	for _, opt := range opts {
		opt(config)
	}
	if config.ProductRepository == nil {
		return nil, errors.New("product repository is required")
	}
	return &productService{config: config}, nil
}

func (s *productService) CreateProduct(ctx context.Context, product *model.Product) error {
	return s.config.ProductRepository.CreateProduct(ctx, product)
}

func (s *productService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	return s.config.ProductRepository.GetProduct(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *model.Product) error {
	return s.config.ProductRepository.UpdateProduct(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.config.ProductRepository.DeleteProduct(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context, page, pageSize int) ([]*model.Product, int64, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 10 // 默认每页 10 条
	}
	return s.config.ProductRepository.ListProducts(ctx, offset, pageSize)
}

