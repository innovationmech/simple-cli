package repository

import (
	"context"

	"github.com/innovationmech/simple-cli/internal/model"
	"gorm.io/gorm"
)

// ProductRepository 商品数据访问接口
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, offset, limit int) ([]*model.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓储实例
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) DeleteProduct(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", id).Error
}

func (r *productRepository) ListProducts(ctx context.Context, offset, limit int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

