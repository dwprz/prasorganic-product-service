package router

import (
	"github.com/dwprz/prasorganic-product-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-product-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func Create(app *fiber.App, h *handler.ProductRESTful, m *middleware.Middleware) {
	// super admin
	app.Add("POST", "/api/products", m.VerifyJwt, m.VerifySuperAdmin, m.SaveTemporaryImage, m.ValidateImage, m.UploadToImageKit, h.Create)
	app.Add("PATCH", "/api/products/:productId", m.VerifyJwt, m.VerifySuperAdmin, h.Update)
	app.Add("PATCH", "/api/products/:productId/image", m.VerifyJwt, m.VerifySuperAdmin, m.SaveTemporaryImage, m.ValidateImage, m.UploadToImageKit, h.UpdateImage)

	// all
	app.Add("GET", "/api/products", h.Get)
}