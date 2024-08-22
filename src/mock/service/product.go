package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
	"github.com/stretchr/testify/mock"
)

type ProductMock struct {
	mock.Mock
}

func (p *ProductMock) Create(ctx context.Context, data *dto.CreateProductReq) error {
	argument := p.Mock.Called(ctx, data)

	return argument.Error(0)
}

func (p *ProductMock) FindMany(ctx context.Context, data *dto.GetProductsReq) (*dto.DataWithPaging[[]*entity.Product], error) {
	arguments := p.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.DataWithPaging[[]*entity.Product]), arguments.Error(1)
}

func (p *ProductMock) FindManyByIds(ctx context.Context, productIds []uint32) ([]*pb.ProductCart, error) {
	arguments := p.Mock.Called(ctx, productIds)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).([]*pb.ProductCart), arguments.Error(1)
}

func (p *ProductMock) Update(ctx context.Context, data *dto.UpdateProductReq) (*entity.Product, error) {
	arguments := p.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.Product), arguments.Error(1)
}

func (p *ProductMock) UpdateImage(ctx context.Context, data *dto.UpdateImagePoductReq) (*entity.Product, error) {
	arguments := p.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.Product), arguments.Error(1)
}

func (p *ProductMock) ReduceStocks(ctx context.Context, data []*dto.ReduceStocksReq) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (p *ProductMock) RollbackStoks(ctx context.Context, data []*dto.RollbackStoksReq) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}
