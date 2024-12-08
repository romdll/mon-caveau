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
	r.Use(middlewares.AuthApi(AuthProtectedPages, AuthAvoidPages))

	// Attach handlers
	r.POST(ApiLogin, handlers.POST_VerifyAccountLogin)

	// || Frontend
	r.GET(Frontend, frontend.ServeFrontendFiles)

	return r
}
