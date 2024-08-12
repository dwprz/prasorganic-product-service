package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/common/helper"
	"github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
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

func (p *ProductImpl) FindMany(ctx context.Context, data *dto.GetProductReq) (*dto.DataWithPaging[[]*entity.Product], error) {
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

func (p *ProductImpl) FindManyByIds(ctx context.Context, productIds []uint32) ([]*pb.ProductCart, error) {
	if err := p.validate.Var(productIds, `dive,required`); err != nil {
		return nil, err
	}

	res, err := p.productRepo.FindManyByIds(ctx, productIds)
	return res, err
}

func (p *ProductImpl) Update(ctx context.Context, data *dto.UpdateProductReq) (*entity.Product, error) {
	if err := p.validate.Struct(data); err != nil {
		return nil, err
	}
	product := new(entity.Product)
	if err := copier.Copy(product, data); err != nil {
		return nil, err
	}

	if err := p.productRepo.UpdateById(ctx, product); err != nil {
		return nil, err
	}

	res, err := p.productRepo.FindById(ctx, data.ProductId)
	return res, err
}

func (p *ProductImpl) UpdateImage(ctx context.Context, data *dto.UpdateProductImageReq) (*entity.Product, error) {
	if err := p.validate.Struct(data); err != nil {
		return nil, err
	}

	err := p.productRepo.UpdateById(ctx, &entity.Product{
		ProductId: data.ProductId,
		ImageId:   data.ImageId,
		Image:     data.Image,
	})

	if err != nil {
		return nil, err
	}

	res, err := p.productRepo.FindById(ctx, data.ProductId)
	return res, err
}
