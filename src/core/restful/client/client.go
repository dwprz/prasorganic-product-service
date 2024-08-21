package client

import "github.com/dwprz/prasorganic-product-service/src/interface/delivery"

// this main restful client
type Restful struct {
	ImageKit delivery.ImageKitRESTful
}

func NewRestful(ikc delivery.ImageKitRESTful) *Restful {
	return &Restful{
		ImageKit: ikc,
	}
}
