package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/mock/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	serviceimpl "github.com/dwprz/prasorganic-product-service/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_UpdateProduct$ -v ./src/service/test/ -count=1

type UpdateProductTestSuite struct {
	suite.Suite
	productService service.Product
	productRepo    *repository.ProductMock
}

func (u *UpdateProductTestSuite) SetupSuite() {
	u.productRepo = repository.NewProductMock()

	u.productService = serviceimpl.NewProduct(u.productRepo)
}

func (u *UpdateProductTestSuite) Test_Success() {
	req := u.createUpdateProductReq()
	req.ProductName = "new product name"

	product := u.createProduct()
	product.ProductName = req.ProductName

	u.MockProductRepo_UpdateById(req, nil)
	u.productRepo.Mock.On("FindById", mock.Anything, req.ProductId).Return(product, nil)

	res, err := u.productService.Update(context.Background(), req)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), product, res)
}

func (u *UpdateProductTestSuite) MockProductRepo_UpdateById(data *dto.UpdateProductReq, returnArg1 error) {

	u.productRepo.Mock.On("UpdateById", mock.Anything, mock.MatchedBy(func(req *entity.Product) bool {
		return data.ProductId == req.ProductId && data.ProductName == req.ProductName
	})).Return(returnArg1)
}

func (u *UpdateProductTestSuite) createUpdateProductReq() *dto.UpdateProductReq {
	return &dto.UpdateProductReq{
		ProductId: 1,
	}
}

func (u *UpdateProductTestSuite) createProduct() *entity.Product {
	return &entity.Product{
		ProductId:   1,
		ProductName: "Apel Organik",
		Image:       "http://example.com/image.jpg",
		Price:       29999,
		Stock:       150,
		Length:      15,
		Width:       7,
		Height:      10,
		Weight:      5,
	}
}

func TestService_UpdateProduct(t *testing.T) {
	suite.Run(t, new(UpdateProductTestSuite))
}
