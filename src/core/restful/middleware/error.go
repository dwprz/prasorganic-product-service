package middleware

import (
	"context"

	"github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/common/helper"
	"github.com/dwprz/prasorganic-product-service/src/common/log"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"github.com/sirupsen/logrus"
)

func (m *Middleware) Error(c *fiber.Ctx, err error) error {
	log.Logger.WithFields(logrus.Fields{
		"host":     c.Hostname(),
		"ip":       c.IP(),
		"protocol": c.Protocol(),
		"location": c.OriginalURL(),
		"method":   c.Method(),
		"from":     "error middleware",
	}).Error(err.Error())

	if c.OriginalURL() == "/api/products" && c.Method() == "POST" {
		
		filename, ok := c.Locals("filename").(string)
		if ok && filename != "" {
			go helper.DeleteFile("./tmp/" + filename)
		}

		req, ok := c.Locals("upload_imagekit_result").(*uploader.UploadResult)
		if ok && req.FileId != "" {
			go m.restfulClient.ImageKit.DeleteFile(context.Background(), req.FileId)
		}
	}

	if validationError, ok := err.(validator.ValidationErrors); ok {

		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]any{
				"field":       validationError[0].Field(),
				"description": validationError[0].Error(),
			},
		})
	}

	if responseError, ok := err.(*errors.Response); ok {
		return c.Status(int(responseError.HttpCode)).JSON(fiber.Map{
			"errors": responseError.Message,
		})
	}

	return c.Status(500).JSON(fiber.Map{
		"errors": "sorry, internal server error try again later",
	})
}
