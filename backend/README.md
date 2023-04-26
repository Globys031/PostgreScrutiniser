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
NAME=resourceConfig
oapi-codegen -generate "types" -package ${NAME} api/${NAME}.yaml > web/${NAME}/${NAME}Types.go
oapi-codegen -generate "gin" -package ${NAME} api/${NAME}.yaml > web/${NAME}/${NAME}Gin.go
oapi-codegen -generate "spec" -package ${NAME} api/${NAME}.yaml > web/${NAME}/${NAME}Spec.go
```
Change `example` to what code is actually being generated for.

### References

Below is a list of references used for creating the backend side of this application
- https://www.postgresql.org/docs/
- https://github.com/jberkus/annotated.conf
- https://github.com/jfcoz/postgresqltuner/blob/master/postgresqltuner.pl
- https://github.com/deepmap/oapi-codegen#overview