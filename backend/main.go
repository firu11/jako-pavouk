package main

import (
	"backend/databaze"
	"backend/utils"
	"log"
	"math"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

var pocetZnaku float32 = 1500
var pocetPismenVeSlovu int = 4

var regexJmeno *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$`)

const tokenTimeDuration time.Duration = time.Hour * 24 * 15 // v nanosekundach, 14 + 1 dni asi good (den predem uz odhlasime aby se nestalo ze neco dela a neulozi se to)
const cifraCislaZaJmenem int = 4

var MaxCisloZaJmeno int = int(math.Pow(10, float64(cifraCislaZaJmenem))) // 10_000

// # Main - začátek programu
//  1. načte .env proměnné
//  2. připojí se k PostgreSQL databázi
//  3. nakonfiguruje middleware
//  4. spustí server
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Nenašel jsem soubor .env v /backend.")
	}

	databaze.DBConnect()
	inject()

	app := fiber.New(fiber.Config{
		AppName: "Jako Pavouk Backend",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://jakopavouk.cz/",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use("/api/prihlaseni", limiter.New(limiter.Config{
		Max:               10,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			var body bodyPrihlaseni
			if err := c.BodyParser(&body); err != nil {
				return c.IP()
			}
			return body.EmailNeboJmeno
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTeapot) // troulin
		},
	}))

	SetupRouter(app)

	err = app.Listen("127.0.0.1:44871")
	if err != nil {
		log.Fatal(err)
	}
}

// předá proměnné které chci mít všechny na jednom místě do dalších souborů
func inject() {
	utils.TokenTimeDuration = tokenTimeDuration
	databaze.RegexJmeno = regexJmeno
	databaze.MaxCisloZaJmeno = MaxCisloZaJmeno
}
