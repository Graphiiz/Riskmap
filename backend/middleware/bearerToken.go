package bearerToken

import (
	"github.com/labstack/echo/v4"
)

func BearerToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		c.Request().Header.Set("Authorization", token)
		return next(c)
	}
}
