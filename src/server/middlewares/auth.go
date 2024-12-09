package middlewares

import (
	"moncaveau/database"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextIsLoggedIn     = "UserIsLoggedIn"
	ContextLoggedInUserId = "UserLoggedInUserId"
)

func AuthApi(baseUrls []string, toAvoid []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUrl := c.Request.URL.String()

		isInBaseUrls := false
		for _, baseUrl := range baseUrls {
			if strings.HasPrefix(requestUrl, baseUrl) {
				isInBaseUrls = true
			}
		}

		if !isInBaseUrls {
			logger.Printf("Skipping hard authentication check, URL does not match baseUrl: %s\n", requestUrl)

			sessionToken, err := c.Cookie(database.AuthCookieName)
			if err == nil && sessionToken != "" {
				validSessionToken, userID, err := database.VerifyIfSessionExistsAndIsValid(sessionToken)
				if err == nil && validSessionToken {
					logger.Printf("Authenticated successfully for request (no hard check): %s", requestUrl)
					c.Set(ContextIsLoggedIn, true)
					c.Set(ContextLoggedInUserId, userID)
				}
			}

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

		validSessionToken, userID, err := database.VerifyIfSessionExistsAndIsValid(sessionToken)
		if err != nil {
			logger.Printf("Error during session token verification: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Une erreur est survenue lors de la vérification de votre authentification.",
			})
			return
		}

		if !validSessionToken {
			logger.Printf("Invalid session token for request: %s\n", requestUrl)

			c.SetCookie(database.AuthCookieName, "", -1, "/", "", false, true)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous n'etes pas autorisé à accéder à cette page",
			})
			return
		}

		logger.Printf("Authenticated successfully for request: %s", requestUrl)
		c.Set(ContextIsLoggedIn, true)
		c.Set(ContextLoggedInUserId, userID)
		c.Next()
	}
}
