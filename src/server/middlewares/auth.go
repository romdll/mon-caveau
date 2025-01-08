package middlewares

import (
	"moncaveau/database"
	"moncaveau/server/middlewares/activity"
	"moncaveau/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextIsLoggedIn     = "UserIsLoggedIn"
	ContextLoggedInUserId = "UserLoggedInUserId"
)

func AuthApi(baseUrls []string, toAvoid []string, activityBuffer *activity.ActivityDoubleBuffer) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUrl := c.Request.URL.String()

		isInBaseUrls := false
		for _, baseUrl := range baseUrls {
			if strings.HasPrefix(requestUrl, baseUrl) {
				isInBaseUrls = true
				break
			}
		}

		if !isInBaseUrls {
			logger.Debugf("Skipping hard authentication check, URL does not match baseUrl: %s", requestUrl)

			sessionToken, err := c.Cookie(database.AuthCookieName)
			if err == nil && sessionToken != "" {
				validSessionToken, userID, err := database.VerifyIfSessionExistsAndIsValid(sessionToken)
				if err == nil && validSessionToken {
					logger.Debugf("Authenticated successfully for request (no hard check): %s", requestUrl)
					activityBuffer.Write(sessionToken)

					c.Set(ContextIsLoggedIn, true)
					c.Set(ContextLoggedInUserId, userID)
				}
			}

			c.Next()
			return
		}

		for _, avoid := range toAvoid {
			if strings.Contains(requestUrl, avoid) {
				logger.Debugf("Skipping authentication check, URL contains allowed segment: %s", requestUrl)
				c.Next()
				return
			}
		}

		sessionToken, err := c.Cookie(database.AuthCookieName)
		if err != nil || sessionToken == "" {
			logger.Warnf("Unauthorized access attempt, no session token: %s", requestUrl)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous n'êtes pas autorisé à accéder à cette page",
			})
			return
		}

		validSessionToken, userID, err := database.VerifyIfSessionExistsAndIsValid(sessionToken)
		if err != nil {
			logger.Errorf("Error during session token verification: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Une erreur est survenue lors de la vérification de votre authentification.",
			})
			return
		}

		if !validSessionToken {
			logger.Warnf("Invalid session token for request: %s", requestUrl)

			c.SetCookie(database.AuthCookieName, "", -1, "/", "", utils.IsHttps(), true)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous n'êtes pas autorisé à accéder à cette page",
			})
			return
		}

		logger.Infof("Authenticated successfully for request: %s", requestUrl)
		activityBuffer.Write(sessionToken)

		c.Set(ContextIsLoggedIn, true)
		c.Set(ContextLoggedInUserId, userID)
		c.Next()
	}
}
