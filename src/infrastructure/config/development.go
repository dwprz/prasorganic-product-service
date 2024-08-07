package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setUpForDevelopment(logger *logrus.Logger) *Config {
	viper := viper.New()

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logger.WithFields(logrus.Fields{"location": "config.setUpForDevelopment", "section": "viper.ReadInConfig"}).Fatal(err)
	}

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddress = viper.GetString("CURRENT_APP_RESTFUL_ADDRESS")
	currentAppConf.GrpcPort = viper.GetString("CURRENT_APP_GRPC_PORT")

	postgresConf := new(postgres)
	postgresConf.Url = viper.GetString("POSTGRES_URL")
	postgresConf.Dsn = viper.GetString("POSTGRES_DSN")
	postgresConf.User = viper.GetString("POSTGRES_USER")
	postgresConf.Password = viper.GetString("POSTGRES_PASSWORD")

	apiGatewayConf := new(apiGateway)
	apiGatewayConf.BaseUrl = viper.GetString("API_GATEWAY_BASE_URL")
	apiGatewayConf.BasicAuth = viper.GetString("API_GATEWAY_BASIC_AUTH")
	apiGatewayConf.BasicAuthUsername = viper.GetString("API_GATEWAY_BASIC_AUTH_USERNAME")
	apiGatewayConf.BasicAuthPassword = viper.GetString("API_GATEWAY_BASIC_AUTH_PASSWORD")

	return &Config{
		CurrentApp: currentAppConf,
		Postgres:   postgresConf,
		ApiGateway: apiGatewayConf,
	}
}
