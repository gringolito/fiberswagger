package fiberswagger

import (
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

const testSpec = `openapi: 3.0.0
info:
  title: test
  version: 1.0.0
paths: {}
`

func writeSpec(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "openapi.yaml")
	err := os.WriteFile(path, []byte(content), 0o600)
	require.NoError(t, err)
	return path
}

func TestConfigDefault(t *testing.T) {
	cfg := configDefault()
	require.Equal(t, "/docs", cfg.BasePath)
	require.Equal(t, "./openapi.yaml", cfg.FilePath)
}

func TestRouter_RendersSwaggerUIAndSpec(t *testing.T) {
	specPath := writeSpec(t, testSpec)
	app := fiber.New()

	Router(app, Config{BasePath: "/openapi", FilePath: specPath})

	uiResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/openapi/", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, uiResp.StatusCode)

	specResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/openapi/swagger.json", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, specResp.StatusCode)
}

func TestMiddleware_RendersSpec(t *testing.T) {
	specPath := writeSpec(t, testSpec)
	app := fiber.New()
	app.Use(Middleware(Config{BasePath: "/docs", FilePath: specPath}))

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/docs/swagger.json", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRouter_PanicsWhenSpecFileMissing(t *testing.T) {
	app := fiber.New()
	missing := filepath.Join(t.TempDir(), "missing-openapi.yaml")

	require.PanicsWithError(t, fmt.Sprintf("%s file is not exist", missing), func() {
		Router(app, Config{FilePath: missing})
	})
}

func TestRouter_PanicsWhenSpecInvalid(t *testing.T) {
	specPath := writeSpec(t, "openapi: [broken")
	app := fiber.New()

	require.Panics(t, func() {
		Router(app, Config{FilePath: specPath})
	})
}
