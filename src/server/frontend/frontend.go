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
	logger.Infow("Received request for frontend file", "filePath", filePath)

	if !strings.Contains(filePath, ".") {
		logger.Infow("File path does not have any suffix, adding .html extension.", "filePath", filePath)
		filePath += ".html"
	}

	if filePath == ".html" {
		logger.Infow("File path is just '/', translating to index.", "filePath", filePath)
		filePath = "index.html"
	}

	if utils.IsDebugMode() {
		logger.Infow("Debug mode is enabled, serving from file system.")
		serveFromFileSystem(c, filePath)
	} else {
		logger.Infow("Serving from embedded filesystem.")
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
	baseDir := filepath.Join("server", "frontend", "site")
	cleanedPath := filepath.Clean(filePath)
	file := filepath.Join(baseDir, cleanedPath)

	if !strings.HasPrefix(file, baseDir) {
		logger.Warnw("Attempt to access file outside base directory", "requestedPath", filePath, "resolvedPath", file)
		c.Status(http.StatusForbidden)
		return
	}

	logger.Infow("Serving file from the filesystem", "file", file)

	if _, err := os.Stat(file); err != nil {
		logger.Debugw("File not found", "file", file, "error", err)
		c.Status(http.StatusNotFound)
		return
	}

	logger.Infow("File found", "file", file)
	contentType := getContentType(filePath)
	logger.Infow("Guessed content type for file", "filePath", filePath, "contentType", contentType)

	// This is not necessary as we already prevent
	// G304 (CWE-22): Potential file inclusion via variable (Confidence: HIGH, Severity: MEDIUM)
	// with a clean above and a verification
	// but its here to remove the error from gosec verifications
	file = filepath.Clean(file)

	if strings.Contains(contentType, "text/html") {
		data, err := os.ReadFile(file)
		if err != nil {
			logger.Errorw("Error reading HTML file", "file", file, "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("index").Parse(string(data))
		if err != nil {
			logger.Errorw("Error parsing HTML template", "file", file, "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", contentType)
		err = tmpl.Execute(c.Writer, getHtmlContext(c))
		if err != nil {
			logger.Errorw("Error executing template", "file", file, "error", err)
			c.Status(http.StatusInternalServerError)
		}
	} else {
		c.Header("Content-Type", contentType)
		c.File(file)
	}
}

func serveFromEmbeddedFS(c *gin.Context, filePath string) {
	logger.Infow("Serving file from embedded filesystem", "filePath", filePath)

	data, err := siteData.ReadFile("site/" + filePath)
	if err != nil {
		logger.Debugw("Error reading file from embedded filesystem", "filePath", filePath, "error", err)
		c.Status(http.StatusNotFound)
		return
	}

	logger.Infow("File found from embedded filesystem", "filePath", filePath)
	contentType := getContentType(filePath)
	logger.Infow("Guessed content type for file", "filePath", filePath, "contentType", contentType)

	if strings.Contains(contentType, "text/html") {
		tmpl, err := template.New("index").Parse(string(data))
		if err != nil {
			logger.Errorw("Error parsing HTML template", "filePath", filePath, "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", contentType)
		err = tmpl.Execute(c.Writer, getHtmlContext(c))
		if err != nil {
			logger.Errorw("Error executing template", "filePath", filePath, "error", err)
			c.Status(http.StatusInternalServerError)
		}
	} else {
		c.Header("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, data)
	}
}
