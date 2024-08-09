package repository

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/model/dto"
)

type Product interface {
	Create(ctx context.Context, data *dto.CreateProductReq) error
	FindManyRandom(ctx context.Context, limit, offset int) (*dto.ProductsWithCountRes, error)
	FindManyByCategory(ctx context.Context, category string, limit, offset int) (*dto.ProductsWithCountRes, error)
	FindManyByName(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error)
}
