package repository

import (
	"context"
	"errors"
	"strings"

	errcustom "github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/common/helper"
	"github.com/dwprz/prasorganic-product-service/src/interface/repository"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
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

func (p *ProductImpl) Create(ctx context.Context, data *dto.CreateProductReq) error {
	data.Category = strings.ToUpper(data.Category)

	if err := p.db.WithContext(ctx).Table("products").Create(data).Error; err != nil {
		if errPG, ok := err.(*pgconn.PgError); ok && errPG.Code == "23505" {
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

func (p *ProductImpl) FindManyByIds(ctx context.Context, productIds []uint32) ([]*pb.ProductCart, error) {
	var products []*pb.ProductCart

	if err := p.db.WithContext(ctx).Table("products").Where("product_id in ?", productIds).Scan(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductImpl) FindManyRandom(ctx context.Context, limit, offset int) (*dto.ProductsWithCountRes, error) {
	var queryRes []*entity.ProductQueryRes

	query := `
	WITH cte_total_products AS (
		SELECT COUNT(*) AS total_products FROM products
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
	SELECT ctp.total_products, cp.* FROM cte_total_products AS ctp CROSS JOIN cte_products AS cp;
	`

	res := p.db.WithContext(ctx).Raw(query, limit, offset).Find(&queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errcustom.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	products, total := helper.MapProductQueryToEntities(queryRes)

	return &dto.ProductsWithCountRes{
		Products:      products,
		TotalProducts: total,
	}, nil
}

func (p *ProductImpl) FindManyByCategory(ctx context.Context, category string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	var queryRes []*entity.ProductQueryRes

	query := `
	WITH cte_total_products AS (
		SELECT
			COUNT(*) AS total_products
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
   	SELECT ctp.total_products, cp.* FROM cte_total_products AS ctp CROSS JOIN cte_products AS cp;
	`

	res := p.db.WithContext(ctx).Raw(query, category, limit, offset).Find(&queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errcustom.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	products, total := helper.MapProductQueryToEntities(queryRes)

	return &dto.ProductsWithCountRes{
		Products:      products,
		TotalProducts: total,
	}, nil
}

func (p *ProductImpl) FindManyByName(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	var queryRes []*entity.ProductQueryRes
	name = strings.Join(strings.Fields(name), " & ")

	query := `
	WITH cte_total_products AS (
    	SELECT
			COUNT(*) AS total_products
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
	SELECT ctp.total_products, cp.* FROM cte_total_products AS ctp CROSS JOIN cte_products AS cp;
	`

	res := p.db.WithContext(ctx).Raw(query, name, name, limit, offset).Find(&queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errcustom.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	products, total := helper.MapProductQueryToEntities(queryRes)

	return &dto.ProductsWithCountRes{
		Products:      products,
		TotalProducts: total,
	}, nil
}

func (p *ProductImpl) UpdateById(ctx context.Context, data *entity.Product) error {

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if data.Stock != 0 {
			if err := tx.Exec("SELECT stock FROM products WHERE product_id = $1 FOR UPDATE;", data.ProductId).Error; err != nil {
				return err
			}
		}

		if err := tx.Table("products").Updates(data).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (p *ProductImpl) ReduceStocks(ctx context.Context, data []*dto.ReduceStocksReq) error {

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ids, err := helper.GetProductIds(data)
		if err != nil {
			return err
		}

		var products []*entity.Product

		if err := tx.Raw("SELECT product_id, product_name, stock FROM products WHERE product_id in (?) FOR UPDATE;", ids).Scan(&products).Error; err != nil {
			return err
		}

		if err := helper.CheckStockProducts(data, products); err != nil {
			return err
		}

		query, args := helper.BuildQueryReduceStocks(data)

		if err := tx.Exec(query, args...).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (p *ProductImpl) RollbackStocks(ctx context.Context, data []*dto.RollbackStoksReq) error {

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ids, err := helper.GetProductIds(data)
		if err != nil {
			return err
		}

		if err := tx.Exec("SELECT stock FROM products WHERE product_id in (?) FOR UPDATE;", ids).Error; err != nil {
			return err
		}

		query, args := helper.BuildQueryRollbackStocks(data)

		if err := tx.Exec(query, args...).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
