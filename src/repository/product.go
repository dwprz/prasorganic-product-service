package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	errcustom "github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
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

func (p *ProductImpl) Create(ctx context.Context, data *dto.CreateProductReq) error {
	data.Category = strings.ToUpper(data.Category)

	if err := p.db.WithContext(ctx).Table("products").Create(data).Error; err != nil {

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &errcustom.Response{
				HttpCode: 409,
				GrpcCode: codes.AlreadyExists,
				Message:  "product already exists",
			}
		}

		return err
	}

	return nil
}

func (p *ProductImpl) FindById(ctx context.Context, productId uint) (*entity.Product, error) {
	product := new(entity.Product)

	if err := p.db.WithContext(ctx).Table("products").Where("product_id = ?", productId).First(&product).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errcustom.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "product not found"}
		}

		return nil, err
	}

	return product, nil
}

func (p *ProductImpl) FindManyRandom(ctx context.Context, limit, offset int) (*dto.ProductsWithCountRes, error) {
	queryRes := new(dto.ProductQueryRes)

	query := `
	WITH cte_total_products AS (
		SELECT COUNT(*) FROM products
    ),
    cte_products AS (
    	SELECT
			*
		FROM
			products
        ORDER BY
            random()
		LIMIT $1 OFFSET $2
    )
    SELECT
        (SELECT * FROM cte_total_products) AS total_products,
        (SELECT json_agg(row_to_json(cte_products.*)) FROM cte_products) AS products;
	`

	res := p.db.WithContext(ctx).Raw(query, limit, offset).Find(queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if queryRes.TotalProducts == 0 {
		return nil, &errcustom.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	var products []entity.Product
	if err := json.Unmarshal(queryRes.Products, &products); err != nil {
		return nil, err
	}

	return &dto.ProductsWithCountRes{
		Products:      &products,
		TotalProducts: queryRes.TotalProducts,
	}, nil
}

func (p *ProductImpl) FindManyByCategory(ctx context.Context, category string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	queryRes := new(dto.ProductQueryRes)

	query := `
	WITH cte_total_products AS (
		SELECT
			COUNT(*)
		FROM
			products
		WHERE
			category = $1
    ),
    cte_products AS (
    	SELECT
			*
		FROM
			products
		WHERE
            category = $1
        ORDER BY
            created_at DESC
		LIMIT $2 OFFSET $3
    )
    SELECT
        (SELECT * FROM cte_total_products) AS total_products,
        (SELECT json_agg(row_to_json(cte_products.*)) FROM cte_products) AS products;
	`

	res := p.db.WithContext(ctx).Raw(query, category, limit, offset).Find(queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if queryRes.TotalProducts == 0 {
		return nil, &errcustom.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	var products []entity.Product
	if err := json.Unmarshal(queryRes.Products, &products); err != nil {
		return nil, err
	}

	return &dto.ProductsWithCountRes{
		Products:      &products,
		TotalProducts: queryRes.TotalProducts,
	}, nil
}

func (p *ProductImpl) FindManyByName(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	queryRes := new(dto.ProductQueryRes)
	name = strings.Join(strings.Fields(name), " & ")

	query := `
	WITH cte_total_products AS (
    	SELECT
			COUNT(*)
		FROM
			products
		WHERE
			to_tsvector('indonesian', product_name) @@ to_tsquery('indonesian', ?)
    ),
    cte_products AS (
    	SELECT
			*
		FROM
			products
		WHERE
			to_tsvector('indonesian', product_name) @@ to_tsquery('indonesian', ?)
		LIMIT ? OFFSET ?
    )
    SELECT
        (SELECT * FROM cte_total_products) AS total_products,
        (SELECT json_agg(row_to_json(cte_products.*)) FROM cte_products) AS products;
	`

	res := p.db.WithContext(ctx).Raw(query, name, name, limit, offset).Find(queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes.Products) == 0 {
		return nil, &errcustom.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	var products []entity.Product
	if err := json.Unmarshal(queryRes.Products, &products); err != nil {
		return nil, err
	}

	return &dto.ProductsWithCountRes{
		Products:      &products,
		TotalProducts: queryRes.TotalProducts,
	}, nil
}

func (p *ProductImpl) UpdateById(ctx context.Context, data *dto.UpdateProductReq) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if data.Stock != 0 {
			if err := tx.WithContext(ctx).Raw("SELECT stock FROM products WHERE product_id = $1 FOR UPDATE;", data.ProductId).Error; err != nil {
				return err
			}
		}

		if err := tx.WithContext(ctx).Table("products").Where("product_id = ?", data.ProductId).Updates(data).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
