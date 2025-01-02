package database

import "time"

type WineCreation struct {
	Name string `json:"name"`

	Domaine    WineDomain     `json:"domain"`
	Region     WineRegion     `json:"region"`
	Type       WineType       `json:"type"`
	BottleSize WineBottleSize `json:"bottle_size"`

	Vintage  int `json:"vintage"`
	Quantity int `json:"quantity"`

	BuyPrice    float64 `json:"buy_price,omitempty"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
}

type WineTransactionForChart struct {
	Quantity int       `json:"quantity"`
	Type     string    `json:"type"`
	Date     time.Time `json:"date"`
}
