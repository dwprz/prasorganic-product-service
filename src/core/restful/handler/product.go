package handler

import (
	"context"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-product-service/src/interface/service"
	"github.com/dwprz/prasorganic-product-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type Product struct {
	productService service.Product
	restfulClient  *client.Restful
}

func NewProduct(ps service.Product, rc *client.Restful) *Product {
	return &Product{
		productService: ps,
		restfulClient:  rc,
	}
}

func (p *Product) Create(c *fiber.Ctx) error {
	req := new(dto.CreateReq)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	uploadRes := c.Locals("upload_imagekit_result").(*uploader.UploadResult)
	req.ImageId = uploadRes.FileId
	req.Image = uploadRes.Url

	err := p.productService.Create(context.Background(), req)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{"data": "successfully created product"})
}
