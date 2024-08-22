package test

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/delivery"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-product-service/src/mock/service"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	"github.com/dwprz/prasorganic-product-service/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetTestSuite struct {
	suite.Suite
	products       []entity.Product
	productService *service.ProductMock
	restfulServer  *server.Restful
}

// go test -v ./src/core/restful/handler/test/... -count=1 -p=1
// go test -run ^TestHandler_Get$ -v ./src/core/restful/handler/test/ -count=1

func (g *GetTestSuite) SetupSuite() {
	g.productService = &service.ProductMock{}

	imageKitDelivery := delivery.NewImageKit()

	restfulClient := client.NewRestful(imageKitDelivery)
	productHandler := handler.NewProductRESTful(g.productService, restfulClient)
	middleware := middleware.New(restfulClient)

	g.restfulServer = server.New(productHandler, middleware)

	for i := 1; i <= 20; i++ {
		g.products = append(g.products, entity.Product{
			ProductName: "Apel " + strconv.Itoa(i),
			ImageId:     "img" + strconv.Itoa(i),
			Image:       "apel_malang.jpg",
			Price:       100,
			Stock:       50,
			Category:    "FRUIT",
			Length:      10,
			Width:       8,
			Height:      7,
			Weight:      0.2,
			Description: "Apel organik segar pilihan",
		})

	}
}

func (g *GetTestSuite) Test_Success() {
	data := &dto.DataWithPaging[*[]entity.Product]{
		Data: &g.products,
		Paging: &dto.Paging{
			TotalData: 1,
			Page:      1,
			TotalPage: 1,
		},
	}

	g.productService.Mock.On("Get", mock.Anything, &dto.GetProductsReq{
		Page:        1,
		ProductName: "soup",
	}).Return(data, nil)

	request := httptest.NewRequest("GET", "/api/products?page=1&name=soup", nil)
	request.Header.Set("Content-Type", "application/json")

	res, err := g.restfulServer.Test(request)
	assert.NoError(g.T(), err)

	body := util.UnmarshalResponseBody(res.Body)
	products, ok := body["data"].([]any)

	assert.True(g.T(), ok)
	assert.Equal(g.T(), 20, len(products))
}

func TestHandler_Get(t *testing.T) {
	suite.Run(t, new(GetTestSuite))
}
