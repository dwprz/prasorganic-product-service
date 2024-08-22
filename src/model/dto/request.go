package dto

type CreateProductReq struct {
	ProductName string  `json:"product_name" validate:"required,min=3,max=100"`
	ImageId     string  `json:"image_id" validate:"required,min=10,max=100"`
	Image       string  `json:"image" validate:"required,min=10,max=500"`
	Price       uint    `json:"price" validate:"required"`
	Stock       uint    `json:"stock" validate:"required"`
	Category    string  `json:"category" validate:"required,min=3,max=20"`
	Length      uint8   `json:"length" validate:"required"`
	Width       uint8   `json:"width" validate:"required"`
	Height      uint8   `json:"height" validate:"required"`
	Weight      float32 `json:"weight" validate:"required"`
	Description string  `json:"description" validate:"required"`
}

type GetProductsReq struct {
	Page        int    `json:"page" validate:"min=1,max=100" gorm:"-"`
	ProductName string `json:"product_name" validate:"omitempty,min=3,max=100"`
	Category    string `json:"category" validate:"omitempty,min=3,max=20"`
}

type UpdateProductReq struct {
	ProductId   uint    `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"omitempty,min=3,max=100"`
	Rating      float32 `json:"rating" validate:"omitempty"`
	Sold        uint32  `json:"sold" validate:"omitempty"`
	Price       uint    `json:"price" validate:"omitempty"`
	Stock       uint    `json:"stock" validate:"omitempty"`
	Category    string  `json:"category" validate:"omitempty,min=3,max=20"`
	Length      uint8   `json:"length" validate:"omitempty"`
	Width       uint8   `json:"width" validate:"omitempty"`
	Height      uint8   `json:"height" validate:"omitempty"`
	Weight      float32 `json:"weight" validate:"omitempty"`
	Description string  `json:"description" validate:"omitempty"`
}

type UpdateImagePoductReq struct {
	ProductId uint   `json:"product_id" validate:"required"`
	ImageId   string `json:"image_id" validate:"required,min=10,max=100"`
	Image     string `json:"image" validate:"required,min=10,max=500"`
}

type ReduceStocksReq struct {
	ProductId uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type RollbackStoksReq struct {
	ProductId uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}
