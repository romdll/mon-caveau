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

	PreferredStartDate string `json:"preferred_start_date,omitempty"`
	PreferredEndDate   string `json:"preferred_end_date,omitempty"`
}

type WineTransactionForChart struct {
	Quantity int       `json:"quantity"`
	Type     string    `json:"type"`
	Date     time.Time `json:"date"`
}

type SessionActivity struct {
	SessionToken string
	LastActivity time.Time
}
