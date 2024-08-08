package imagekit

import (
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/config"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/logger"
)

func New() *imagekit.ImageKit {
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  config.Conf.ImageKit.PrivateKey,
		PublicKey:   config.Conf.ImageKit.PublicKey,
		UrlEndpoint: config.Conf.ImageKit.BaseUrl,
	})

	ik.Logger.SetLevel(logger.ERROR)

	return ik
}
