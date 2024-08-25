package helper

import (
	"fmt"

	"github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	"google.golang.org/grpc/codes"
)

func GetProductIds(products any) (ids []uint, err error) {
	switch v := products.(type) {
	case []*dto.ReduceStocksReq:
		for _, item := range v {
			ids = append(ids, item.ProductId)
		}

	case []*dto.RollbackStoksReq:
		for _, item := range v {
			ids = append(ids, item.ProductId)
		}
	default:
		return nil, fmt.Errorf("unexpected type %T (product ids)", products)
	}

	return ids, nil
}

func CheckStockProducts(orders []*dto.ReduceStocksReq, products []*entity.Product) error {

	if len(products) == 0 {
		return &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	productsMap := make(map[uint]*entity.Product)
	for _, product := range products {
		productsMap[product.ProductId] = product
	}

	for _, order := range orders {

		product, exists := productsMap[order.ProductId]
		if !exists {
			msg := fmt.Sprintf("product with id %d not found", order.ProductId)
			return &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: msg}
		}

		if order.Quantity > int(product.Stock) {
			msg := fmt.Sprintf("not enough stock product %s (id: %d)", product.ProductName, order.ProductId)
			return &errors.Response{HttpCode: 400, GrpcCode: codes.FailedPrecondition, Message: msg}
		}
	}

	return nil
}
