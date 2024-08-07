package config

import (
	"context"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

func setUpForNonDevelopment(appStatus string, logger *logrus.Logger) *Config {
	defaultConf := vault.DefaultConfig()
	defaultConf.Address = os.Getenv("PRASORGANIC_CONFIG_ADDRESS")

	client, err := vault.NewClient(defaultConf)
	if err != nil {
		log.Fatalf("vault new client: %v", err)
	}

	client.SetToken(os.Getenv("PRASORGANIC_CONFIG_TOKEN"))

	mountPath := "prasorganic-secrets" + "-" + strings.ToLower(appStatus)

	productServiceSecrets, err := client.KVv2(mountPath).Get(context.Background(), "product-service")
	if err != nil {
		logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	apiGatewaySecrets, err := client.KVv2(mountPath).Get(context.Background(), "api-gateway")
	if err != nil {
		logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddress = productServiceSecrets.Data["RESTFUL_ADDRESS"].(string)
	currentAppConf.GrpcPort = productServiceSecrets.Data["GRPC_PORT"].(string)

	postgresConf := new(postgres)
	postgresConf.Url = productServiceSecrets.Data["POSTGRES_URL"].(string)
	postgresConf.Dsn = productServiceSecrets.Data["POSTGRES_DSN"].(string)
	postgresConf.User = productServiceSecrets.Data["POSTGRES_USER"].(string)
	postgresConf.Password = productServiceSecrets.Data["POSTGRES_PASSWORD"].(string)

	apiGatewayConf := new(apiGateway)
	apiGatewayConf.BaseUrl = apiGatewaySecrets.Data["BASE_URL"].(string)
	apiGatewayConf.BasicAuth = apiGatewaySecrets.Data["BASIC_AUTH"].(string)
	apiGatewayConf.BasicAuthUsername = apiGatewaySecrets.Data["BASIC_AUTH_PASSWORD"].(string)
	apiGatewayConf.BasicAuthPassword = apiGatewaySecrets.Data["BASIC_AUTH_USERNAME"].(string)

	return &Config{
		CurrentApp: currentAppConf,
		Postgres:   postgresConf,
		ApiGateway: apiGatewayConf,
	}
}
