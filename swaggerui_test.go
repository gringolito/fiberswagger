package fiberswagger

import (
	"fmt"
	"io"
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
	tests := []struct {
		name             string
		input            []Config
		expectedBasePath string
		expectedFilePath string
	}{
		{
			name:             "no args uses both defaults",
			input:            nil,
			expectedBasePath: "/docs",
			expectedFilePath: "./openapi.yaml",
		},
		{
			name:             "empty config uses both defaults",
			input:            []Config{{}},
			expectedBasePath: "/docs",
			expectedFilePath: "./openapi.yaml",
		},
		{
			name:             "only BasePath set falls back FilePath",
			input:            []Config{{BasePath: "/custom"}},
			expectedBasePath: "/custom",
			expectedFilePath: "./openapi.yaml",
		},
		{
			name:             "only FilePath set falls back BasePath",
			input:            []Config{{FilePath: "/spec.yaml"}},
			expectedBasePath: "/docs",
			expectedFilePath: "/spec.yaml",
		},
		{
			name:             "both fields set uses neither default",
			input:            []Config{{BasePath: "/custom", FilePath: "/spec.yaml"}},
			expectedBasePath: "/custom",
			expectedFilePath: "/spec.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := configDefault(tt.input...)
			require.Equal(t, tt.expectedBasePath, cfg.BasePath)
			require.Equal(t, tt.expectedFilePath, cfg.FilePath)
		})
	}
}

func TestConfigDefault_MutationDoesNotAffectDefaults(t *testing.T) {
	original := defaultConfig
	t.Cleanup(func() { defaultConfig = original })

	// Simulate a partial-override call that modifies nothing, then verify
	// that the internal defaults haven't drifted.
	_ = configDefault(Config{BasePath: "/custom"})

	cfg := configDefault()
	require.Equal(t, "/docs", cfg.BasePath)
	require.Equal(t, "./openapi.yaml", cfg.FilePath)
}

func TestRouter_RendersSwaggerUIAndSpec(t *testing.T) {
	specPath := writeSpec(t, testSpec)
	app := fiber.New()

	err := Router(app, Config{BasePath: "/openapi", FilePath: specPath})
	require.NoError(t, err)

	uiResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/openapi/", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, uiResp.StatusCode)

	specResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/openapi/swagger.json", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, specResp.StatusCode)
}

func TestMiddleware_RendersUI(t *testing.T) {
	specPath := writeSpec(t, testSpec)
	app := fiber.New()

	handler, err := Middleware(Config{BasePath: "/docs", FilePath: specPath})
	require.NoError(t, err)
	app.Use(handler)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/docs/", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "swagger-ui")
}

func TestMiddleware_RendersSpec(t *testing.T) {
	specPath := writeSpec(t, testSpec)
	app := fiber.New()

	handler, err := Middleware(Config{BasePath: "/docs", FilePath: specPath})
	require.NoError(t, err)
	app.Use(handler)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/docs/swagger.json", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRouter_ReturnsErrorWhenSpecFileMissing(t *testing.T) {
	app := fiber.New()
	missing := filepath.Join(t.TempDir(), "missing-openapi.yaml")

	err := Router(app, Config{FilePath: missing})
	require.EqualError(t, err, fmt.Sprintf("%s file is not exist", missing))
}

func TestRouter_ReturnsErrorWhenSpecInvalid(t *testing.T) {
	specPath := writeSpec(t, "openapi: [broken")
	app := fiber.New()

	err := Router(app, Config{FilePath: specPath})
	require.Error(t, err)
}

func TestMiddleware_ReturnsErrorWhenSpecFileMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing-openapi.yaml")

	_, err := Middleware(Config{FilePath: missing})
	require.EqualError(t, err, fmt.Sprintf("%s file is not exist", missing))
}

func TestMiddleware_ReturnsErrorWhenSpecInvalid(t *testing.T) {
	specPath := writeSpec(t, "openapi: [broken")

	_, err := Middleware(Config{FilePath: specPath})
	require.Error(t, err)
}

func TestMustRouter_PanicsWhenSpecFileMissing(t *testing.T) {
	app := fiber.New()
	missing := filepath.Join(t.TempDir(), "missing-openapi.yaml")

	require.PanicsWithError(t, fmt.Sprintf("%s file is not exist", missing), func() {
		MustRouter(app, Config{FilePath: missing})
	})
}

func TestMustRouter_PanicsWhenSpecInvalid(t *testing.T) {
	specPath := writeSpec(t, "openapi: [broken")
	app := fiber.New()

	require.Panics(t, func() {
		MustRouter(app, Config{FilePath: specPath})
	})
}

func TestMustMiddleware_PanicsWhenSpecFileMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing-openapi.yaml")

	require.PanicsWithError(t, fmt.Sprintf("%s file is not exist", missing), func() {
		MustMiddleware(Config{FilePath: missing})
	})
}

func TestMustMiddleware_PanicsWhenSpecInvalid(t *testing.T) {
	specPath := writeSpec(t, "openapi: [broken")

	require.Panics(t, func() {
		MustMiddleware(Config{FilePath: specPath})
	})
}
