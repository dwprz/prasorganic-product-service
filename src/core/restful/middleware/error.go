package middleware

import (
	"context"
	"encoding/json"

	errcustom "github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/common/errors/restful"
	"github.com/dwprz/prasorganic-product-service/src/common/helper"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func (m *Middleware) Error(c *fiber.Ctx, err error) error {
	restful.LogError(c, err)

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
		return restful.HandleValidationError(c, validationError)
	}

	if responseError, ok := err.(*errcustom.Response); ok {
		return restful.HandleResponseError(c, responseError)
	}

	if jwtError := restful.HanldeJwtError(err); jwtError != nil {
		return c.Status(401).JSON(fiber.Map{"errors": jwtError.Error()})
	}

	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		return restful.HandleJsonError(c, jsonError)
	}

	return c.Status(500).JSON(fiber.Map{
		"errors": "sorry, internal server error try again later",
	})
}
