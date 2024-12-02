package server

import (
	"moncaveau/server/frontend"
	"moncaveau/server/handlers"
	"moncaveau/server/middlewares"

	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Attach middlewares
	r.Use(middlewares.Logger())
	r.Use(gin.Recovery())

	// Attach handlers
	r.POST("/api/login", handlers.POST_VerifyAccountLogin)

	// || Frontend
	r.GET("/*filepath", frontend.ServeFrontendFiles)

	return r
}
