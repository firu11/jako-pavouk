package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/firu11/jako-pavouk/backend/config"
	"github.com/firu11/jako-pavouk/backend/databaze"
	"github.com/firu11/jako-pavouk/backend/handlers"
	"github.com/firu11/jako-pavouk/backend/middlewares"
	"github.com/firu11/jako-pavouk/backend/utils"

	"github.com/labstack/echo/v5"
)

func main() {
	cfg := config.LoadEnvConfig()

	utils.SetupEmaily(cfg.Email)
	utils.SetupAuth(cfg.Auth)
	utils.SetupGoogle(cfg.Auth)
	databaze.Init(cfg.DatabaseURL)
	defer databaze.Close()

	e := echo.New()

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         addr,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: 10 * time.Second,
	}
	if err := sc.Start(ctx, e); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
