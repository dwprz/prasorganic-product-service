package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/mock/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	serviceimpl "github.com/dwprz/prasorganic-product-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_CreateProduct$ -v ./src/service/test/ -count=1

type CreateProductTestSuite struct {
	suite.Suite
	productService service.Product
	productRepo    *repository.ProductMock
}

func (c *CreateProductTestSuite) SetupSuite() {
	c.productRepo = repository.NewProductMock()

	c.productService = serviceimpl.NewProduct(c.productRepo)
}

func (c *CreateProductTestSuite) Test_Succcess() {
	req := c.createProductReq()
	c.productRepo.Mock.On("Create", mock.Anything, req).Return(nil)

	err := c.productService.Create(context.Background(), req)
	assert.NoError(c.T(), err)
}

func (c *CreateProductTestSuite) Test_WithoutPrice() {
	req := c.createProductReq()
	req.Price = 0

	err := c.productService.Create(context.Background(), req)
	assert.Error(c.T(), err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(c.T(), ok)

	assert.Equal(c.T(), validationErr[0].Field(), "Price")
}

func (c *CreateProductTestSuite) createProductReq() *dto.CreateProductReq {
	return &dto.CreateProductReq{
		ProductName: "Apel Organik",
		ImageId:     "image1234567890",
		Image:       "http://example.com/image.jpg",
		Price:       29999,
		Stock:       150,
		Category:    "FRUITS",
		Length:      15,
		Width:       7,
		Height:      10,
		Weight:      5,
		Description: "Apel organik berkualitas dan terjangkau",
	}
}

func TestService_CreateProduct(t *testing.T) {
	suite.Run(t, new(CreateProductTestSuite))
}
