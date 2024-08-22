package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/mock/repository"
	serviceimpl "github.com/dwprz/prasorganic-product-service/src/service"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_FindManyProductByIds$ -v ./src/service/test/ -count=1

type FindManyProductByIdsTestSuite struct {
	suite.Suite
	productService service.Product
	productRepo    *repository.ProductMock
}

func (f *FindManyProductByIdsTestSuite) SetupSuite() {
	f.productRepo = repository.NewProductMock()

	f.productService = serviceimpl.NewProduct(f.productRepo)
}

func (f *FindManyProductByIdsTestSuite) Test_Success() {
	productIds := []uint32{1, 2, 3, 4, 5, 6}
	productsCart := f.createProductsCart()

	f.productRepo.Mock.On("FindManyByIds", mock.Anything, productIds).Return(productsCart, nil)

	res, err := f.productService.FindManyByIds(context.Background(), productIds)
	assert.NoError(f.T(), err)

	assert.Equal(f.T(), productsCart, res)
}

func (f *FindManyProductByIdsTestSuite) createProductsCart() []*pb.ProductCart {
	return []*pb.ProductCart{
		{
			ProductId:   1,
			ProductName: "Apel Organik",
			Image:       "http://example.com/image.jpg",
			Price:       29999,
			Stock:       150,
			Length:      15,
			Width:       7,
			Height:      10,
			Weight:      5,
		},
		{
			ProductId:   2,
			ProductName: "Jeruk Organik",
			Image:       "http://example.com/image.jpg",
			Price:       29999,
			Stock:       150,
			Length:      15,
			Width:       7,
			Height:      10,
			Weight:      5,
		},
	}
}

func TestService_FindManyProductByIds(t *testing.T) {
	suite.Run(t, new(FindManyProductByIdsTestSuite))
}
