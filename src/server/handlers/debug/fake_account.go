package debug

import (
	"moncaveau/database"
	"moncaveau/database/crypt"
	"moncaveau/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var output strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := utils.SafeIntN(len(charset))
		output.WriteByte(charset[randomIndex])
	}
	return output.String()
}

func generateRandomEmail() string {
	return randomString(8) + "@moncaveau.com"
}

func generateRandomName() string {
	names := []string{"Jean", "Alice", "Robert", "Marie", "Sophie", "Daniel", "Emilie", "Christophe", "Sophie", "Valentin", "Thomas", "Guilhem", "Camille", "Lucas", "Alizé"}
	return names[utils.SafeIntN(len(names))]
}

func generateRandomWineName() string {
	wineNames := []string{
		"Chardonnay", "Merlot", "Pinot Noir", "Cabernet Sauvignon", "Syrah", "Zinfandel", "Riesling",
		"Sauvignon Blanc", "Grenache", "Malbec", "Viognier", "Chenin Blanc", "Muscadet", "Tannat",
		"Carmenère", "Carignan", "Mourvèdre",
	}
	return wineNames[utils.SafeIntN(len(wineNames))]
}

func generateRandomVintage() int {
	return utils.SafeIntN(30) + 1990
}

func generateRandomQuantity() int {
	return utils.SafeIntN(100) + 1
}

func generateRandomPrice() float64 {
	return float64(utils.SafeIntN(100) + 10)
}

func ensureWineDomainsExist() []database.WineDomain {
	domains := []string{
		"Domaine du Château", "Domaine de la Vigne", "Les Vins de Bordeaux", "Vignoble Saint-Émilion",
		"Domaine des Côtes de Provence", "Domaine du Pape", "Domaine de la Romanée-Conti",
		"Château Margaux", "Domaine des Hautes-Côtes", "Château Lafite Rothschild", "Château d'Yquem",
		"Domaine de la Côte des Blancs", "Château Pétrus", "Domaine de Montrachet",
	}

	var wineDomains []database.WineDomain
	for _, domainName := range domains {
		existingDomains, err := database.GetAllEntities[database.WineDomain]()
		if err != nil {
			logger.Errorf("Failed to get wine domains: %v", err)
			return nil
		}

		domainExists := false
		for _, existingDomain := range existingDomains {
			if existingDomain.Name == domainName {
				domainExists = true
				wineDomains = append(wineDomains, existingDomain)
				logger.Infof("Wine domain '%s' already exists.", domainName)
				break
			}
		}

		if !domainExists {
			wineDomain := &database.WineDomain{Name: domainName}
			res, err := database.InsertEntityById(wineDomain)
			if err != nil {
				logger.Errorf("Error inserting wine domain '%s': %v", domainName, err)
				return nil
			}

			id, err := res.LastInsertId()
			if err != nil {
				logger.Errorf("Error getting ID for wine domain '%s': %v", domainName, err)
				return nil
			}

			wineDomain.ID = int(id)
			wineDomains = append(wineDomains, *wineDomain)
			logger.Infof("Successfully inserted wine domain '%s'.", domainName)
		}
	}
	return wineDomains
}

func ensureWineRegionsExist() []database.WineRegion {
	regions := []struct {
		Name    string
		Country string
	}{
		{"Bordeaux", "France"}, {"Bourgogne", "France"}, {"Côtes du Rhône", "France"},
		{"Provence", "France"}, {"Languedoc", "France"}, {"Alsace", "France"}, {"Champagne", "France"},
		{"Loire Valley", "France"}, {"Beaujolais", "France"}, {"Jura", "France"},
		{"Savoie", "France"}, {"Corsica", "France"}, {"Roussillon", "France"},
	}

	var wineRegions []database.WineRegion
	for _, region := range regions {
		existingRegions, err := database.GetAllEntities[database.WineRegion]()
		if err != nil {
			logger.Errorf("Failed to get wine regions: %v", err)
			return nil
		}

		regionExists := false
		for _, existingRegion := range existingRegions {
			if existingRegion.Name == region.Name && existingRegion.Country == region.Country {
				regionExists = true
				wineRegions = append(wineRegions, existingRegion)
				logger.Infof("Wine region '%s' from '%s' already exists.", region.Name, region.Country)
				break
			}
		}

		if !regionExists {
			wineRegion := &database.WineRegion{Name: region.Name, Country: region.Country}
			res, err := database.InsertEntityById(wineRegion)
			if err != nil {
				logger.Errorf("Error inserting wine region '%s': %v", region.Name, err)
				return nil
			}

			id, err := res.LastInsertId()
			if err != nil {
				logger.Errorf("Error getting ID for wine region '%s': %v", region.Name, err)
				return nil
			}

			wineRegion.ID = int(id)
			wineRegions = append(wineRegions, *wineRegion)
			logger.Infof("Successfully inserted wine region '%s' from '%s'.", region.Name, region.Country)
		}
	}
	return wineRegions
}

