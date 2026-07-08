package middlewares

import (
	"net/http"

	"github.com/firu11/jako-pavouk/backend/config"
	"github.com/firu11/jako-pavouk/backend/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Register(e *echo.Echo, cfg *config.Config) {
	if !cfg.Production {
		e.Use(devCORS())
	}

	e.Use(AuthContext())
	e.GET("/healthz", Healthz)
}

func devCORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	})
}

func AuthContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("uzivID", utils.Autentizace(c.Request().Header.Get("Authorization")))
			return next(c)
		}
	}
}

func Healthz(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
