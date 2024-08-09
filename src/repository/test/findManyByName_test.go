package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/database"
	repointerface "github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/repository"
	"github.com/dwprz/prasorganic-product-service/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_FindManyByName$ -v ./src/repository/test -count=1

type FindManyByNameTestSuite struct {
	suite.Suite
	productRepo     repointerface.Product
	postgresDB      *gorm.DB
	productTestUtil *util.ProductTest
}

func (f *FindManyByNameTestSuite) SetupSuite() {
	f.postgresDB = database.NewPostgres()
	f.productRepo = repository.NewProduct(f.postgresDB)
	f.productTestUtil = util.NewProductTest(f.postgresDB)

	f.productTestUtil.CreateMany()
}

func (f *FindManyByNameTestSuite) TearDownSuite() {
	f.productTestUtil.Delete()

	sqlDB, _ := f.postgresDB.DB()
	sqlDB.Close()
}

func (f *FindManyByNameTestSuite) Test_Success() {

	res, err := f.productRepo.FindManyByName(context.Background(), "apel", 10, 0)
	assert.NoError(f.T(), err)

	assert.Equal(f.T(), 10, len(*res.Products))
	assert.Equal(f.T(), 20, res.TotalProducts)
}

func (f *FindManyByNameTestSuite) Test_NotFound() {

	res, err := f.productRepo.FindManyByName(context.Background(), "not found", 10, 0)
	assert.Error(f.T(), err)

	errorRes := &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	assert.Equal(f.T(), errorRes, err)

	assert.Nil(f.T(), res)
}

func TestRepository_FindManyByName(t *testing.T) {
	suite.Run(t, new(FindManyByNameTestSuite))
}
