package repository

import (
	"context"
	"github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type ProductImpl struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) repository.Product {
	return &ProductImpl{
		db: db,
	}
}

func (p *ProductImpl) Create(ctx context.Context, data *dto.CreateReq) error {
	if err := p.db.WithContext(ctx).Table("products").Create(data).Error; err != nil {

		if errPG, ok := err.(*pgconn.PgError); ok && errPG.Code == "23505" {
			return &errors.Response{
				HttpCode: 409,
				GrpcCode: codes.AlreadyExists,
				Message:  "product already exists",
			}
		}

		return err
	}

	return nil
}
