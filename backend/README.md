# PostgreScrutiniser Backend

This is the backend code for `PostgreScrutiniser`. It follows a basic folder structure [as described here](https://github.com/golang-standards/project-layout) and additionally has a `utils` folder for functions used across the application.

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
oapi-codegen -generate "types" -package types api/example.yaml > web/types/example.go
oapi-codegen -generate "gin" -package routes api/example.yaml > web/routes/example.go
oapi-codegen -generate "spec" -package specs api/example.yaml > web/specs/example.go
```
Change `example` to what code is actually being generated for.

#### TO DO:

Use this for declaring proper imports:
```
x-go-type-import
```