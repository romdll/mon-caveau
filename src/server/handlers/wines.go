package handlers

import (
	"errors"
	"fmt"
	"moncaveau/database"
	"moncaveau/database/transformers"
	"moncaveau/server/middlewares"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GET_WinesDashboard(c *gin.Context) {
	userId := c.GetInt(middlewares.ContextLoggedInUserId)

	totalWines, totalWinesDrankSold, realTotalBottlesAdded, totalCurrentBottles, err := database.GetWinesForDashboard(userId)
	if err != nil {
		logger.Errorw("Error when getting basic dashboard data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation des données basiques.",
		})
		return
	}

	winesCountPerRegions, err := database.GetWinesCountPerRegion(userId)
	if err != nil {
		logger.Errorw("Error when getting wines count per region", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation du nombre de bouteilles par regions.",
		})
		return
	}

	winesCountPerTypes, err := database.GetWinesCountPerTypes(userId)
	if err != nil {
		logger.Errorw("Error when getting wines count per types", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation du nombre de bouteilles par type.",
		})
		return
	}

	last4Transactions, winesIdToName, err := database.Get4LatestsTransactions(userId)
	if err != nil {
		logger.Errorw("Error when getting the last 4 transactions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récuperation des 4 dernieres transactions.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalWines":            totalWines,
		"totalCurrentBottles":   totalCurrentBottles,
		"totalWinesDrankSold":   totalWinesDrankSold,
		"realTotalBottlesAdded": realTotalBottlesAdded,
		"winesCountPerRegions":  winesCountPerRegions,
		"winesCountPerTypes":    winesCountPerTypes,
		"last4Transactions": gin.H{
			"data":          last4Transactions,
			"winesIdToName": winesIdToName,
		},
	})
}

func validateDrinkingDates(startDate, endDate string) error {
	const layout = "2006-01-02"

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %v", err)
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %v", err)
	}

	if start.Equal(end) {
		return errors.New("start date and end date cannot be the same")
	}

	if start.After(end) {
		return errors.New("start date cannot be after end date")
	}

	return nil
}

func POST_CreateWine(c *gin.Context) {
	data := database.WineCreation{}
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
	} else if data.Domaine.ID < 0 {
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
	} else if data.Region.ID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni un nom de region invalide.",
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
	} else if data.Type.ID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni un type de vin invalide.",
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
	} else if data.BottleSize.ID == 0 && (data.BottleSize.Name != "" || data.BottleSize.Size != 0) {
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
	} else if data.BottleSize.ID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni une taille de bouteille de vin invalide.",
		})
		return
	}

	if data.PreferredStartDate != "" && data.PreferredEndDate != "" && validateDrinkingDates(data.PreferredStartDate, data.PreferredEndDate) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni des dates de consommation invalide.",
		})
		return
	}

	// TODO verify all the sub ids to make sure they exists

	userId := c.GetInt(middlewares.ContextLoggedInUserId)
	realWine := transformers.ToWineWineEntity(data)
	realWine.AccountID = userId

	_, err := database.InsertEntityById(realWine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de créer votre vin.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func POST_editWine(c *gin.Context) {
	data := database.WineEdit{}
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
	} else if data.Domaine.ID < 0 {
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
	} else if data.Region.ID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni un nom de region invalide.",
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
	} else if data.Type.ID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni un type de vin invalide.",
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
	} else if data.BottleSize.ID == 0 && (data.BottleSize.Name != "" || data.BottleSize.Size != 0) {
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
	} else if data.BottleSize.ID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni une taille de bouteille de vin invalide.",
		})
		return
	}

	if data.PreferredStartDate != "" && data.PreferredEndDate != "" && validateDrinkingDates(data.PreferredStartDate, data.PreferredEndDate) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous avez fourni des dates de consommation invalide.",
		})
		return
	}

	// TODO verify all the sub ids to make sure they exists

	userId := c.GetInt(middlewares.ContextLoggedInUserId)
	realWine := transformers.EditToWineWineEntity(data)
	realWine.AccountID = userId

	_, err := database.UpdateEntityById(realWine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de mettre à jour votre vin.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func GET_allWines(c *gin.Context) {
	userId := c.GetInt(middlewares.ContextLoggedInUserId)

	pageParam := c.DefaultQuery("page", "1")
	limitParam := c.DefaultQuery("limit", "3")
	searchQuery := c.DefaultQuery("search", "")
	filterPreferredDatesParam := c.DefaultQuery("filterPreferredDates", "false")

	logger.Infof("Received parameters: page=%s, limit=%s, search='%s', filterPreferredDates=%s", pageParam, limitParam, searchQuery, filterPreferredDatesParam)

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		logger.Warnf("Invalid 'page' parameter: %s", pageParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Paramètre 'page' invalide.",
		})
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		logger.Warnf("Invalid 'limit' parameter: %s", limitParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Paramètre 'limit' invalide.",
		})
		return
	}

	filterPreferredDates := filterPreferredDatesParam == "true"

	logger.Infof("Validated parameters: userId=%d, page=%d, limit=%d, search='%s', filterPreferredDates=%t", userId, page, limit, searchQuery, filterPreferredDates)

	wines, winesCount, err := database.GetWinesWithPaginationAndSearch(limit, page, userId, searchQuery, filterPreferredDates)
	if err != nil {
		logger.Errorf("Error fetching wines for userId=%d with search='%s' and filterPreferredDates=%t: %v", userId, searchQuery, filterPreferredDates, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur interne lors de la recupération de vos vins.",
		})
		return
	}

	logger.Infof("Successfully fetched %d wines for userId=%d with search='%s' and filterPreferredDates=%t", len(wines), userId, searchQuery, filterPreferredDates)
	logger.Infof("Total wine count for userId=%d is %d", userId, winesCount)

	c.JSON(http.StatusOK, gin.H{
		"wines": wines,
		"total": winesCount,
	})
}

