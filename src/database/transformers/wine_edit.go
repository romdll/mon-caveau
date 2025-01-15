package transformers

import "moncaveau/database"

func EditToWineWineEntity(entry database.WineEdit) database.WineWine {
	entity := database.WineWine{}

	entity.ID = entry.Id
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

	if entry.PreferredStartDate != "" {
		entity.PreferredStartDate = &entry.PreferredStartDate
	}

	if entry.PreferredEndDate != "" {
		entity.PreferredEndDate = &entry.PreferredEndDate
	}

	return entity
}
