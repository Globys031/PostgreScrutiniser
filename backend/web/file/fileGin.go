// Package file provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package file

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (DELETE /backup)
	DeleteBackups(c *gin.Context)

	// (GET /backup)
	GetBackups(c *gin.Context)

	// (DELETE /backup/{backup_name})
	DeleteBackup(c *gin.Context, backupName string)

	// (PUT /backup/{backup_name})
	PutBackup(c *gin.Context, backupName string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// DeleteBackups operation middleware
func (siw *ServerInterfaceWrapper) DeleteBackups(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteBackups(c)
}

// GetBackups operation middleware
func (siw *ServerInterfaceWrapper) GetBackups(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetBackups(c)
}

// DeleteBackup operation middleware
func (siw *ServerInterfaceWrapper) DeleteBackup(c *gin.Context) {

	var err error

	// ------------- Path parameter "backup_name" -------------
	var backupName string

	err = runtime.BindStyledParameter("simple", false, "backup_name", c.Param("backup_name"), &backupName)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter backup_name: %s", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteBackup(c, backupName)
}

// PutBackup operation middleware
func (siw *ServerInterfaceWrapper) PutBackup(c *gin.Context) {

	var err error

	// ------------- Path parameter "backup_name" -------------
	var backupName string

	err = runtime.BindStyledParameter("simple", false, "backup_name", c.Param("backup_name"), &backupName)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter backup_name: %s", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PutBackup(c, backupName)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {

	errorHandler := options.ErrorHandler

	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.DELETE(options.BaseURL+"/backup", wrapper.DeleteBackups)

	router.GET(options.BaseURL+"/backup", wrapper.GetBackups)

	router.DELETE(options.BaseURL+"/backup/:backup_name", wrapper.DeleteBackup)

	router.PUT(options.BaseURL+"/backup/:backup_name", wrapper.PutBackup)

	return router
}
