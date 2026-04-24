# fiberswagger

Fiber middleware and router helpers for exposing Swagger UI and OpenAPI spec files.

## Install

```bash
go get github.com/gringolito/fiberswagger
```

## Usage

### Router mode

`Router` mounts dedicated routes for the UI and spec. It returns an error if the spec cannot be loaded, so you decide how to handle misconfiguration:

```go
import (
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/gringolito/fiberswagger"
)

app := fiber.New()

if err := fiberswagger.Router(app, fiberswagger.Config{
    BasePath: "/openapi",
    FilePath: "./api/spec/openapi.yaml",
}); err != nil {
    log.Fatalf("failed to mount swagger: %v", err)
}
```

If you prefer a panic on misconfiguration (e.g. during `init` or in tests), use the `Must` variant:

```go
fiberswagger.MustRouter(app, fiberswagger.Config{
    BasePath: "/openapi",
    FilePath: "./api/spec/openapi.yaml",
})
```

### Middleware mode

`Middleware` returns a `fiber.Handler` that can be used with `app.Use`. It returns an error for the same reasons as `Router`:

```go
handler, err := fiberswagger.Middleware(fiberswagger.Config{
    BasePath: "/docs",
    FilePath: "./api/spec/openapi.yaml",
})
if err != nil {
    log.Fatalf("failed to build swagger middleware: %v", err)
}
app.Use(handler)
```

Or with the `Must` variant:

```go
app.Use(fiberswagger.MustMiddleware(fiberswagger.Config{
    BasePath: "/docs",
    FilePath: "./api/spec/openapi.yaml",
}))
```

## Features

- Serves Swagger UI under a configurable base path
- Serves JSON OpenAPI spec at `<basePath>/swagger.json`
- Supports direct router mounting or middleware mode
