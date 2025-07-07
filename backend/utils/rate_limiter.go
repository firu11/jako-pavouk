package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var RateLimiter echo.MiddlewareFunc = middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
	Skipper: middleware.DefaultSkipper,
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{Rate: rate.Every(time.Minute / 2), Burst: 10, ExpiresIn: 3 * time.Minute},
	),
	IdentifierExtractor: func(c echo.Context) (string, error) {
		var body struct {
			EmailNeboJmeno string `json:"email"`
		}

		// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
		// workaround aby jsem mohl potom v prihlašování číst znovu
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &body); err != nil || body.EmailNeboJmeno == "" {
			log.Println("dotaz prihlaseni: " + c.RealIP())
			return c.RealIP(), nil
		}

		log.Println("dotaz prihlaseni: " + body.EmailNeboJmeno)
		return body.EmailNeboJmeno, nil
	},
	DenyHandler: func(c echo.Context, identifier string, err error) error {
		return c.NoContent(http.StatusTeapot) //xdd
	},
})
