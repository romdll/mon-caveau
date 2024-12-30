package utils

import (
	"moncaveau/version"
	"os"

	"github.com/gin-gonic/gin"
)

func logServerConfig() {
	logger.Info("Logging all environment variables:")

	for _, envVar := range os.Environ() {
		logger.Infof("%s", envVar)
	}
}

func logServerStartDetails(serverEngine *gin.Engine) {
	logger.Info("START - Logging all the server configuration")

	logger.Infof("Gin Mode: %s", gin.Mode())

	logger.Infof("Server will start listening on: %s", ":80")

	routes := serverEngine.Routes()
	logger.Infof("Total number of routes registered: %d", len(routes))

	logger.Info("Registered routes:")
	for _, route := range routes {
		logger.Infof("Method: %s, Path: %s", route.Method, route.Path)
	}

	logger.Infof("Total number of middleware registered: %d", len(serverEngine.Handlers))

	if !IsDebugMode() {
		logger.Infof("Server version: %s", version.Version)
		logger.Info("Logging environment variables or configuration settings:")
		logServerConfig()
	}

	logger.Info("END - Logging all the server configuration")
}
