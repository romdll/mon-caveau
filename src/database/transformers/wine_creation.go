package transformers

import (
	"moncaveau/database"
)

func ToEntity(data database.WineCreation) database.WineWine {
	entity := database.WineWine{}

	entity.Name = data.Name

	entity.DomainID = data.Domaine.ID
	entity.RegionID = data.Region.ID
	entity.TypeID = data.Type.ID
	entity.BottleSizeID = data.BottleSize.ID

	entity.Vintage = data.Vintage
	entity.Quantity = data.Quantity

	entity.BuyPrice = data.BuyPrice
	entity.Description = data.Description
	entity.Image = data.Image

	return entity
}
