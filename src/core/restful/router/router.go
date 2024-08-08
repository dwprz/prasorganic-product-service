package router

import (
	"github.com/dwprz/prasorganic-product-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App, h *handler.Product, m *middleware.Middleware) {
	app.Add("POST", "/api/products", m.VerifyJwt, m.VerifySuperAdmin, m.SaveTemporaryImage, m.ValidateImage, m.UploadToImageKit, h.Create)
}
