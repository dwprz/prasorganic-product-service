package dto

import "github.com/dwprz/prasorganic-product-service/src/model/entity"

type ProductsWithCountRes struct {
	Products      []*entity.Product `json:"products"`
	TotalProducts int               `json:"total_products"`
}

type Paging struct {
	TotalData int `json:"total_data"`
	Page      int `json:"page"`
	TotalPage int `json:"total_page"`
}

type DataWithPaging[T any] struct {
	Data   T       `json:"data"`
	Paging *Paging `json:"paging"`
}
