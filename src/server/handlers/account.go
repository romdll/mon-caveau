package handlers

import (
	"moncaveau/database"
	"moncaveau/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func POST_VerifyAccountLogin(c *gin.Context) {
	data := database.Account{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Le contenu de la requete est invalide.",
		})
		return
	}

	if data.AccountKey == "" && data.Email == "" && data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La clé de compte, l'email et le mot de passe sont vides. Veuillez en saisir au moins 1.",
		})
		return
	}

	// Login with only account key
	if data.AccountKey != "" && data.Email == "" && data.Password == "" {
		valid, accountId, err := database.CheckIfAccountKeyExists(data.AccountKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur interne",
			})
			return
		}

		if valid {
			sessionToken, err := utils.GenerateSessionToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erreur interne impossible de créer une nouvelle session.",
				})
				return
			}

			// TODO check if remember me clicked (if so generate a non ending cookie)
			// TODO setup a real expiration instead of 200 * days
			expirationDate := time.Now().Add(200 * (24 * time.Hour))

			newSession := database.Session{
				AccountID:    accountId,
				SessionToken: sessionToken,
				CreatedAt:    time.Now(),
				ExpiresAt:    expirationDate,
				LastActivity: time.Now(),
			}

			_, err = database.InsertEntityById(newSession)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erreur interne impossible de sauvegarder votre session.",
				})
				return
			}

			c.SetCookie(
				database.AuthCookieName,
				sessionToken,
				int(time.Until(expirationDate).Seconds()),
				"/",
				"",
				database.IsCookieSecure,
				true,
			)
			c.JSON(http.StatusOK, gin.H{
				"token": sessionToken,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Le compte n'est pas valide.",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Type de connection non implémentée.",
	})
}
