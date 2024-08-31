package benchmark

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/common/helper"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/dwprz/prasorganic-product-service/src/model/entity"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// *cd current directory
// go test -v -bench=. -count=1 -p=1

var db *gorm.DB

func init() {
	db = database.NewPostgres()
}

func withCTE(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error) {
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

	res := db.WithContext(ctx).Raw(query, name, name, limit, offset).Find(&queryRes)

	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "products not found"}
	}

	products, total := helper.MapProductQueryToEntities(queryRes)

	return &dto.ProductsWithCountRes{
		Products:      products,
		TotalProducts: total,
	}, nil
}

type ProductQueryRes struct {
	Products      []byte
	TotalProducts int
}

func cteWithJsonAgg(ctx context.Context, name string, limit, offset int) (*dto.ProductsWithCountRes, error) {
	queryRes := new(ProductQueryRes)
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

	res := db.WithContext(ctx).Raw(query, name, name, limit, offset).Find(queryRes)

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

	if err := db.WithContext(ctx).Raw(query, name, limit, offset).Scan(&products).Error; err != nil {
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

	if err := db.WithContext(ctx).Raw(query, name).Scan(&totalProducts).Error; err != nil {
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

	b.Run("CTE With JSON Agg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cteWithJsonAgg(context.Background(), "soup", 20, 0)
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
// Benchmark_CompareQueryCTE/With_CTE-12               1550            748378 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.270s

// test 2:
// Benchmark_CompareQueryCTE/With_CTE-12               1454            737744 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.189s

// test 3:
// Benchmark_CompareQueryCTE/With_CTE-12               1485            758148 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.240s

//================================ CTE With JSON Agg ================================
// test 1:
// Benchmark_CompareQueryCTE/CTE_With_JSON_Agg-12              1220            864279 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.186s

// test 2:
// Benchmark_CompareQueryCTE/CTE_With_JSON_Agg-12              1376            861641 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   2.320s

// test 3:
// Benchmark_CompareQueryCTE/CTE_With_JSON_Agg-12              1327            871121 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.278s

//================================ Non CTE ================================
// test 1:
// Benchmark_CompareQueryCTE/Non_CTE-12                 986           1082766 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.223s

// test 2:
// Benchmark_CompareQueryCTE/Non_CTE-12                1069           1084290 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.305s

// test 3:
// Benchmark_CompareQueryCTE/Non_CTE-12                1024           1068560 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-product-service/src/repository/benchmark   1.245s
