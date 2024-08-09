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
// go test -run ^TestRepository_FindManyRandom$ -v ./src/repository/test -count=1

type FindManyRandomTestSuite struct {
	suite.Suite
	productRepo     repointerface.Product
	postgresDB      *gorm.DB
	productTestUtil *util.ProductTest
}

func (f *FindManyRandomTestSuite) SetupSuite() {
	f.postgresDB = database.NewPostgres()
	f.productRepo = repository.NewProduct(f.postgresDB)
	f.productTestUtil = util.NewProductTest(f.postgresDB)

	f.productTestUtil.CreateMany()
}

func (f *FindManyRandomTestSuite) TearDownSuite() {
	f.productTestUtil.Delete()

	sqlDB, _ := f.postgresDB.DB()
	sqlDB.Close()
}

func (f *FindManyRandomTestSuite) Test_Success() {

	res, err := f.productRepo.FindManyRandom(context.Background(), 20, 0)
	assert.NoError(f.T(), err)

	assert.Equal(f.T(), 20, len(*res.Products))
	assert.Equal(f.T(), 20, res.TotalProducts)
}

func TestRepository_FindManyRandom(t *testing.T) {
	suite.Run(t, new(FindManyRandomTestSuite))
}
