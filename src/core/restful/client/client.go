package client

import "github.com/dwprz/prasorganic-product-service/src/interface/delivery"

// this main restful client
type Restful struct {
	ImageKit delivery.ImageKitRestful
}

func NewRestful(ikc delivery.ImageKitRestful) *Restful {
	return &Restful{
		ImageKit: ikc,
	}
}
