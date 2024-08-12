package benchmark

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// *cd current directory
// go test -v -bench=. -count=1 -p=1

var postgres *gorm.DB

func init() {
	postgres = database.NewPostgres()
}

func withCTE(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error) {
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

	res := postgres.WithContext(ctx).Raw(query, name, name, limit, offset).Find(queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes.Products) == 0 {
		return nil, &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	var products []*entity.Product
	if err := json.Unmarshal(queryRes.Products, &products); err != nil {
		return nil, err
	}

	return &dto.ProductsWithCountRes{
		Products:      products,
		TotalProducts: queryRes.TotalProducts,
	}, nil
}

func nonCTE(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	var products []*entity.Product
	name = strings.Join(strings.Fields(name), " & ")

	query := `
	SELECT
		*
	FROM
		products
	WHERE	
		to_tsvector('indonesian', product_name) @@ to_tsquery('indonesian', ?)
	LIMIT ? OFFSET ?;
	`

	if err := postgres.WithContext(ctx).Raw(query, name, limit, offset).Scan(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	var totalProducts int

	query = `
	SELECT
		COUNT(*) AS total_products
	FROM
		products
	WHERE
		to_tsvector('indonesian', product_name) @@ to_tsquery('indonesian', ?);
	`

	if err := postgres.WithContext(ctx).Raw(query, name).Scan(&totalProducts).Error; err != nil {
		return nil, err
	}

	return &dto.ProductsWithCountRes{
		Products:      products,
		TotalProducts: totalProducts,
	}, nil
}

func Benchmark_CompareQueryCTE(b *testing.B) {
	b.Run("With CTE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			withCTE(context.Background(), "soup", 20, 0)
		}
	})

	b.Run("Non CTE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nonCTE(context.Background(), "soup", 20, 0)
		}
	})
}

// 1 ms = 1.000.000 ns
// 1 s = 1000 ms
//================================ With CTE ================================
// test 1:
// Benchmark_CompareQuery/With_CTE-12                  3024            374551 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.188s

// test 2:
// Benchmark_CompareQuery/With_CTE-12                  2852            377084 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.133s

// test 3:
// Benchmark_CompareQuery/With_CTE-12                  3081            373036 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.205s

//================================ Non CTE ================================
// test 1:
// Benchmark_CompareQuery/Non_CTE-12                   2654            432769 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.210s

// test 2:
// Benchmark_CompareQuery/Non_CTE-12                   2568            434005 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.176s

// test 3:
// Benchmark_CompareQuery/Non_CTE-12                   2546            432254 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.164s