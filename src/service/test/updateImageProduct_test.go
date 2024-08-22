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
// go test -run ^TestService_UpdateImageProduct$ -v ./src/service/test/ -count=1

type UpdateImageProductTestSuite struct {
	suite.Suite
	productService service.Product
	productRepo    *repository.ProductMock
}

func (u *UpdateImageProductTestSuite) SetupSuite() {
	u.productRepo = repository.NewProductMock()

	u.productService = serviceimpl.NewProduct(u.productRepo)
}

func (u *UpdateImageProductTestSuite) Test_Success() {
	req := u.createUpdateImageProductReq()

	product := u.createProduct()

	u.MockProductRepo_UpdateById(req, nil)
	u.productRepo.Mock.On("FindById", mock.Anything, req.ProductId).Return(product, nil)

	res, err := u.productService.UpdateImage(context.Background(), req)
	assert.NoError(u.T(), err)

	assert.Equal(u.T(), product, res)
}

func (u *UpdateImageProductTestSuite) MockProductRepo_UpdateById(data *dto.UpdateImagePoductReq, returnArg1 error) {

	u.productRepo.Mock.On("UpdateById", mock.Anything, mock.MatchedBy(func(req *entity.Product) bool {
		return data.ProductId == req.ProductId && data.ImageId == req.ImageId && data.Image == req.Image
	})).Return(returnArg1)
}

func (u *UpdateImageProductTestSuite) createUpdateImageProductReq() *dto.UpdateImagePoductReq {
	return &dto.UpdateImagePoductReq{
		ProductId: 1,
		ImageId:   "diaskdmsafsfas",
		Image:     "http://example.com/image.jpg",
	}
}

func (u *UpdateImageProductTestSuite) createProduct() *entity.Product {
	return &entity.Product{
		ProductId:   1,
		ProductName: "Apel Organik",
		ImageId:     "diaskdmsafsfas",
		Image:       "http://example.com/image.jpg",
		Price:       29999,
		Stock:       150,
		Length:      15,
		Width:       7,
		Height:      10,
		Weight:      5,
	}
}

func TestService_UpdateImageProduct(t *testing.T) {
	suite.Run(t, new(UpdateImageProductTestSuite))
}
