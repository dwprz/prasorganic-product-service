package dto

type UpImageRes struct {
	ImageId string `json:"image_id" validate:"required,min=10,max=100"`
	Image   string `json:"image" validate:"required,min=10,max=500"`
}

type CreateReq struct {
	ProductName string `json:"product_name" validate:"required,min=3,max=100"`
	ImageId     string `json:"image_id" validate:"required,min=10,max=100"`
	Image       string `json:"image" validate:"required,min=10,max=500"`
	Price       uint   `json:"price" validate:"required"`
	Stock       uint   `json:"stock" validate:"required"`
	Length      uint8  `json:"length" validate:"required"`
	Width       uint8  `json:"width" validate:"required"`
	Height      uint8  `json:"height" validate:"required"`
	Weight      uint   `json:"weight" validate:"required"`
	Description string `json:"description" validate:"required"`
}