func GET_wineStatisticsData(c *gin.Context) {
	userId := c.GetInt(middlewares.ContextLoggedInUserId)
	logger.Debugf("Fetching wine statistics data for userId: %d", userId)

	top5Domains, err := database.GetTop5DomainsPerNumberOfBottles(userId)
	if err != nil {
		logger.Errorf("Error fetching top 5 domains for userId %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les 5 domaines avec le plus de vins dans votre collection.",
		})
		return
	}
	logger.Debugf("Successfully fetched top 5 domains for userId %d: %v", userId, top5Domains)

	wineDistributionPerVintage, err := database.GetWineDistributionPerVintage(userId)
	if err != nil {
		logger.Errorf("Error fetching wine distribution per vintage for userId %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer la distributions par année de votre collection.",
		})
		return
	}
	logger.Debugf("Successfully fetched wine distribution per vintage for userId %d: %v", userId, wineDistributionPerVintage)

	wineTypesDistributionPerRegion, err := database.GetWineTypesDistributionPerRegions(userId)
	if err != nil {
		logger.Errorf("Error fetching wine types distribution per region for userId %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer la distribution des types de vins par régions.",
		})
		return
	}
	logger.Debugf("Successfully fetched wine types distribution per region for userId %d: %v", userId, wineTypesDistributionPerRegion)

	userUsedRegionsWithBottlecount, err := database.GetUserUsedRegionsWithBottleCount(userId)
	if err != nil {
		logger.Errorf("Error fetching user used regions with bottle count for userId %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les régions que vous utilisez dans votre collection.",
		})
		return
	}
	logger.Debugf("Successfully fetched user used regions with bottle count for userId %d: %v", userId, userUsedRegionsWithBottlecount)

	logger.Infof("Successfully fetched all wine statistics data for userId %d", userId)
	c.JSON(http.StatusOK, gin.H{
		"top5Domains":                    top5Domains,
		"wineDistributionPerVintage":     wineDistributionPerVintage,
		"wineTypesDistributionPerRegion": wineTypesDistributionPerRegion,
		"userUsedRegionsWithBottlecount": userUsedRegionsWithBottlecount,
	})
}

func POST_AdjustQuantity(c *gin.Context) {
	wineID, err := strconv.Atoi(c.Param("id"))
	if err != nil || wineID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant de vin invalide.",
		})
		return
	}

	var req AdjustQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "L'ajustement de quantité est impossible avec les données fournies.",
		})
		return
	}

	wine, err := database.SelectEntityById[database.WineWine](wineID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant de vin invalide.",
		})
		return
	}

	userId := c.GetInt(middlewares.ContextLoggedInUserId)
	if wine.AccountID != userId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous essayez de modifié un vin qui n'est pas a vous.",
		})
		return
	}

	wine.Quantity += req.Change

	if wine.Quantity < 0 {
		wine.Quantity = 0
	}

	_, err = database.UpdateEntityById(wine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de mettre a jour la quantité de bouteilles.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quantity": wine.Quantity})
}

func DELETE_deleteWine(c *gin.Context) {
	wineID, err := strconv.Atoi(c.Param("id"))
	if err != nil || wineID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant de vin invalide.",
		})
		return
	}

	wine, err := database.SelectEntityById[database.WineWine](wineID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant de vin invalide.",
		})
		return
	}

	userId := c.GetInt(middlewares.ContextLoggedInUserId)
	if wine.AccountID != userId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous essayez de suprimé un vin qui n'est pas a vous.",
		})
		return
	}

	wine.Quantity = 0

	_, err = database.UpdateEntityById(wine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de mettre a jour la quantité de bouteilles.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func GET_wineById(c *gin.Context) {
	wineID, err := strconv.Atoi(c.Param("id"))
	if err != nil || wineID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant de vin invalide.",
		})
		return
	}

	wine, err := database.SelectEntityById[database.WineWine](wineID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant de vin invalide.",
		})
		return
	}

	userId := c.GetInt(middlewares.ContextLoggedInUserId)
	if wine.AccountID != userId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vous essayez d'acceder à un vin qui n'est pas a vous.",
		})
		return
	}

	c.JSON(http.StatusOK, wine)
}
