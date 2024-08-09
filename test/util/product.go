package util

import (
	"strconv"

	"github.com/dwprz/prasorganic-product-service/src/common/log"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductTest struct {
	db *gorm.DB
}

func NewProductTest(db *gorm.DB) *ProductTest {
	return &ProductTest{
		db: db,
	}
}

func (p *ProductTest) CreateMany() {
	var products []entity.Product

	for i := 1; i <= 20; i++ {
		products = append(products, entity.Product{
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

	p.db.Create(products)
}

func (p *ProductTest) Delete() {
	if err := p.db.Exec("DELETE FROM products;").Error; err != nil {

		log.Logger.WithFields(logrus.Fields{"location": "util.ProductTest/Delete", "section": "db.Exec"}).Error(err)
	}
}
