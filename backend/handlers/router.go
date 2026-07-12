package handlers

import (
	"net/http"
	"time"

	"github.com/firu11/jako-pavouk/backend/config"
	"github.com/firu11/jako-pavouk/backend/middlewares"
	"github.com/firu11/jako-pavouk/backend/utils"

	"github.com/labstack/echo/v5"
)

var secureAuthCookie bool

// vytvoří skupinu /api a v ní všechny endpointy
func SetupRouter(c *echo.Echo, production bool) {
	secureAuthCookie = production
	api := c.Group("/api")

	api.GET("/lekce", getVsechnyLekce)
	api.GET("/lekce/:pismena", getCviceniVLekci)
	api.GET("/cvic/:pismena/:cislo", getCviceni)
	api.POST("/dokonceno/:pismena/:cislo", dokoncitCvic)
	api.POST("/dokonceno-procvic/:cislo", dokoncitProcvic)
	api.GET("/procvic", getVsechnyProcvic)
	api.GET("/procvic/:cisloProcvic/:neCislo", getProcvic)
	api.POST("/test-psani", testPsani)
	api.POST("/uloz-procvic-postup", ulozProcvicPostup)
	api.GET("/procvic-postup/:cislo", getProcvicPostup)

	api.POST("/overit-email", overitEmail)
	api.POST("/registrace", registrace)
	api.POST("/prihlaseni", prihlaseni, middlewares.RateLimiter)
	api.POST("/zmena-hesla", zmenaHesla)
	api.POST("/overeni-zmeny-hesla", overitZmenuHesla)
	api.POST("/google", google)
	api.POST("/odhlaseni", odhlaseni)

	api.GET("/nastaveni", nastaveni)
	api.GET("/statistiky", statistiky)
	api.POST("/ucet-zmena", upravaUctu)

	api.GET("/token-expirace", testVyprseniTokenu)

	setupSkolniRouter(api)
}

func setAuthCookie(c *echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     utils.AuthCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(config.TokenLifetime / time.Second),
		HttpOnly: true,
		Secure:   secureAuthCookie,
		SameSite: http.SameSiteLaxMode,
	})
}

func odhlaseni(c *echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     utils.AuthCookieName,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureAuthCookie,
		SameSite: http.SameSiteLaxMode,
	})
	return c.NoContent(http.StatusNoContent)
}

func chyba(msg string) map[string]any {
	if msg == "" {
		msg = "Neco se pokazilo."
	}
	return map[string]any{"error": msg}
}
