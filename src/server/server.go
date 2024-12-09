package server

import (
	"moncaveau/server/frontend"
	"moncaveau/server/handlers"
	"moncaveau/server/middlewares"
	"net/http"

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
	// || Api
	r.POST(ApiLogin, handlers.POST_VerifyAccountLogin)
	r.GET(ApiWinesDashboard, handlers.GET_WinesDashboard)

	// || Frontend
	r.GET(Frontend, frontend.ServeFrontendFiles)

	r.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			c.Redirect(http.StatusTemporaryRedirect, "/v1/")
			return
		}
	})

	return r
}
