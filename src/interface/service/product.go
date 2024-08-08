package service

import (
	"context"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
)

type Product interface {
	Create(ctx context.Context, data *dto.CreateReq) error
}
