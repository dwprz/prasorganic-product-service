package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-product-service/src/core/grpc"
	"github.com/dwprz/prasorganic-product-service/src/core/restful"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-product-service/src/repository"
	"github.com/dwprz/prasorganic-product-service/src/service"
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

	postgresDB := database.NewPostgres()
	defer database.ClosePostgres(postgresDB)

	productRepository := repository.NewProduct(postgresDB)
	productService := service.NewProduct(productRepository)

	restfulClient := restful.InitClient()
	restfulServer := restful.InitServer(productService, restfulClient)
	defer restfulServer.Stop()

	go restfulServer.Run()

	grpcServer := grpc.InitServer(productService)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-closeCH
}
