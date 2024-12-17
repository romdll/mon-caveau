package handlers

import (
	"moncaveau/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET_GetCountriesAndRegions(c *gin.Context) {
	logger.Println("Starting to fetch countries and regions...")

	countriesAndRegions, err := database.GetAllEntities[database.WineRegion]()
	if err != nil {
		logger.Printf("Error retrieving countries and regions: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les pays et les régions.",
		})
		return
	}

	logger.Printf("Successfully retrieved %d countries and regions.\n", len(countriesAndRegions))
	c.JSON(http.StatusOK, countriesAndRegions)
}

func GET_GetWineTypes(c *gin.Context) {
	logger.Println("Starting to fetch wine types...")

	wineTypes, err := database.GetAllEntities[database.WineType]()
	if err != nil {
		logger.Printf("Error retrieving wine types: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les type de vin.",
		})
		return
	}

	logger.Printf("Successfully retrieved %d wine types.\n", len(wineTypes))
	c.JSON(http.StatusOK, wineTypes)
}

func GET_GetBottleSizes(c *gin.Context) {
	logger.Println("Starting to fetch wine bottle sizes...")

	wineBottleSizes, err := database.GetAllEntities[database.WineBottleSize]()
	if err != nil {
		logger.Printf("Error retrieving wine bottle sizes: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les tailles de bouteilles de vin.",
		})
		return
	}

	logger.Printf("Successfully retrieved %d wine bottle sizes.\n", len(wineBottleSizes))
	c.JSON(http.StatusOK, wineBottleSizes)
}

func GET_GetDomains(c *gin.Context) {
	logger.Println("Starting to fetch wine domains...")

	wineDomains, err := database.GetAllEntities[database.WineDomain]()
	if err != nil {
		logger.Printf("Error retrieving wine domains: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les domaines.",
		})
		return
	}

	logger.Printf("Successfully retrieved %d wine domains.\n", len(wineDomains))
	c.JSON(http.StatusOK, wineDomains)
}
