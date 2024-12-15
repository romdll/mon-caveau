package handlers

import (
	"log"
	"moncaveau/database"
	"moncaveau/utils"
)

var (
	logger *log.Logger = utils.CreateLogger("handlers")
)

type WineCreation struct {
	Name string `json:"name"`

	Domaine    database.WineDomain     `json:"domain"`
	Region     database.WineRegion     `json:"region"`
	Type       database.WineType       `json:"type"`
	BottleSize database.WineBottleSize `json:"bottle_size"`

	Vintage  int `json:"vintage"`
	Quantity int `json:"quantity"`

	BuyPrice    float64 `json:"buy_price,omitempty"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
}
