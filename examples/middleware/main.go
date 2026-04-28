package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	fiberswagger "github.com/gringolito/fiberswagger"
)

func main() {
	app := fiber.New()

	// Middleware mode wraps the swagger handlers in a single fiber.Handler
	// registered via app.Use. Requests that do not match a swagger route
	// (/docs/ or /docs/swagger.json) are propagated to the next handler,
	// so downstream Fiber middleware continues to run for non-swagger traffic.
	app.Use(fiberswagger.MustMiddleware(fiberswagger.Config{
		BasePath: "/docs",
		FilePath: "../openapi.yaml",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Visit /docs for Swagger UI.")
	})

	log.Fatal(app.Listen(":3000"))
}
