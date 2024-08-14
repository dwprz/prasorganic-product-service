package handler

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (p *ProductGrpcImpl) UpdateStock(ctx context.Context, data *pb.UpdateStockReq) (*emptypb.Empty, error) {
	var req []*dto.UpdateStockReq
	if err := copier.Copy(&req, data.Data); err != nil {
		return nil, err
	}

	err := p.productService.UpdateManyStock(ctx, req)
	return nil, err
}
