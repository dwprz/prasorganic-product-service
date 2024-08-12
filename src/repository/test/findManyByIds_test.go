package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/infrastructure/database"
	repointerface "github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/repository"
	"github.com/dwprz/prasorganic-product-service/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// go test -v ./src/repository/test/... -count=1 -p=1
// go test -run ^TestRepository_FindManyByIds$ -v ./src/repository/test -count=1

type FindManyByIdsTestSuite struct {
	suite.Suite
	productIds      []uint32
	productRepo     repointerface.Product
	postgresDB      *gorm.DB
	productTestUtil *util.ProductTest
}

func (f *FindManyByIdsTestSuite) SetupSuite() {
	f.postgresDB = database.NewPostgres()
	f.productRepo = repository.NewProduct(f.postgresDB)
	f.productTestUtil = util.NewProductTest(f.postgresDB)

	f.productIds = f.productTestUtil.CreateMany()
}

func (f *FindManyByIdsTestSuite) TearDownSuite() {
	f.productTestUtil.Delete()

	sqlDB, _ := f.postgresDB.DB()
	sqlDB.Close()
}

func (f *FindManyByIdsTestSuite) Test_Success() {

	res, err := f.productRepo.FindManyByIds(context.Background(), f.productIds)

	assert.NoError(f.T(), err)
	assert.Equal(f.T(), 20, len(res))
}

func TestRepository_FindManyByIds(t *testing.T) {
	suite.Run(t, new(FindManyByIdsTestSuite))
}
