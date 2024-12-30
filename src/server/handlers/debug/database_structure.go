package debug

import (
	"moncaveau/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET_DatabaseStructure(c *gin.Context) {
	tables, err := database.GetAllTablesAndStructures()
	if err != nil {
		logger.Errorf("error fetching tables: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tables)
}
