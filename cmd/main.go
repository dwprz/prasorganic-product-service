package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-product-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/delivery"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/imagekit"
	"github.com/dwprz/prasorganic-product-service/src/repository"
	"github.com/dwprz/prasorganic-product-service/src/service"
	"github.com/go-playground/validator/v10"
)

func handleCloseApp(closeCH chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		close(closeCH)
	}()
}

func main() {
	closeCH := make(chan struct{})
	handleCloseApp(closeCH)

	validate := validator.New()
	postgresDB := database.NewPostgres()

	imageKit := imagekit.New()

	productRepository := repository.NewProduct(postgresDB)
	productService := service.NewProduct(validate, productRepository)
	imageKitDelivery := delivery.NewImageKit(imageKit)
	restfulClient := client.New(imageKitDelivery)
	productRestfulHandler := handler.NewProduct(productService, restfulClient)
	middleware := middleware.New(restfulClient)
	restfuleServer := server.New(productRestfulHandler, middleware)
	defer restfuleServer.Stop()

	go restfuleServer.Run()

	<-closeCH
}
