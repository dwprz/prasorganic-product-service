package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/mock/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	serviceimpl "github.com/dwprz/prasorganic-product-service/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_RollbackStockProducts$ -v ./src/service/test/ -count=1

type rollbackStockProductsTestSuite struct {
	suite.Suite
	productService service.Product
	productRepo    *repository.ProductMock
}

func (u *rollbackStockProductsTestSuite) SetupSuite() {
	u.productRepo = repository.NewProductMock()

	u.productService = serviceimpl.NewProduct(u.productRepo)
}

func (u *rollbackStockProductsTestSuite) Test_Success() {
	req := u.createRollbackStocksReq()

	u.productRepo.Mock.On("RollbackStocks", mock.Anything, req).Return(nil)

	err := u.productService.RollbackStoks(context.Background(), req)
	assert.NoError(u.T(), err)
}

func (u *rollbackStockProductsTestSuite) createRollbackStocksReq() []*dto.RollbackStoksReq {
	return []*dto.RollbackStoksReq{
		{
			ProductId: 1,
			Quantity:  10,
		},
		{
			ProductId: 2,
			Quantity:  20,
		},
	}
}

func TestService_RollbackStockProducts(t *testing.T) {
	suite.Run(t, new(rollbackStockProductsTestSuite))
}
