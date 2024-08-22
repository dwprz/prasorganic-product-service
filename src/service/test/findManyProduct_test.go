package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/mock/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	serviceimpl "github.com/dwprz/prasorganic-product-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_FindManyProduct$ -v ./src/service/test/ -count=1

type FindManyProductTestSuite struct {
	suite.Suite
	productService service.Product
	productRepo    *repository.ProductMock
}

func (f *FindManyProductTestSuite) SetupSuite() {
	f.productRepo = repository.NewProductMock()

	f.productService = serviceimpl.NewProduct(f.productRepo)
}

func (f *FindManyProductTestSuite) Test_Random() {
	req := f.createGetProductsReq()
	productsWithCountRes := f.createProductsWithCountRes()

	f.productRepo.Mock.On("FindManyRandom", mock.Anything, mock.Anything, mock.Anything).Return(productsWithCountRes, nil)

	res, err := f.productService.FindMany(context.Background(), req)
	assert.NoError(f.T(), err)

	assert.Equal(f.T(), productsWithCountRes.Products, res.Data)
}

func (f *FindManyProductTestSuite) Test_ByCategory() {
	req := f.createGetProductsReq()
	req.Category = "FRUIT"

	productsWithCountRes := f.createProductsWithCountRes()
	f.productRepo.Mock.On("FindManyByCategory", mock.Anything, req.Category, mock.Anything, mock.Anything).Return(productsWithCountRes, nil)

	res, err := f.productService.FindMany(context.Background(), req)
	assert.NoError(f.T(), err)

	assert.Equal(f.T(), productsWithCountRes.Products, res.Data)
}

func (f *FindManyProductTestSuite) Test_ByName() {
	req := f.createGetProductsReq()
	req.ProductName = "Organik"

	productsWithCountRes := f.createProductsWithCountRes()
	f.productRepo.Mock.On("FindManyByName", mock.Anything, req.ProductName, mock.Anything, mock.Anything).Return(productsWithCountRes, nil)

	res, err := f.productService.FindMany(context.Background(), req)
	assert.NoError(f.T(), err)

	assert.Equal(f.T(), productsWithCountRes.Products, res.Data)
}

func (f *FindManyProductTestSuite) Test_LimitPage() {
	req := f.createGetProductsReq()
	req.Page = 1000

	res, err := f.productService.FindMany(context.Background(), req)
	assert.Error(f.T(), err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(f.T(), ok)

	assert.Equal(f.T(), validationErr[0].Field(), "Page")
	assert.Nil(f.T(), res)
}

func (f *FindManyProductTestSuite) createGetProductsReq() *dto.GetProductsReq {
	return &dto.GetProductsReq{
		Page:        1,
		ProductName: "",
		Category:    "",
	}
}

func (f *FindManyProductTestSuite) createProductsWithCountRes() *dto.ProductsWithCountRes {
	products := f.createProducts()

	return &dto.ProductsWithCountRes{
		Products:      products,
		TotalProducts: len(products),
	}
}

func (f *FindManyProductTestSuite) createProducts() []*entity.Product {
	return []*entity.Product{
		{
			ProductId:   1,
			ProductName: "Apel Organik",
			ImageId:     "image1234567890",
			Image:       "http://example.com/image.jpg",
			Price:       29999,
			Stock:       150,
			Category:    "FRUIT",
			Length:      15,
			Width:       7,
			Height:      10,
			Weight:      5,
			Description: "Apel organik berkualitas dan terjangkau",
		},
		{
			ProductId:   2,
			ProductName: "Jeruk Organik",
			ImageId:     "image1234567890",
			Image:       "http://example.com/image.jpg",
			Price:       29999,
			Stock:       150,
			Category:    "FRUIT",
			Length:      15,
			Width:       7,
			Height:      10,
			Weight:      5,
			Description: "Jeruk organik berkualitas dan terjangkau",
		},
	}
}

func TestService_FindManyProduct(t *testing.T) {
	suite.Run(t, new(FindManyProductTestSuite))
}
