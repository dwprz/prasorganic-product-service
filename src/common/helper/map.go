package helper

import "github.com/dwprz/prasorganic-product-service/src/model/entity"

func MapProductQueryToEntities(data []*entity.ProductQueryRes) (products []*entity.Product, totalProducts int) {
	for _, product := range data {
		products = append(products, &entity.Product{
			ProductId:   product.ProductId,
			ProductName: product.ProductName,
			ImageId:     product.ImageId,
			Image:       product.Image,
			Rating:      product.Rating,
			Sold:        product.Sold,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			Length:      product.Length,
			Width:       product.Width,
			Height:      product.Height,
			Weight:      product.Weight,
			Description: product.Description,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		})
	}

	return products, data[0].TotalProducts
}
