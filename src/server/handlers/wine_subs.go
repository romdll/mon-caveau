package handlers

import (
	"moncaveau/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET_GetCountriesAndRegions(c *gin.Context) {
	logger.Infow("Starting to fetch countries and regions...")

	countriesAndRegions, err := database.GetAllEntities[database.WineRegion]()
	if err != nil {
		logger.Errorw("Error retrieving countries and regions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les pays et les régions.",
		})
		return
	}

	logger.Infow("Successfully retrieved countries and regions", "count", len(countriesAndRegions))
	c.JSON(http.StatusOK, countriesAndRegions)
}

func GET_GetWineTypes(c *gin.Context) {
	logger.Infow("Starting to fetch wine types...")

	wineTypes, err := database.GetAllEntities[database.WineType]()
	if err != nil {
		logger.Errorw("Error retrieving wine types", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les type de vin.",
		})
		return
	}

	logger.Infow("Successfully retrieved wine types", "count", len(wineTypes))
	c.JSON(http.StatusOK, wineTypes)
}

func GET_GetBottleSizes(c *gin.Context) {
	logger.Infow("Starting to fetch wine bottle sizes...")

	wineBottleSizes, err := database.GetAllEntities[database.WineBottleSize]()
	if err != nil {
		logger.Errorw("Error retrieving wine bottle sizes", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les tailles de bouteilles de vin.",
		})
		return
	}

	logger.Infow("Successfully retrieved wine bottle sizes", "count", len(wineBottleSizes))
	c.JSON(http.StatusOK, wineBottleSizes)
}

func GET_GetDomains(c *gin.Context) {
	logger.Infow("Starting to fetch wine domains...")

	wineDomains, err := database.GetAllEntities[database.WineDomain]()
	if err != nil {
		logger.Errorw("Error retrieving wine domains", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les domaines.",
		})
		return
	}

	logger.Infow("Successfully retrieved wine domains", "count", len(wineDomains))
	c.JSON(http.StatusOK, wineDomains)
}
