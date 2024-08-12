package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-product-service/src/core/grpc"
	"github.com/dwprz/prasorganic-product-service/src/core/restful"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/delivery"
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
	imageKitDelivery := delivery.NewImageKit(imageKit)

	productRepository := repository.NewProduct(postgresDB)
	productService := service.NewProduct(validate, productRepository)

	restfulServer := restful.Initialize(productService, imageKitDelivery)
	defer restfulServer.Stop()

	go restfulServer.Run()

	grpcServer := grpc.Initialize(productService)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-closeCH
}
