package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterStaticRoutes(e *echo.Echo, publicDir string) {
	indexPath := filepath.Join(publicDir, "index.html")

	e.GET("/*", func(c echo.Context) error {
		requestPath := c.Request().URL.Path
		if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		cleanPath := strings.TrimPrefix(filepath.Clean(requestPath), "/")
		if cleanPath == "." || cleanPath == "" {
			setIndexCacheHeaders(c)
			return c.File(indexPath)
		}

		fullPath := filepath.Join(publicDir, cleanPath)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			setStaticCacheHeaders(c, cleanPath)
			return c.File(fullPath)
		}

		if filepath.Ext(cleanPath) != "" {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		setIndexCacheHeaders(c)
		return c.File(indexPath)
	})
}

func setIndexCacheHeaders(c echo.Context) {
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache, no-store, must-revalidate")
}

func setStaticCacheHeaders(c echo.Context, path string) {
	if isImmutableAsset(path) {
		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000, immutable")
		return
	}

	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache, must-revalidate")
}

func isImmutableAsset(path string) bool {
	if strings.HasPrefix(path, "assets/") {
		return true
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".avif", ".svg", ".ico", ".ttf", ".woff", ".woff2", ".mp4":
		return true
	default:
		return false
	}
}

func ResolvePublicDir(dir string) (string, error) {
	indexPath := filepath.Join(dir, "index.html")

	info, err := os.Stat(indexPath)
	if err != nil || info.IsDir() {
		return "", fmt.Errorf("frontend build not found in %q", dir)
	}

	return dir, nil
}
