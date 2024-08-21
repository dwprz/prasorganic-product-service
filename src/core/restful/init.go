package restful

import (
	"github.com/dwprz/prasorganic-product-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/delivery"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-product-service/src/interface/service"
)

func InitServer(ps service.Product, rc *client.Restful) *server.Restful {
	productHandler := handler.NewProductRESTful(ps, rc)

	middleware := middleware.New(rc)
	restfulServer := server.New(productHandler, middleware)

	return restfulServer
}

func InitClient() *client.Restful {
	imageKitDelivery := delivery.NewImageKit()
	restfulClient := client.NewRestful(imageKitDelivery)

	return restfulClient
}
