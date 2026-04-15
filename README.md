# fiberswagger

Fiber middleware and router helpers for exposing Swagger UI and OpenAPI spec files.

## Install

```bash
go get github.com/gringolito/fiberswagger
```

## Usage

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gringolito/fiberswagger"
)

app := fiber.New()

fiberswagger.Router(app, fiberswagger.Config{
    BasePath: "/openapi",
    FilePath: "./api/spec/openapi.yaml",
})
```

## Features

- Serves Swagger UI under a configurable base path
- Serves JSON OpenAPI spec at `<basePath>/swagger.json`
- Supports direct router mounting or middleware mode
