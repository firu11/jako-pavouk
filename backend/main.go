package main

import (
	"backend/databaze"
	"backend/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

var pocetSlov int = 20 // jeste zvednu
var pocetPismenVeSlovu int = 4
var delkaTextu = (pocetPismenVeSlovu+1)*pocetSlov - 1
var tokenTimeDuration time.Duration = time.Hour * 24 * 15 // v nanosekundach, 14 + 1 dni asi good (den predem uz odhlasime aby se nestalo ze neco dela a neulozi se to)

func main() {
	databaze.DBConnect()
	inject()

	app := fiber.New(fiber.Config{
		AppName: "Pavouk",
	})

	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://jakopavouk.cz",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	SetupRouter(app)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	} // http://localhost:8080
}

func inject() {
	utils.TokenTimeDuration = tokenTimeDuration
}
