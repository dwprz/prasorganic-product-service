package restful

import (
	"github.com/dwprz/prasorganic-product-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-product-service/src/interface/delivery"
	"github.com/dwprz/prasorganic-product-service/src/interface/service"
)

func Initialize(ps service.Product, dik delivery.ImageKitRestful) *server.Restful {
	restfulClient := client.New(dik)
	productHandler := handler.NewProduct(ps, restfulClient)

	middleware := middleware.New(restfulClient)
	restfulServer := server.New(productHandler, middleware)

	return restfulServer
}