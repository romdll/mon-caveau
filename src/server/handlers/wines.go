package handlers

import (
	"moncaveau/database"
	"moncaveau/server/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET_WinesDashboard(c *gin.Context) {
	userId := c.GetInt(middlewares.ContextLoggedInUserId)

	totalWines, totalWinesDrankSold, totalWinesDrankSoldThisMonth, err := database.GetWinesForDashboard(userId)
	if err != nil {
		logger.Printf("Error when getting basic dashboard data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation des données basiques.",
		})
		return
	}

	winesCountPerRegions, err := database.GetWinesCountPerRegion(userId)
	if err != nil {
		logger.Printf("Error when getting wines count per region: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation du nombre de bouteilles par regions.",
		})
		return
	}

	winesCountPerTypes, err := database.GetWinesCountPerTypes(userId)
	if err != nil {
		logger.Printf("Error when getting wines count per types: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation du nombre de bouteilles par type.",
		})
		return
	}

	last4Transactions, winesIdToName, err := database.Get4LatestsTransactions(userId)
	if err != nil {
		logger.Printf("Error when getting the last 4 transactions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation des 4 dernieres transactions.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalWines":                   totalWines,
		"totalWinesDrankSold":          totalWinesDrankSold,
		"totalWinesDrankSoldThisMonth": totalWinesDrankSoldThisMonth,
		"winesCountPerRegions":         winesCountPerRegions,
		"winesCountPerTypes":           winesCountPerTypes,
		"last4Transactions": gin.H{
			"data":          last4Transactions,
			"winesIdToName": winesIdToName,
		},
	})
}

func GET_WinesBasicStats(c *gin.Context) {

}

func GET_WinesFullStats(c *gin.Context) {

}
