package main

import (
	"backend/databaze"
	"backend/handlers"
	"backend/utils"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var pocetZnaku float32 = 1500
var pocetPismenVeSlovu int = 4

var regexJmeno *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$`)

const tokenTimeDuration time.Duration = time.Hour * 24 * 15 // v nanosekundach, 14 + 1 dni asi good (den predem uz odhlasime aby se nestalo ze neco dela a neulozi se to)
const cifraCislaZaJmenem int = 4

var MaxCisloZaJmeno int = int(math.Pow(10, float64(cifraCislaZaJmenem))) // 10_000

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

	handlers.SetupRouter(e)

	err = utils.SetupEmaily()
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start("127.0.0.1:1323")) // 1324 na testing
}

func inject() {
	utils.TokenTimeDuration = tokenTimeDuration
	databaze.RegexJmeno = regexJmeno
	databaze.MaxCisloZaJmeno = MaxCisloZaJmeno
	handlers.PocetZnaku = pocetZnaku
	handlers.RegexJmeno = regexJmeno
	utils.PocetPismenVeSlovu = pocetPismenVeSlovu
}
