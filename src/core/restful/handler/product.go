package handler

import (
	"context"
	"strconv"

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
	req := new(dto.CreateProductReq)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	uploadRes := c.Locals("upload_imagekit_result").(*uploader.UploadResult)
	req.ImageId = uploadRes.FileId
	req.Image = uploadRes.Url

	err := p.productService.Create(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{"data": "successfully created product"})
}

func (p *Product) Get(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return err
	}

	category := c.Query("category")
	productName := c.Query("name")

	res, err := p.productService.FindMany(c.Context(), &dto.GetProductReq{
		Page:        page,
		Category:    category,
		ProductName: productName,
	})

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": res.Data, "paging": res.Paging})
}

func (p *Product) Update(c *fiber.Ctx) error {
	req := new(dto.UpdateProductReq)

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return err
	}

	req.ProductId = uint(productId)

	res, err := p.productService.Update(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": res})
}

func (p *Product) UpdateImage(c *fiber.Ctx) error {
	req := new(dto.UpdateProductImageReq)

	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return err
	}

	req.ProductId = uint(productId)

	uploadRes := c.Locals("upload_imagekit_result").(*uploader.UploadResult)
	req.ImageId = uploadRes.FileId
	req.Image = uploadRes.Url

	res, err := p.productService.UpdateImage(c.Context(), req)
	if err != nil {
		return err
	}

	oldImageId := c.FormValue("image_id")
	go p.restfulClient.ImageKit.DeleteFile(context.Background(), oldImageId)

	return c.Status(200).JSON(fiber.Map{"data": res})
}
