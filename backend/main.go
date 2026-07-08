package main

import (
	"log"

	"github.com/firu11/jako-pavouk/backend/config"
	"github.com/firu11/jako-pavouk/backend/databaze"
	"github.com/firu11/jako-pavouk/backend/handlers"
	"github.com/firu11/jako-pavouk/backend/middlewares"
	"github.com/firu11/jako-pavouk/backend/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.LoadEnvConfig()

	utils.SetupEmaily(cfg.Email)
	databaze.Init(cfg.DatabaseURL)
	defer databaze.Close()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	publicDir, err := handlers.ResolvePublicDir(cfg.PublicDir)
	if err != nil {
		if cfg.Production {
			log.Fatal(err)
		}
		log.Println(err) // thats okay in dev
	} else {
		handlers.RegisterStaticRoutes(e, publicDir)
	}

	middlewares.RegisterBasic(e)
	handlers.SetupRouter(e)

	addr := cfg.Address()
	log.Printf("starting server on %s...", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
