package server

import (
	"moncaveau/server/frontend"
	"moncaveau/server/handlers"
	"moncaveau/server/handlers/debug"
	"moncaveau/server/middlewares"
	"moncaveau/utils"
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
	// || Auth
	r.POST(ApiLogin, handlers.POST_VerifyAccountLogin)
	r.GET(ApiLogout, handlers.GET_Logout)
	r.GET(ApiRegister, handlers.GET_GenerateAccount)
	r.GET(ApiAccountDetails, handlers.GET_AccountDetails)

	// || Compiled data
	r.GET(ApiWinesDashboard, handlers.GET_WinesDashboard)

	// || Wines sub
	r.GET(ApiWinesFetchRegionsCountries, handlers.GET_GetCountriesAndRegions)
	r.GET(ApiWinesFetchTypes, handlers.GET_GetWineTypes)
	r.GET(ApiWinesFetchBottleSizes, handlers.GET_GetBottleSizes)
	r.GET(ApiWinesFetchDomains, handlers.GET_GetDomains)
	r.GET(ApiWinesFetchPaginated, handlers.GET_allWines)
	r.GET(ApiWinesFetchTransactions, handlers.GET_FetchAllTransactionsForCharts)
	r.GET(ApiWinesFetchStatistics, handlers.GET_wineStatisticsData)

	// || Wines
	r.POST(ApiWinesWineCreation, handlers.POST_CreateWine)

	// || Frontend
	r.GET(Frontend, frontend.ServeFrontendFiles)

	r.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			c.Redirect(http.StatusTemporaryRedirect, "/v1/")
			return
		}
	})
	r.GET(Favicon, func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, RealFavicon)
	})

	// Debug handlers
	if utils.IsDebugMode() {
		r.GET(DebugSqlStructure, debug.GET_DatabaseStructure)
		r.GET(DebugCreateFakeAccount, debug.GET_CreateFakeAccount)
	}

	return r
}
