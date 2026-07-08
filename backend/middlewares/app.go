package middlewares

import (
	"net/http"
	"strings"

	"github.com/firu11/jako-pavouk/backend/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterBasic(e *echo.Echo) {
	e.Use(AuthContext())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health" || strings.HasPrefix(c.Path(), "/api")
		},
		Level: 5,
	}))
	e.GET("/health", Health)
}

func AuthContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("uzivID", utils.Autentizace(c.Request().Header.Get("Authorization")))
			return next(c)
		}
	}
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}
