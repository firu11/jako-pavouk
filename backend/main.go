package main

import (
	"backend/databaze"
	"backend/utils"
	"log"
	"math"
	"os"
	"regexp"
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

const tokenTimeDuration time.Duration = time.Hour * 24 * 15 // v nanosekundach, 14 + 1 dni asi good (den predem uz odhlasime aby se nestalo ze neco dela a neulozi se to)
var regexJmeno *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$`)

const cifraCislaZaJmenem int = 4

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Nenašel jsem soubor .env v /backend.")
	}

	pocetSlov, pocetPismenVeSlovu = getEnvDelky()
	delkaTextu = (pocetPismenVeSlovu+1)*pocetSlov - 1

	databaze.DBConnect()
	inject()

	app := fiber.New(fiber.Config{
		AppName: "Pavouk",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://jakopavouk.cz/",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:               30,
		Expiration:        15 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
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

func inject() {
	utils.TokenTimeDuration = tokenTimeDuration
	utils.RegexJmeno = regexJmeno
	utils.CifraCislaZaJmenem = cifraCislaZaJmenem
	utils.MaxCislo = int(math.Pow(10, float64(cifraCislaZaJmenem))) // 10_000
}

// abych pro testing měl kratší texty tak to mám v .env
func getEnvDelky() (int, int) {
	x, err := strconv.Atoi(os.Getenv("POCET_SLOV"))
	if err != nil {
		log.Fatalln("ENV se pokazilo - asi spatna hodnota", err)
	}
	y, err := strconv.Atoi(os.Getenv("POCET_PISMEN"))
	if err != nil {
		log.Fatalln("ENV se pokazilo - asi spatna hodnota", err)
	}
	return x, y
}
