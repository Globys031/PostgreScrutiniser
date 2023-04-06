// All routes are exposed here

package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/web/auth"
	"github.com/Globys031/PostgreScrutiniser/backend/web/file"
	"github.com/Globys031/PostgreScrutiniser/backend/web/resourceConfig"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

// func RegisterRoutes(svc *AuthService) *gin.Engine {
func RegisterRoutes(jwt *auth.JwtWrapper, dbHandler *sql.DB, appUser *utils.User, postgresUser *utils.User, backupDir string, logger *utils.Logger) *gin.Engine {
	////////////////////////
	// Route configurations
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"authorization", "Origin", "Content-Length", "Content-Type"}
	// Default() allows all CORS origins,
	// TO DO: Consider changing this later
	router.Use(cors.New(config))
	////////////////////////

	validate := registerCustomValidators()
	configFilePath, _ := utils.FindConfigFile(dbHandler, logger)

	////////////////////////
	// Register routes
	registerAuthRoute(router, validate, jwt, dbHandler, logger)
	registerResourceConfigRoute(router, validate, jwt, dbHandler, backupDir, postgresUser, appUser, configFilePath, logger)
	registerFileRoute(router, validate, jwt, dbHandler, backupDir, postgresUser, appUser, configFilePath, logger)
	// Registers routes for openapi specification
	registerDocsRoutes(router, logger)

	return router
}

func registerCustomValidators() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation(`backup`, utils.ValidateAutoConfBackup)
	validate.RegisterValidation(`username`, utils.ValidateUsername)
	return validate
}

func registerAuthRoute(router *gin.Engine, validate *validator.Validate, jwt *auth.JwtWrapper, dbHandler *sql.DB, logger *utils.Logger) {
	optionsAuthConfig := &auth.GinServerOptions{
		BaseURL: "/api",
	}
	authConfigApi := &auth.AuthImpl{
		Jwt:      jwt,
		Logger:   logger,
		Validate: validate,
	}
	auth.RegisterHandlersWithOptions(router, authConfigApi, *optionsAuthConfig)
}

func registerResourceConfigRoute(router *gin.Engine, validate *validator.Validate, jwt *auth.JwtWrapper, dbHandler *sql.DB, backupDir string, postgresUser *utils.User, appUser *utils.User, configFilePath string, logger *utils.Logger) {
	optionsResourceConfig := &resourceConfig.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []resourceConfig.MiddlewareFunc{
			resourceConfig.MiddlewareFunc(jwt.ValidateTokenMiddleware(logger)),
		},
	}
	resourceConfigApi := &resourceConfig.ResourceConfigImpl{
		ConfigFile:   configFilePath,
		Logger:       logger,
		AppUser:      appUser,
		PostgresUser: postgresUser,
		DbHandler:    dbHandler,
	}
	resourceConfig.RegisterHandlersWithOptions(router, resourceConfigApi, *optionsResourceConfig)
}

func registerFileRoute(router *gin.Engine, validate *validator.Validate, jwt *auth.JwtWrapper, dbHandler *sql.DB, backupDir string, postgresUser *utils.User, appUser *utils.User, configFilePath string, logger *utils.Logger) {
	optionsFile := &file.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []file.MiddlewareFunc{
			file.MiddlewareFunc(jwt.ValidateTokenMiddleware(logger)),
			// middleware.OapiRequestValidator(swaggerFile),
		},
	}

	fileApi := &file.FileImpl{
		BackupDir:        backupDir,
		CurrentFile:      filepath.Dir(configFilePath) + "/postgresql.auto.conf",
		PostgresUsername: postgresUser.Username,
		AppUser: appUser,
		Logger:           logger,
		DbHandler:        dbHandler,
		Validate:         validate,
	}
	file.RegisterHandlersWithOptions(router, fileApi, *optionsFile)
}

func registerDocsRoutes(router *gin.Engine, logger *utils.Logger) {
	router.GET("/api/docs/auth", func(c *gin.Context) {
		openAPISpecHandler("auth", logger).ServeHTTP(c.Writer, c.Request)
	})
	router.GET("/api/docs/file", func(c *gin.Context) {
		openAPISpecHandler("file", logger).ServeHTTP(c.Writer, c.Request)
	})
	router.GET("/api/docs/resource-config", func(c *gin.Context) {
		openAPISpecHandler("resourceConfig", logger).ServeHTTP(c.Writer, c.Request)
	})
}

// Returns a handler function for displaying openapi documentation
// @docsType - what we are return docs for (file, auth, resourceConfigs, etc...)
func openAPISpecHandler(docsType string, logger *utils.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Select what documentation to load based on docsType
		var swagger *openapi3.T
		var err error = fmt.Errorf("Did not attempt loading any api documentation")
		switch docsType {
		case "auth":
			swagger, err = auth.GetSwagger()
		case "file":
			swagger, err = file.GetSwagger()
		case "resourceConfig":
			swagger, err = resourceConfig.GetSwagger()
		default:
			logger.LogError(fmt.Errorf("Something went wrong loading swagger spec"))
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return docs
		jsonData, err := swagger.MarshalJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonData)
	})
}
