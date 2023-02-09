// All routes are exposed here

package routes

import (
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
func RegisterRoutes() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"authorization", "Origin", "Content-Length", "Content-Type"}
	// Default() allows all CORS origins,
	// TO DO: Consider changing this later
	router.Use(cors.New(config))

	// Each route function decides where to use authentication middleware (if resource should be protected)
	// inside the function itself

	// TO DO: add middleware for authentification here
	// routes := router.Group("/api")

	// TO DO: come back and see how routes could be properly registered if using openapi generation

	return router
}
