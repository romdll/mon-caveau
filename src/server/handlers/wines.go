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

	c.JSON(http.StatusOK, gin.H{
		"totalWines":                   totalWines,
		"totalWinesDrankSold":          totalWinesDrankSold,
		"totalWinesDrankSoldThisMonth": totalWinesDrankSoldThisMonth,
	})
}

func GET_WinesBasicStats(c *gin.Context) {

}

func GET_WinesFullStats(c *gin.Context) {

}
