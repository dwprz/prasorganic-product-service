package middleware

import (
	"fmt"

	errcustom "github.com/dwprz/prasorganic-product-service/src/common/errors"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) VerifyJwt(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")

	if accessToken == "" {
		return &errcustom.Response{HttpCode: 401, Message: "token is required"}
	}

	jwtToken, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected token method: %v", t.Header["alg"])
		}

		return config.Conf.Jwt.PublicKey, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		c.Locals("user_data", claims)
		return c.Next()
	}

	return &errcustom.Response{HttpCode: 401, Message: "token is invalid"}
}
