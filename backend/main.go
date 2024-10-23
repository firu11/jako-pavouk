package main

import (
	"backend/databaze"
	"backend/utils"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var pocetZnaku float32 = 1500
var pocetPismenVeSlovu int = 4

var regexJmeno *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$`)

const tokenTimeDuration time.Duration = time.Hour * 24 * 15 // v nanosekundach, 14 + 1 dni asi good (den predem uz odhlasime aby se nestalo ze neco dela a neulozi se to)
const cifraCislaZaJmenem int = 4

var MaxCisloZaJmeno int = int(math.Pow(10, float64(cifraCislaZaJmenem))) // 10_000

var rateLimiter echo.MiddlewareFunc = middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
	Skipper: middleware.DefaultSkipper,
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{Rate: rate.Every(5 * time.Minute / 5), Burst: 5, ExpiresIn: 5 * time.Minute},
	),
	IdentifierExtractor: func(c echo.Context) (string, error) {
		var body bodyPrihlaseni

		// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
		// workaround aby jsem mohl potom v prihlašování číst znovu
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return c.RealIP(), nil
		}
		return body.EmailNeboJmeno, nil
	},
	DenyHandler: func(c echo.Context, identifier string, err error) error {
		return c.NoContent(http.StatusTeapot) //xdd
	},
})

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Nenašel jsem soubor .env v /backend.")
	}

	databaze.DBConnect()
	inject()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "https://jakopavouk.cz/"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("uzivID", utils.Autentizace(c.Request().Header.Get("Authorization")))
			return next(c)
		}
	})

	SetupRouter(e)

	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func inject() {
	utils.TokenTimeDuration = tokenTimeDuration
	databaze.RegexJmeno = regexJmeno
	databaze.MaxCisloZaJmeno = MaxCisloZaJmeno
}
