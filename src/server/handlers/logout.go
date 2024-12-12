package handlers

import (
	"moncaveau/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET_Logout(c *gin.Context) {
	sessionToken, err := c.Cookie(database.AuthCookieName)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = database.DeleteSessionToken(sessionToken)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.SetCookie(
		database.AuthCookieName,
		"",
		-1,
		"/",
		"",
		database.IsCookieSecure,
		true,
	)
	c.Status(http.StatusOK)
}
