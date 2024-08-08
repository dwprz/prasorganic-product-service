package util

import (
	"github.com/dwprz/prasorganic-product-service/src/common/log"
	"gorm.io/gorm"
	"github.com/sirupsen/logrus"
)

type ProductTest struct {
	db *gorm.DB
}

func NewProductTest(db *gorm.DB) *ProductTest {
	return &ProductTest{
		db: db,
	}
}

func (p *ProductTest) Delete() {
	if err := p.db.Exec("DELETE FROM products;").Error; err != nil {

		log.Logger.WithFields(logrus.Fields{"location": "util.ProductTest/Delete", "section": "db.Exec"}).Error(err)
	}
}
