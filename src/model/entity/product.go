package entity

import (
	"time"
)

type Product struct {
	ProductId   uint      `json:"product_id" gorm:"column:product_id;primaryKey;autoIncrement;<-:create"`
	ProductName string    `json:"product_name" gorm:"column:product_name"`
	ImageId     string    `json:"image_id" gorm:"column:image_id"`
	Image       string    `json:"image" gorm:"column:image"`
	Rating      float32   `json:"rating" gorm:"column:rating;default:null"`
	Sold        uint32    `json:"sold" gorm:"column:sold;default:null"`
	Price       uint      `json:"price" gorm:"column:price"`
	Stock       uint      `json:"stock" gorm:"column:stock"`
	Category    string    `json:"category" gorm:"column:category"`
	Length      uint8     `json:"length" gorm:"column:length"`
	Width       uint8     `json:"width" gorm:"column:width"`
	Height      uint8     `json:"height" gorm:"column:height"`
	Weight      float32   `json:"weight" gorm:"column:weight"`
	Description string    `json:"description" gorm:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt   time.Time `json:"updated_at" grom:"column:updated_at;autoUpdateTime"`
}

func (p *Product) TableName() string {
	return "products"
}
