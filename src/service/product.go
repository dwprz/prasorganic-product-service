package service

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/common/helper"
	v "github.com/dwprz/prasorganic-product-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
	"github.com/jinzhu/copier"
)

type ProductImpl struct {
	productRepo repository.Product
}

func NewProduct(pr repository.Product) service.Product {
	return &ProductImpl{
		productRepo: pr,
	}
}

func (p *ProductImpl) Create(ctx context.Context, data *dto.CreateProductReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	err := p.productRepo.Create(ctx, data)
	return err
}

func (p *ProductImpl) FindMany(ctx context.Context, data *dto.GetProductsReq) (*dto.DataWithPaging[[]*entity.Product], error) {
	if err := v.Validate.Struct(data); err != nil {
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
	if err := v.Validate.Var(productIds, `dive,required`); err != nil {
		return nil, err
	}

	res, err := p.productRepo.FindManyByIds(ctx, productIds)
	return res, err
}

func (p *ProductImpl) Update(ctx context.Context, data *dto.UpdateProductReq) (*entity.Product, error) {
	if err := v.Validate.Struct(data); err != nil {
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

func (p *ProductImpl) UpdateImage(ctx context.Context, data *dto.UpdateImagePoductReq) (*entity.Product, error) {
	if err := v.Validate.Struct(data); err != nil {
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

func (p *ProductImpl) ReduceStocks(ctx context.Context, data []*dto.ReduceStocksReq) error {
	if err := v.Validate.Var(data, `required,dive`); err != nil {
		return err
	}

	err := p.productRepo.ReduceStocks(ctx, data)
	return err
}

func (p *ProductImpl) RollbackStoks(ctx context.Context, data []*dto.RollbackStoksReq) error{
	if err := v.Validate.Var(data, `required,dive`); err != nil {
		return err
	}

	err := p.productRepo.RollbackStocks(ctx, data)
	return err
}