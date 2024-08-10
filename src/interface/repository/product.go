package repository

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
)

type Product interface {
	Create(ctx context.Context, data *dto.CreateProductReq) error
	FindById(ctx context.Context, productId uint) (*entity.Product, error)
	FindManyRandom(ctx context.Context, limit, offset int) (*dto.ProductsWithCountRes, error)
	FindManyByCategory(ctx context.Context, category string, limit, offset int) (*dto.ProductsWithCountRes, error)
	FindManyByName(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error)
	UpdateById(ctx context.Context, data *entity.Product) error
}
