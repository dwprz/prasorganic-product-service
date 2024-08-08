package database

import (
	"github.com/dwprz/prasorganic-product-service/src/common/log"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres() *gorm.DB {
	dsn := config.Conf.Postgres.Dsn

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "database.NewPostgres", "section": "gorm.Open"}).Fatal(err)
	}

	return db
}
