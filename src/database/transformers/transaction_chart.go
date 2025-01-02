package transformers

import "moncaveau/database"

func FromWineTransactionEntity(entry database.WineTransaction) database.WineTransactionForChart {
	res := database.WineTransactionForChart{}

	res.Quantity = entry.Quantity
	res.Type = entry.Type
	res.Date = entry.Date

	return res
}
