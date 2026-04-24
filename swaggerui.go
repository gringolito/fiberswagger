package fiberswagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func loadSpec(cfg Config) (*loads.Document, error) {
	if _, err := os.Stat(cfg.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s file is not exist", cfg.FilePath)
	}

	spec, err := loads.Spec(cfg.FilePath)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// Middleware returns a fiber.Handler that renders an OpenAPI specification using SwaggerUI.
// It returns an error if the spec file cannot be loaded or parsed.
// Use MustMiddleware if you prefer a panic on misconfiguration.
func Middleware(config ...Config) (fiber.Handler, error) {
	cfg := configDefault(config...)

	spec, err := loadSpec(cfg)
	if err != nil {
		return nil, err
	}

	specJSON, err := json.Marshal(spec.Raw())
	if err != nil {
		return nil, err
	}

	return adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		swaggerUIHandler := middleware.SwaggerUI(middleware.SwaggerUIOpts{
			Path:    cfg.BasePath,
			SpecURL: path.Join(cfg.BasePath, "swagger.json"),
		}, next)

		return middleware.Spec(cfg.BasePath, specJSON, swaggerUIHandler)
	}), nil
}

// MustMiddleware is like Middleware but panics on error.
// Use this when misconfiguration should be a fatal startup error.
func MustMiddleware(config ...Config) fiber.Handler {
	handler, err := Middleware(config...)
	if err != nil {
		panic(err)
	}
	return handler
}

// Router creates routes with handlers that render an OpenAPI specification using SwaggerUI.
// It returns an error if the spec file cannot be loaded or parsed.
// Use MustRouter if you prefer a panic on misconfiguration.
func Router(router fiber.Router, config ...Config) error {
	cfg := configDefault(config...)

	spec, err := loadSpec(cfg)
	if err != nil {
		return err
	}

	router.Route(cfg.BasePath, func(router fiber.Router) {
		router.Get("/", handleSwaggerUI(cfg)).Name("ui")
		router.Get("/swagger.json", handleSwaggerJSON(spec.Raw())).Name("spec")
	}, "swagger.")

	return nil
}

// MustRouter is like Router but panics on error.
// Use this when misconfiguration should be a fatal startup error.
func MustRouter(router fiber.Router, config ...Config) {
	if err := Router(router, config...); err != nil {
		panic(err)
	}
}

func handleSwaggerUI(cfg Config) fiber.Handler {
	return adaptor.HTTPHandler(middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    cfg.BasePath,
		SpecURL: path.Join(cfg.BasePath, "swagger.json"),
	}, nil))
}

func handleSwaggerJSON(swagger json.RawMessage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Status(http.StatusOK).JSON(swagger)
		return nil
	}
}
