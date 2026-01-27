package web

import "github.com/gin-gonic/gin"

type WebApp = *gin.Engine
type WebRouter = gin.IRouter
type WebContext = *gin.Context

// NewWebApp creates a new web application.
//
// It creates a new Gin Engine which can be used as a WebApp.
func NewWebApp() WebApp {
	return gin.New()
}
