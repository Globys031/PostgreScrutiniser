// All routes are exposed here

package web

import (
	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/web/resourceConfig"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/Globys031/plotzemis/go/auth"
	// "github.com/Globys031/plotzemis/go/db"
)

// type AuthService struct {
// 	Handler db.Handler
// 	Jwt     auth.JwtWrapper
// }

// func RegisterRoutes(svc *AuthService) *gin.Engine {
func RegisterRoutes(logger *utils.Logger) *gin.Engine {
	////////////////////////
	// Initialise logging
	// logger := utils.InitLogging()

	////////////////////////
	// Implementations
	resourceConfigApi := &resourceConfig.ResourceConfigImpl{
		Logger: logger,
		// Configuration: nil,
	}

	////////////////////////
	// Route configurations

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"authorization", "Origin", "Content-Length", "Content-Type"}
	// Default() allows all CORS origins,
	// TO DO: Consider changing this later
	router.Use(cors.New(config))

	// TO DO: add middleware for authentification somewhere here
	optionsResourceConfig := &resourceConfig.GinServerOptions{
		BaseURL: "/api",
		// Middlewares  []MiddlewareFunc
		// ErrorHandler func(*gin.Context, error, int)
	}
	resourceConfig.RegisterHandlersWithOptions(router, resourceConfigApi, *optionsResourceConfig)

	return router
}
