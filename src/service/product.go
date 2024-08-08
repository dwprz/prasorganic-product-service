package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
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

func (p *ProductImpl) Create(ctx context.Context, data *dto.CreateReq) error {
	if err := p.validate.Struct(data); err != nil {
		return err
	}

	err := p.productRepo.Create(ctx, data)
	return err
}
