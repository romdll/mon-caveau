package handlers

import (
	"fmt"
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

func POST_CreateWine(c *gin.Context) {
	data := WineCreation{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Le contenu de la requete est invalide.",
		})
		return
	}

	if data.Domaine.ID == 0 && data.Domaine.Name != "" {
		res, err := database.InsertEntityById(data.Domaine)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Erreur lors de la création du nouveau domaine.",
			})
			return
		}

		newId, err := res.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la récupération d'un élement en base de données. (domain)",
			})
			return
		}

		data.Domaine.ID = int(newId)
	} else if data.Domaine.ID == 0 && data.Domaine.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni un nom de domaine invalide.",
		})
		return
	}

	if data.Region.ID == 0 && data.Region.Name != "" && data.Region.Country != "" {
		res, err := database.InsertEntityById(data.Region)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Erreur lors de la création d'une nouvelle region.",
			})
			return
		}

		newId, err := res.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la récupération d'un élement en base de données. (region)",
			})
			return
		}

		data.Region.ID = int(newId)
	} else if data.Region.ID == 0 && (data.Region.Name == "" || data.Region.Country == "") {
		if data.Region.Name == "" && data.Region.Country == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous avez fourni une region invalide.",
			})
			return
		}

		if data.Region.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous avez fourni un nom de region invalide.",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni un nom de pays pour la région invalide.",
		})
		return
	}

	if data.Type.ID == 0 && data.Type.Name != "" {
		res, err := database.InsertEntityById(data.Type)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Erreur lors de la création d'un nouveau type de vin.",
			})
			return
		}

		newId, err := res.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la récupération d'un élement en base de données. (type)",
			})
			return
		}

		data.Type.ID = int(newId)
	} else if data.Type.ID == 0 && data.Type.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni un nom type de vin invalide.",
		})
		return
	}

	if data.BottleSize.ID == 0 && data.BottleSize.Name != "" && data.BottleSize.Size != 0 {
		res, err := database.InsertEntityById(data.BottleSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Erreur lors de la création d'une nouvelle taille de bouteille.",
			})
			return
		}

		newId, err := res.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la récupération d'un élement en base de données. (size)",
			})
			return
		}

		data.BottleSize.ID = int(newId)
	} else if data.Region.ID == 0 && (data.BottleSize.Name != "" || data.BottleSize.Size != 0) {
		if data.BottleSize.Name != "" && data.BottleSize.Size != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous avez fourni une taille de bouteille invalide.",
			})
			return
		}

		if data.BottleSize.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Vous avez fourni un nom de taille de bouteille invalide.",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni une taille de bouteille invalide.",
		})
		return
	}

	fmt.Println(data)
	// userId := c.GetInt(middlewares.ContextLoggedInUserId)
}

func GET_WinesBasicStats(c *gin.Context) {

}

func GET_WinesFullStats(c *gin.Context) {

}
