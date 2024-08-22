package repository

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

func NewProductMock() *ProductMock {
	return &ProductMock{
		Mock: mock.Mock{},
	}
}

func (p *ProductMock) Create(ctx context.Context, data *dto.CreateProductReq) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (p *ProductMock) FindById(ctx context.Context, productId uint) (*entity.Product, error) {
	arguments := p.Mock.Called(ctx, productId)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.Product), arguments.Error(1)
}

func (p *ProductMock) FindManyByIds(ctx context.Context, productIds []uint32) ([]*pb.ProductCart, error) {
	arguments := p.Mock.Called(ctx, productIds)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).([]*pb.ProductCart), arguments.Error(1)
}

func (p *ProductMock) FindManyRandom(ctx context.Context, limit, offset int) (*dto.ProductsWithCountRes, error) {
	arguments := p.Mock.Called(ctx, limit, offset)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.ProductsWithCountRes), arguments.Error(1)
}

func (p *ProductMock) FindManyByCategory(ctx context.Context, category string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	arguments := p.Mock.Called(ctx, category, limit, offset)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.ProductsWithCountRes), arguments.Error(1)
}

func (p *ProductMock) FindManyByName(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	arguments := p.Mock.Called(ctx, name, limit, offset)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.ProductsWithCountRes), arguments.Error(1)
}

func (p *ProductMock) UpdateById(ctx context.Context, data *entity.Product) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (p *ProductMock) ReduceStocks(ctx context.Context, data []*dto.ReduceStocksReq) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (p *ProductMock) RollbackStocks(ctx context.Context, data []*dto.RollbackStoksReq) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}
