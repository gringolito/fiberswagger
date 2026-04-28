package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	fiberswagger "github.com/gringolito/fiberswagger"
)

func main() {
	app := fiber.New()

	// Router mode mounts two explicit GET routes under BasePath:
	//   GET /docs/           → Swagger UI
	//   GET /docs/swagger.json → OpenAPI spec
	//
	// These routes are terminal: Fiber middleware registered after MustRouter
	// will not run for requests matched by them.
	fiberswagger.MustRouter(app, fiberswagger.Config{
		BasePath: "/docs",
		FilePath: "../openapi.yaml",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Visit /docs for Swagger UI.")
	})

	log.Fatal(app.Listen(":3000"))
}
