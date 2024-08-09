package dto

type CreateProductReq struct {
	ProductName string `json:"product_name" validate:"required,min=3,max=100"`
	ImageId     string `json:"image_id" validate:"required,min=10,max=100"`
	Image       string `json:"image" validate:"required,min=10,max=500"`
	Price       uint   `json:"price" validate:"required"`
	Stock       uint   `json:"stock" validate:"required"`
	Category    string `json:"category" validate:"required,min=3,max=20"`
	Length      uint8  `json:"length" validate:"required"`
	Width       uint8  `json:"width" validate:"required"`
	Height      uint8  `json:"height" validate:"required"`
	Weight      uint   `json:"weight" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type GetProductReq struct {
	Page        int    `json:"page" validate:"min=1,max=100" gorm:"-"`
	ProductName string `json:"product_name" validate:"omitempty,min=3,max=100"`
	Category    string `json:"category" validate:"omitempty,min=3,max=20"`
}

