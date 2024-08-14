package client

import "github.com/dwprz/prasorganic-product-service/src/interface/delivery"

// this main restful client
type Restful struct {
	ImageKit delivery.ImageKit
}

func NewRestful(ikc delivery.ImageKit) *Restful {
	return &Restful{
		ImageKit: ikc,
	}
}
