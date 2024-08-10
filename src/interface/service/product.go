package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
)

type Product interface {
	Create(ctx context.Context, data *dto.CreateProductReq) error
	Get(ctx context.Context, data *dto.GetProductReq) (*dto.DataWithPaging[*[]entity.Product], error)
	Update(ctx context.Context, data *dto.UpdateProductReq) (*entity.Product, error)
	UpdateImage(ctx context.Context, data *dto.UpdateProductImageReq) (*entity.Product, error) 
}
