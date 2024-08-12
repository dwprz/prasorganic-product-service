package handler

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
)

type ProductGrpcImpl struct {
	productService service.Product
	pb.UnimplementedProductServiceServer
}

func NewProductGrpc(ps service.Product) pb.ProductServiceServer {
	return &ProductGrpcImpl{
		productService: ps,
	}
}

func (p *ProductGrpcImpl) FindManyByIdsForCart(ctx context.Context, data *pb.ProductIds) (*pb.ProductsCartResponse, error) {
	res, err := p.productService.FindManyByIds(ctx, data.Ids)
	if err != nil {
		return nil, err
	}

	return &pb.ProductsCartResponse{
		Data: res,
	}, nil
}
