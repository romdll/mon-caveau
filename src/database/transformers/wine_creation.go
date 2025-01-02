package transformers

import (
	"moncaveau/database"
)

func ToWineWineEntity(entry database.WineCreation) database.WineWine {
	entity := database.WineWine{}

	entity.Name = entry.Name

	entity.DomainID = entry.Domaine.ID
	entity.RegionID = entry.Region.ID
	entity.TypeID = entry.Type.ID
	entity.BottleSizeID = entry.BottleSize.ID

	entity.Vintage = entry.Vintage
	entity.Quantity = entry.Quantity

	entity.BuyPrice = entry.BuyPrice
	entity.Description = entry.Description
	entity.Image = entry.Image

	return entity
}
