package frontend

import (
	"embed"
	"mime"
	"moncaveau/server/middlewares"
	"moncaveau/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

//go:embed "site/*"
var siteData embed.FS

func ServeFrontendFiles(c *gin.Context) {
	filePath := c.Param("filepath")[1:]
	logger.Printf("Received request for frontend file: '%s'\n", filePath)

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

func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "text/plain"
	}
	return contentType
}

func getHtmlContext(c *gin.Context) map[string]interface{} {
	return map[string]interface{}{
		"LoggedIn": c.GetBool(middlewares.ContextIsLoggedIn),
	}
}

func serveFromFileSystem(c *gin.Context, filePath string) {
	file := filepath.Join("server", "frontend", "site", filePath)
	logger.Printf("Serving file from the filesystem: %s", file)

	if _, err := os.Stat(file); err != nil {
		logger.Printf("File not found: %s", file)
		c.Status(http.StatusNotFound)
		return
	}

	logger.Printf("File found: '%s'", file)
	contentType := getContentType(filePath)
	logger.Printf("Guessed content type for file: '%s' is '%s'", filePath, contentType)

	if strings.Contains(contentType, "text/html") {
		data, err := os.ReadFile(file)
		if err != nil {
			logger.Printf("Error reading HTML file: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("index").Parse(string(data))
		if err != nil {
			logger.Printf("Error parsing HTML template: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", contentType)
		err = tmpl.Execute(c.Writer, getHtmlContext(c))
		if err != nil {
			logger.Printf("Error executing template: %v", err)
			c.Status(http.StatusInternalServerError)
		}
	} else {
		c.Header("Content-Type", contentType)
		c.File(file)
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

	logger.Printf("File found from embedded filesystem: '%s'", filePath)
	contentType := getContentType(filePath)
	logger.Printf("Guessed content type for file: '%s' is '%s'", filePath, contentType)

	if strings.Contains(contentType, "text/html") {
		tmpl, err := template.New("index").Parse(string(data))
		if err != nil {
			logger.Printf("Error parsing HTML template: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", contentType)
		err = tmpl.Execute(c.Writer, getHtmlContext(c))
		if err != nil {
			logger.Printf("Error executing template: %v", err)
			c.Status(http.StatusInternalServerError)
		}
	} else {
		c.Header("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, data)
	}
}
