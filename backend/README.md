# PostgreScrutiniser Backend

This is the backend code for `PostgreScrutiniser`.

## Starting project

### Backend

```
go run ./backend/
```

## Development

Below are details concerning development of this project

### Generating openapi code

Below command uses [oapi-codegen](https://github.com/deepmap/oapi-codegen#overview) to generate API code for Gin framework:
```
oapi-codegen -generate "types" -package types docs/example.yaml > http/types/example.go
oapi-codegen -generate "gin" -package routes docs/example.yaml > http/routes/example.go
oapi-codegen -generate "spec" -package specs docs/example.yaml > http/specs/example.go
```
Change `example` to what code is actually being generated for.

#### TO DO:

Use this for declaring proper imports:
```
x-go-type-import
```