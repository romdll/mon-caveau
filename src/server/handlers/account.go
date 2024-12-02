package handlers

import (
	"moncaveau/database"
	"net/http"

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
		valid, err := database.CheckIfAccountKeyExists(data.AccountKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur interne",
			})
			return
		}

		if valid {
			c.Status(http.StatusOK)

			// TODO set cookie and generate it + save it
			// TODO check if remember me clicked (if so generate a non ending cookie)
			// c.SetCookie()
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Le compte n'est pas valide.",
		})
		return
	}
}
