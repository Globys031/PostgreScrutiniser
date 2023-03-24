// All routes are exposed here

package web

import (
	"database/sql"
	"path/filepath"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/web/auth"
	"github.com/Globys031/PostgreScrutiniser/backend/web/file"
	"github.com/Globys031/PostgreScrutiniser/backend/web/resourceConfig"
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
	// Register custom validators
	validate := validator.New()
	validate.RegisterValidation(`backup`, utils.ValidateAutoConfBackup)
	validate.RegisterValidation(`username`, utils.ValidateUsername)

	////////////////////////
	// Register routes

	optionsAuthConfig := &auth.GinServerOptions{
		BaseURL: "/api",
	}
	authConfigApi := &auth.AuthImpl{
		Jwt:      jwt,
		Logger:   logger,
		Validate: validate,
	}
	auth.RegisterHandlersWithOptions(router, authConfigApi, *optionsAuthConfig)

	//////////////////////////////////////////////////////////

	// All routes below use middleware for authentication
	optionsResourceConfig := &resourceConfig.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []resourceConfig.MiddlewareFunc{
			resourceConfig.MiddlewareFunc(jwt.ValidateTokenMiddleware(logger)),
		},
	}

	resourceConfigApi := &resourceConfig.ResourceConfigImpl{
		Logger:       logger,
		AppUser:      appUser,
		PostgresUser: postgresUser,
		DbHandler:    dbHandler,
	}
	resourceConfig.RegisterHandlersWithOptions(router, resourceConfigApi, *optionsResourceConfig)

	//////////////////////////////////////////////////////////

	optionsFile := &file.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []file.MiddlewareFunc{
			file.MiddlewareFunc(jwt.ValidateTokenMiddleware(logger)),
		},
	}

	configFilePath, _ := utils.FindConfigFile(logger)
	fileApi := &file.FileImpl{
		BackupDir:        backupDir,
		CurrentFile:      filepath.Dir(configFilePath) + "/postgresql.auto.conf",
		PostgresUsername: postgresUser.Username,
		Logger:           logger,
		DbHandler:        dbHandler,
		Validate:         validate,
	}
	file.RegisterHandlersWithOptions(router, fileApi, *optionsFile)

	//////////////////////////////////////////////////////////

	return router
}
