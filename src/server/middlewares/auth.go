package middlewares

import (
	"moncaveau/database"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthApi(baseUrl string, toAvoid []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUrl := c.Request.URL.String()
		if !strings.HasPrefix(requestUrl, baseUrl) {
			logger.Printf("Skipping authentication check, URL does not match baseUrl: %s\n", requestUrl)
			c.Next()
			return
		}

		for _, avoid := range toAvoid {
			if strings.Contains(requestUrl, avoid) {
				logger.Printf("Skipping authentication check, URL contains allowed segment: %s\n", requestUrl)
				c.Next()
				return
			}
		}

		sessionToken, err := c.Cookie(database.AuthCookieName)
		if err != nil || sessionToken == "" {
			logger.Printf("Unauthorized access attempt, no session token: %s\n", requestUrl)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous n'etes pas autorisé à accéder à cette page",
			})
			return
		}

		validSessionToken, err := database.VerifyIfSessionExistsAndIsValid(sessionToken)
		if err != nil {
			logger.Printf("Error during session token verification: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Une erreur est survenue lors de la vérification de votre authentification.",
			})
			return
		}

		if !validSessionToken {
			logger.Printf("Invalid session token for request: %s\n", requestUrl)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous n'etes pas autorisé à accéder à cette page",
			})
			return
		}

		logger.Printf("Authenticated successfully for request: %s", requestUrl)
		c.Next()
	}
}
