package handlers

import (
	"moncaveau/database"
	"moncaveau/database/crypt"
	"moncaveau/server/middlewares"
	"moncaveau/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func POST_VerifyAccountLogin(c *gin.Context) {
	data := database.Account{}
	if err := c.BindJSON(&data); err != nil {
		logger.Errorf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Le contenu de la requete est invalide.",
		})
		return
	}

	if data.AccountKey == "" && data.Email == "" && data.Password == "" {
		logger.Warn("Invalid input: all fields (AccountKey, Email, Password) are empty.")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La clé de compte, l'email et le mot de passe sont vides. Veuillez en saisir au moins 1.",
		})
		return
	}

	if data.AccountKey != "" && data.Email == "" && data.Password == "" {
		logger.Infof("Attempting login with AccountKey: %s", data.AccountKey)

		valid, accountId, err := database.CheckIfAccountKeyExists(data.AccountKey)
		if err != nil {
			logger.Errorf("Error checking if AccountKey exists: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur interne",
			})
			return
		}

		if valid {
			logger.Infof("AccountKey valid. Generating session token for Account ID: %d", accountId)

			sessionToken, err := utils.GenerateSessionToken()
			if err != nil {
				logger.Errorf("Error generating session token: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erreur interne impossible de créer une nouvelle session.",
				})
				return
			}

			rememberMe := data.RememberMe
			var expirationDate time.Time
			if rememberMe {
				logger.Info("User selected 'Remember Me'. Creating long-lasting session.")
				expirationDate = time.Now().Add(365 * 24 * time.Hour)
			} else {
				expirationDate = time.Now().Add(10 * 24 * time.Hour)
			}

			newSession := database.Session{
				AccountID:    accountId,
				SessionToken: sessionToken,
				CreatedAt:    time.Now(),
				ExpiresAt:    expirationDate,
				LastActivity: time.Now(),
			}

			_, err = database.InsertEntityById(newSession)
			if err != nil {
				logger.Errorf("Error saving session: %v", err)
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
				utils.IsHttps(),
				true,
			)
			logger.Infof("Session created and cookie set for Account ID: %d", accountId)
			c.JSON(http.StatusOK, gin.H{
				"token": sessionToken,
			})
			return
		}

		logger.Warnf("Invalid account for AccountKey: %s", data.AccountKey)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Le compte n'est pas valide.",
		})
		return
	}

	logger.Warn("Unsupported login type or combination.")
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Type de connection non implémentée.",
	})
}

func GET_GenerateAccount(c *gin.Context) {
	logger.Info("Starting account generation process.")

	accountKey, err := crypt.GenerateSecureAccountKey()
	if err != nil {
		logger.Errorf("Error generating secure account key: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de générer un nouveau compte.",
		})
		return
	}

	logger.Infof("Generated new secure account key: %s", utils.MaskOnlyNumbers(accountKey, 6))

	if _, err := database.InsertEntityById(&database.Account{AccountKey: accountKey}); err != nil {
		logger.Errorf("Error inserting new account into the database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de créer un nouveau compte dans le système.",
		})
		return
	}

	logger.Infof("Account with key %s successfully created in the system.", utils.MaskOnlyNumbers(accountKey, 6))

	c.JSON(http.StatusOK, gin.H{
		"key": accountKey,
	})
}

func GET_AccountDetails(c *gin.Context) {
	userId := c.GetInt(middlewares.ContextLoggedInUserId)

	account, err := database.SelectEntityById[database.Account](userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les détails de votre compte dans le systeme.",
		})
		return
	}

	sessions, err := database.GetAllUserSessions(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les sessions liées a votre compte.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"details":  account,
		"sessions": sessions,
	})
}
