package main

import (
	"backend/databaze"
	"backend/utils"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

var pocetSlov int
var pocetPismenVeSlovu int
var delkaTextu int
var tokenTimeDuration time.Duration = time.Hour * 24 * 15 // v nanosekundach, 14 + 1 dni asi good (den predem uz odhlasime aby se nestalo ze neco dela a neulozi se to)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error .env", err)
	}

	pocetSlov = getPocetSlov()
	pocetPismenVeSlovu = getDelkaPismen()
	delkaTextu = (pocetPismenVeSlovu+1)*pocetSlov - 1

	databaze.DBConnect()
	inject()

	app := fiber.New(fiber.Config{
		AppName: "Pavouk",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://jakopavouk.cz",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:               30,
		Expiration:        10 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTeapot) // troulin
		},
	}))

	SetupRouter(app)

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	} // http://localhost:8080
}

func inject() {
	utils.TokenTimeDuration = tokenTimeDuration
}

func getDelkaPismen() int {
	x, err := strconv.Atoi(os.Getenv("POCET_PISMEN"))
	if err != nil {
		log.Fatalln("ENV se pokazilo - asi spatna hodnota", err)
	}
	return x
}

func getPocetSlov() int {
	x, err := strconv.Atoi(os.Getenv("POCET_SLOV"))
	if err != nil {
		log.Fatalln("ENV se pokazilo - asi spatna hodnota", err)
	}
	return x
}
