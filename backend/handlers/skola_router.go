package handlers

import "github.com/labstack/echo/v4"

// typy uživatelů
// 1 - basic
// 2 - učitel
func setupSkolniRouter(api *echo.Group) {
	skolaApi := api.Group("/skola")

	skolaApi.POST("/create-trida", createTrida)
	skolaApi.GET("/tridy", tridy)
	skolaApi.GET("/trida", tridaStudent)
	skolaApi.GET("/trida/:id", trida)
	skolaApi.GET("/zaci-stream/:id", zaciStream)
	skolaApi.GET("/test-tridy/:kod", testTridy)
	skolaApi.POST("/zmena-tridy", zmenaTridy)

	skolaApi.POST("/pridat-praci", pridatPraci)
	skolaApi.GET("/get-praci/:id", getPraci)
	skolaApi.GET("/get-statistiky-prace/:id", getStatistikyPrace)
	skolaApi.POST("/dokoncit-praci/:id", dokoncitPraci)
	skolaApi.DELETE("/smazat-praci/:id", smazatPraci)

	skolaApi.POST("/text", getText)
	skolaApi.GET("/typy-cviceni", getTypyCviceni)

	skolaApi.GET("/student/:id", student)
	skolaApi.POST("/student", studentUprava)
	skolaApi.POST("/zapis", zapis)

	skolaApi.GET("/ucitele", ucitele)
	skolaApi.POST("/upravit-ucitele", upravaUcitele)
	skolaApi.POST("/zapis-skoly", zapisSkoly)
}
