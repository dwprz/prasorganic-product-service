package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/common/helper"
	"github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	"github.com/go-playground/validator/v10"
)

type ProductImpl struct {
	validate    *validator.Validate
	productRepo repository.Product
}

func NewProduct(v *validator.Validate, pr repository.Product) service.Product {
	return &ProductImpl{
		validate:    v,
		productRepo: pr,
	}
}

func (p *ProductImpl) Create(ctx context.Context, data *dto.CreateProductReq) error {
	if err := p.validate.Struct(data); err != nil {
		return err
	}

	err := p.productRepo.Create(ctx, data)
	return err
}

func (p *ProductImpl) Get(ctx context.Context, data *dto.GetProductReq) (*dto.DataWithPaging[*[]entity.Product], error) {
	if err := p.validate.Struct(data); err != nil {
		return nil, err
	}

	limit, offset := helper.CreateLimitAndOffset(data.Page)

	var res *dto.ProductsWithCountRes
	var err error

	switch {
	case data.Category != "":
		res, err = p.productRepo.FindManyByCategory(ctx, data.Category, limit, offset)
	case data.ProductName != "":
		res, err = p.productRepo.FindManyByName(ctx, data.ProductName, limit, offset)
	default:
		res, err = p.productRepo.FindManyRandom(ctx, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	return helper.FormatPagedData(res.Products, res.TotalProducts, data.Page, limit), nil
}
