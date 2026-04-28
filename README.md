# fiberswagger

[![CI](https://github.com/gringolito/fiberswagger/actions/workflows/ci.yaml/badge.svg)](https://github.com/gringolito/fiberswagger/actions/workflows/ci.yaml)
[![CodeQL](https://github.com/gringolito/fiberswagger/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/gringolito/fiberswagger/actions/workflows/github-code-scanning/codeql)
[![Go Reference](https://pkg.go.dev/badge/github.com/gringolito/fiberswagger.svg)](https://pkg.go.dev/github.com/gringolito/fiberswagger)
![Go 1.25+](https://img.shields.io/badge/go-%3E%3D1.25-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Fiber middleware and router helpers for exposing Swagger UI and OpenAPI spec files.

## Install

```bash
go get github.com/gringolito/fiberswagger
```

## Features

- Serves Swagger UI under a configurable base path
- Serves JSON OpenAPI spec at `<basePath>/swagger.json`
- Supports direct router mounting or middleware mode

## Usage

Both functions accept `Config` as an optional argument. When omitted, the UI is served at `/docs/` and the spec at `/docs/swagger.json`, reading the spec from `./openapi.yaml`:

```go
fiberswagger.MustRouter(app)
// or
app.Use(fiberswagger.MustMiddleware())
```

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

## Choosing between Middleware and Router

Both modes serve the same UI and spec, but they differ in how they interact with the Fiber handler chain.

**Router** mounts two explicit routes (`GET /` and `GET /swagger.json` under `BasePath`) using `adaptor.HTTPHandler` with a `nil` next handler. The route handlers are terminal: Fiber middleware registered *after* `Router()` will not execute for requests that match the swagger routes. Choose `Router` when you want clean, isolated route registration and don't need downstream middleware to run on swagger traffic.

**Middleware** registers a single `fiber.Handler` via `app.Use` using `adaptor.HTTPMiddleware`, which propagates the Fiber handler chain through the underlying `net/http` layer. Requests that do not match the swagger paths (UI or spec) are passed to the next Fiber handler, so subsequent middleware continues to run. Choose `Middleware` when swagger should sit inside a broader middleware stack and non-swagger requests must fall through to other handlers.

| | Router | Middleware |
| --- | --- | --- |
| Mount style | Explicit routes (`app.Get`) | `app.Use` catch-all |
| Next-handler propagation | No — swagger routes are terminal | Yes — non-swagger requests fall through |
| Downstream middleware runs | No | Yes (for non-swagger paths) |

## Config Reference

Both `Router` and `Middleware` accept an optional `Config` argument. All fields are optional and fall back to the defaults shown below.

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `BasePath` | `string` | `"/docs"` | URL path prefix where Swagger UI and the spec are served. The UI is at `<BasePath>/` and the spec at `<BasePath>/swagger.json`. |
| `FilePath` | `string` | `"./openapi.yaml"` | Filesystem path to the OpenAPI spec file. Relative paths resolve from the working directory at startup. |

## Contributing

Contributions are welcome! Feel free to:

- [Open an issue](https://github.com/gringolito/fiberswagger/issues) to report a bug or request a feature
- Submit a pull request — please include tests for any new behavior
- Suggest ideas by starting a [discussion](https://github.com/gringolito/fiberswagger/discussions)

---

*This project has moved from the Beerware License to MIT, but the spirit lives on: if we ever meet and you think this stuff is worth it, you're still very welcome to buy me a beer.*
