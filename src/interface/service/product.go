package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
)

type Product interface {
	Create(ctx context.Context, data *dto.CreateProductReq) error
	FindMany(ctx context.Context, data *dto.GetProductsReq) (*dto.DataWithPaging[[]*entity.Product], error)
	FindManyByIds(ctx context.Context, productIds []uint32) ([]*pb.ProductCart, error)
	Update(ctx context.Context, data *dto.UpdateProductReq) (*entity.Product, error)
	UpdateImage(ctx context.Context, data *dto.UpdateImagePoductReq) (*entity.Product, error)
	ReduceStocks(ctx context.Context, data []*dto.ReduceStocksReq) error 
	RollbackStoks(ctx context.Context, data []*dto.RollbackStoksReq) error
}
