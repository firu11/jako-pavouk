package handlers

import (
	"backend/utils"

	"github.com/labstack/echo/v4"
)

// vytvoří skupinu /api a v ní všechny endpointy
func SetupRouter(c *echo.Echo) {
	api := c.Group("/api")

	api.GET("/lekce", getVsechnyLekce)
	api.GET("/lekce/:pismena", getCviceniVLekci)
	api.GET("/cvic/:pismena/:cislo", getCviceni)
	api.POST("/dokonceno/:pismena/:cislo", dokoncitCvic)
	api.POST("/dokonceno-procvic/:cislo", dokoncitProcvic)
	api.GET("/procvic", getVsechnyProcvic)
	api.GET("/procvic/:cisloProcvic/:neCislo", getProcvic)
	api.POST("/test-psani", testPsani)

	api.POST("/overit-email", overitEmail)
	api.POST("/registrace", registrace)
	api.POST("/prihlaseni", prihlaseni, utils.RateLimiter)
	api.POST("/zmena-hesla", zmenaHesla)
	api.POST("/overeni-zmeny-hesla", overitZmenuHesla)
	api.POST("/google", google)

	api.GET("/nastaveni", nastaveni)
	api.GET("/statistiky", statistiky)
	api.POST("/ucet-zmena", upravaUctu)

	api.GET("/token-expirace", testVyprseniTokenu)
	api.POST("/navsteva", navsteva)

	setupSkolniRouter(api)
}

func chyba(msg string) echo.Map {
	if msg == "" {
		msg = "Neco se pokazilo."
	}
	return echo.Map{"error": msg}
}