func ensureWineTypesExist() []database.WineType {
	types := []string{
		"Rouge", "Blanc", "Rosé",
	}

	var wineTypes []database.WineType
	for _, wineType := range types {
		existingTypes, err := database.GetAllEntities[database.WineType]()
		if err != nil {
			logger.Errorf("Failed to get wine types: %v", err)
			return nil
		}

		typeExists := false
		for _, existingType := range existingTypes {
			if existingType.Name == wineType {
				wineTypes = append(wineTypes, existingType)
				typeExists = true
				logger.Infof("Wine type '%s' already exists.", wineType)
				break
			}
		}

		if !typeExists {
			wineType := &database.WineType{Name: wineType}
			res, err := database.InsertEntityById(wineType)
			if err != nil {
				logger.Errorf("Error inserting wine type '%s': %v", wineType, err)
				return nil
			}

			id, err := res.LastInsertId()
			if err != nil {
				logger.Errorf("Error getting ID for wine type '%s': %v", wineType, err)
				return nil
			}

			wineType.ID = int(id)
			wineTypes = append(wineTypes, *wineType)
			logger.Infof("Successfully inserted wine type '%s'.", wineType)
		}
	}
	return wineTypes
}

func ensureWineBottleSizesExist() []database.WineBottleSize {
	bottleSizes := []string{"750ml", "1L", "1.5L", "3L", "5L", "6L", "9L", "12L", "18L"}
	bottleRealSizes := []float64{750, 1000, 1500, 3000, 5000, 6000, 9000, 12000, 18000}

	var wineBottleSizes []database.WineBottleSize
	for i, bottleSizeName := range bottleSizes {
		existingBottleSizes, err := database.GetAllEntities[database.WineBottleSize]()
		if err != nil {
			logger.Errorf("Failed to get wine bottle sizes: %v", err)
			return nil
		}

		bottleSizeExists := false
		for _, existingBottleSize := range existingBottleSizes {
			if existingBottleSize.Name == bottleSizeName || existingBottleSize.Size == bottleRealSizes[i] {
				bottleSizeExists = true
				wineBottleSizes = append(wineBottleSizes, existingBottleSize)
				logger.Infof("Wine bottle size '%s' already exists.", bottleSizeName)
				break
			}
		}

		if !bottleSizeExists {
			wineBottleSize := &database.WineBottleSize{Name: bottleSizeName, Size: bottleRealSizes[i]}
			res, err := database.InsertEntityById(wineBottleSize)
			if err != nil {
				logger.Errorf("Error inserting wine bottle size '%s': %v", bottleSizeName, err)
				return nil
			}

			id, err := res.LastInsertId()
			if err != nil {
				logger.Errorf("Error getting ID for wine bottle size '%s': %v", bottleSizeName, err)
				return nil
			}

			wineBottleSize.ID = int(id)
			wineBottleSizes = append(wineBottleSizes, *wineBottleSize)
			logger.Infof("Successfully inserted wine bottle size '%s'.", bottleSizeName)
		}
	}
	return wineBottleSizes
}

func GET_CreateFakeAccount(c *gin.Context) {
	logger.Info("Starting the creation of a fake account and wines.")

	accountKey, err := crypt.GenerateSecureAccountKey()
	if err != nil {
		logger.Errorf("Error generating account key: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate account key"})
		return
	}

	logger.Infof("Generated account key for new fake account: %s", accountKey)

	wineDomains := ensureWineDomainsExist()
	wineRegions := ensureWineRegionsExist()
	wineTypes := ensureWineTypesExist()
	wineBottleSizes := ensureWineBottleSizesExist()

	fakeAccount := &database.Account{
		AccountKey: accountKey,
		Email:      generateRandomEmail(),
		Name:       generateRandomName(),
		Surname:    generateRandomName(),
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	result, err := database.InsertEntityById(fakeAccount)
	if err != nil {
		logger.Errorf("Error inserting account into the database: %v", err)
		c.JSON(500, gin.H{"error": "Failed to insert account"})
		return
	}

	accountId, err := result.LastInsertId()
	if err != nil {
		logger.Errorf("Can't get the ID of the fake account: %v", err)
		c.JSON(500, gin.H{"error": "Failed to insert account"})
		return
	}

	fakeAccount.ID = int(accountId)
	logger.Infof("Successfully created fake account with ID: %d", fakeAccount.ID)

	for i := 0; i < 20+utils.SafeIntN(10); i++ {
		fakeWine := &database.WineWine{
			Name:         generateRandomWineName(),
			DomainID:     wineDomains[utils.SafeIntN(len(wineDomains))].ID,
			RegionID:     wineRegions[utils.SafeIntN(len(wineRegions))].ID,
			TypeID:       wineTypes[utils.SafeIntN(len(wineTypes))].ID,
			BottleSizeID: wineBottleSizes[utils.SafeIntN(len(wineBottleSizes))].ID,
			Vintage:      generateRandomVintage(),
			Quantity:     generateRandomQuantity(),
			BuyPrice:     generateRandomPrice(),
			Description:  "Un excellent vin à déguster et apprécier !",
			AccountID:    fakeAccount.ID,
		}

		_, err := database.InsertEntityById(fakeWine)
		if err != nil {
			logger.Errorf("Error inserting wine into the database: %v", err)
			c.JSON(500, gin.H{"error": "Failed to insert wine"})
			return
		}

		logger.Infof("Inserted wine '%s' for fake account ID %d", fakeWine.Name, fakeAccount.ID)
	}

	logger.Info("Successfully created fake account and 20 wines.")
	c.JSON(200, gin.H{
		"accountKey": accountKey,
	})
}
