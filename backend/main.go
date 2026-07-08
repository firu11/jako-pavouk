package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/firu11/jako-pavouk/backend/config"
	"github.com/firu11/jako-pavouk/backend/databaze"
	"github.com/firu11/jako-pavouk/backend/handlers"
	"github.com/firu11/jako-pavouk/backend/middlewares"
	"github.com/firu11/jako-pavouk/backend/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.LoadEnvConfig()

	databaze.DBConnect(cfg)

	if err := utils.SetupEmaily(cfg.Email); err != nil {
		log.Fatal(err)
	}

	publicDir, err := resolvePublicDir(cfg.PublicDir)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	middlewares.Register(e, cfg)
	handlers.SetupRouter(e)
	registerStaticRoutes(e, publicDir)

	addr := cfg.Address()
	log.Printf("starting server on %s!", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func registerStaticRoutes(e *echo.Echo, publicDir string) {
	indexPath := filepath.Join(publicDir, "index.html")
	fileServer := http.FileServer(http.Dir(publicDir))

	e.GET("/*", func(c echo.Context) error {
		requestPath := c.Request().URL.Path
		if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		cleanPath := strings.TrimPrefix(filepath.Clean(requestPath), "/")
		if cleanPath == "." || cleanPath == "" {
			return c.File(indexPath)
		}

		fullPath := filepath.Join(publicDir, cleanPath)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(c.Response(), c.Request())
			return nil
		}

		if filepath.Ext(cleanPath) != "" {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		return c.File(indexPath)
	})
}

func resolvePublicDir(configuredDir string) (string, error) {
	candidates := []string{}
	if configuredDir != "" {
		candidates = append(candidates, configuredDir)
	}
	candidates = append(candidates, "public", "../frontend/dist")

	for _, dir := range candidates {
		indexPath := filepath.Join(dir, "index.html")
		if info, err := os.Stat(indexPath); err == nil && !info.IsDir() {
			return dir, nil
		}
	}

	return "", fmt.Errorf("no frontend build found; looked for index.html in: %s", strings.Join(candidates, ", "))
}
