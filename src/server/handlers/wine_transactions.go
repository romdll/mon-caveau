package handlers

import (
	"moncaveau/database"
	"moncaveau/database/transformers"
	"moncaveau/server/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET_FetchAllTransactionsForCharts(c *gin.Context) {
	userId := c.GetInt(middlewares.ContextLoggedInUserId)

	transactions, err := database.GetAllTransactions(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation des transactions.",
		})
		return
	}

	user, err := database.SelectEntityById[database.Account](userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation de votre utilisateur.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":                transformers.PassListByTransformer(transactions, transformers.FromWineTransactionEntity),
		"accountCreationDate": user.CreatedAt,
	})
}
