// All routes are exposed here

package web

import (
	"database/sql"
	"path/filepath"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/web/file"
	"github.com/Globys031/PostgreScrutiniser/backend/web/resourceConfig"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

// type AuthService struct {
// 	Handler db.Handler
// 	Jwt     auth.JwtWrapper
// }

// func RegisterRoutes(svc *AuthService) *gin.Engine {
func RegisterRoutes(dbHandler *sql.DB, appUser *utils.User, postgresUser *utils.User, backupDir string, logger *utils.Logger) *gin.Engine {
	////////////////////////
	// Route configurations
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"authorization", "Origin", "Content-Length", "Content-Type"}
	// Default() allows all CORS origins,
	// TO DO: Consider changing this later
	router.Use(cors.New(config))

	// ////////////////////////
	// // Register custom validators
	validate := validator.New()
	validate.RegisterValidation(`backup`, utils.ValidateAutoConfBackup)

	////////////////////////
	// Register routes

	// TO DO: add middleware for authentification somewhere here
	optionsResourceConfig := &resourceConfig.GinServerOptions{
		BaseURL: "/api",
		// Middlewares  []MiddlewareFunc
		// ErrorHandler func(*gin.Context, error, int)
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
