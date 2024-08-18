package helper

import (
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
)

func BuildQueryReduceStocks(data []*dto.ReduceStocksReq) (query string, args []any) {
	query = "UPDATE products SET stock = CASE"

	args = []any{}
	for _, product := range data {
		query += " WHEN product_id = ? THEN stock - ?"
		args = append(args, product.ProductId, product.Quantity)
	}

	query += " END WHERE product_id IN ("

	for i, product := range data {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args = append(args, product.ProductId)
	}

	query += ");"

	return query, args
}

func BuildQueryRollbackStocks(data []*dto.RollbackStoksReq) (query string, args []any) {
	query = "UPDATE products SET stock = CASE"

	args = []any{}
	for _, product := range data {
		query += " WHEN product_id = ? THEN stock + ?"
		args = append(args, product.ProductId, product.Quantity)
	}

	query += " END WHERE product_id IN ("

	for i, product := range data {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args = append(args, product.ProductId)
	}

	query += ");"

	return query, args
}