package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/database"
	repointerface "github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/repository"
	"github.com/dwprz/prasorganic-product-service/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_Create$ -v ./src/repository/test -count=1

type CreateTestSuite struct {
	suite.Suite
	productRepo     repointerface.Product
	postgresDB      *gorm.DB
	productTestUtil *util.ProductTest
}

func (c *CreateTestSuite) SetupSuite() {
	c.postgresDB = database.NewPostgres()
	c.productRepo = repository.NewProduct(c.postgresDB)
	c.productTestUtil = util.NewProductTest(c.postgresDB)
}

func (c *CreateTestSuite) TearDownTest() {
	c.productTestUtil.Delete()

}

func (c *CreateTestSuite) TearDownSuite() {
	sqlDB, _ := c.postgresDB.DB()
	sqlDB.Close()
}

func (c *CreateTestSuite) Test_Success() {
	req := &dto.CreateProductReq{
		ProductName: "apel hijau",
		ImageId:     "uweh28hsajewu2ijs",
		Image:       "https://example-prasor.com/image",
		Category:    "fruit",
		Price:       10000,
		Stock:       250,
		Length:      25,
		Width:       20,
		Height:      15,
		Weight:      5000,
		Description: "example description",
	}

	err := c.productRepo.Create(context.Background(), req)
	assert.NoError(c.T(), err)
}

func (c *CreateTestSuite) Test_AlreadyExists() {
	req := &dto.CreateProductReq{
		ProductName: "apel hijau",
		ImageId:     "uweh28hsajewu2ijs",
		Image:       "https://example-prasor.com/image",
		Category:    "fruit",
		Price:       10000,
		Stock:       250,
		Length:      25,
		Width:       20,
		Height:      15,
		Weight:      5000,
		Description: "example description",
	}

	c.productRepo.Create(context.Background(), req)

	err := c.productRepo.Create(context.Background(), req)
	assert.Error(c.T(), err)

	errorRes := &errors.Response{HttpCode: 409, GrpcCode: codes.AlreadyExists, Message: "product already exists"}
	assert.Equal(c.T(), errorRes, err)
}

func TestRepository_Create(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}
