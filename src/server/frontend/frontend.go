package frontend

import (
	"embed"
	"fmt"
	"moncaveau/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed "site/*"
var siteData embed.FS

func ServeFrontendFiles(c *gin.Context) {
	filePath := c.Param("filepath")[1:]
	logger.Printf("Received request for frontend file: '%s'", filePath)

	if !strings.Contains(filePath, ".") {
		logger.Printf("File path does not have any suffix, adding .html extension.")
		filePath += ".html"
	}

	if filePath == ".html" {
		logger.Printf("File path is just '/' translating to index.")
		filePath = "index.html"
	}

	if utils.IsDebugMode() {
		logger.Println("Debug mode is enabled, serving from file system.")
		serveFromFileSystem(c, filePath)
	} else {
		logger.Println("Serving from embedded filesystem.")
		serveFromEmbeddedFS(c, filePath)
	}
}

func serveFromFileSystem(c *gin.Context, filePath string) {
	file := fmt.Sprintf("server/frontend/site/%s", filePath)
	logger.Printf("Serving file from the filesystem: %s", file)

	if _, err := os.Stat(file); err == nil {
		logger.Printf("File found: %s", file)
		c.File(file)
	} else {
		logger.Printf("File not found: %s", file)
		c.Status(http.StatusNotFound)
	}
}

func serveFromEmbeddedFS(c *gin.Context, filePath string) {
	logger.Printf("Serving file from embedded filesystem: %s", filePath)

	data, err := siteData.ReadFile("site/" + filePath)
	if err != nil {
		logger.Printf("Error reading file from embedded filesystem: %s, error: %v", filePath, err)
		c.Status(http.StatusNotFound)
		return
	}

	logger.Printf("File successfully served from embedded filesystem: %s", filePath)
	c.Data(http.StatusOK, "text/html", data)
}
