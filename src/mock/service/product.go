package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	"github.com/stretchr/testify/mock"
)

type ProductMock struct {
	mock.Mock
}

func (p *ProductMock) Create(ctx context.Context, data *dto.CreateProductReq) error {
	argument := p.Mock.Called(ctx, data)

	return argument.Error(0)
}

func (p *ProductMock) Get(ctx context.Context, data *dto.GetProductReq) (*dto.DataWithPaging[*[]entity.Product], error) {
	arguments := p.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.DataWithPaging[*[]entity.Product]), arguments.Error(1)
}

func (p *ProductMock) Update(ctx context.Context, data *dto.UpdateProductReq) (*entity.Product, error) {
	arguments := p.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.Product), arguments.Error(1)
}
