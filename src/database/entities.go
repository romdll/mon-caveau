package database

import "time"

type Account struct {
	DB_NAME string `db:"accounts" json:"-"`

	ID         int    `db:"id" json:"-"`
	AccountKey string `db:"account_key" json:"account_key,omitempty"`
	Email      string `db:"email" json:"email,omitempty"`
	Password   string `db:"password" json:"password,omitempty"`
	Name       string `db:"name" json:"name"`
	Surname    string `db:"surname" json:"surname"`
	CreatedAt  string `db:"created_at" json:"created_at"`

	// Exceptionnal field for login ONLY
	RememberMe bool `json:"remember_me,omitempty"`
}

type Session struct {
	DB_NAME string `db:"sessions" json:"-"`

	ID           int       `db:"id" json:"id"`
	AccountID    int       `db:"account_id" json:"-"`
	SessionToken string    `db:"session_token" json:"session_token"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	LastActivity time.Time `db:"last_activity" json:"last_activity"`
}

type WineDomain struct {
	DB_NAME string `db:"wine_domains" json:"-"`

	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

type WineRegion struct {
	DB_NAME string `db:"wine_regions" json:"-"`

	ID      int    `db:"id" json:"id,omitempty"`
	Name    string `db:"name" json:"name,omitempty"`
	Country string `db:"country" json:"country,omitempty"`
}

type WineType struct {
	DB_NAME string `db:"wine_types" json:"-"`

	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

type WineBottleSize struct {
	DB_NAME string `db:"wine_bottle_sizes" json:"-"`

	ID   int     `db:"id" json:"id,omitempty"`
	Size float64 `db:"size" json:"size,omitempty"`
	Name string  `db:"name" json:"name,omitempty"`
}

type WineWine struct {
	DB_NAME string `db:"wine_wines" json:"-"`

	ID           int     `db:"id" json:"id"`
	Name         string  `db:"name" json:"name"`
	DomainID     int     `db:"domain_id" json:"domain_id"`
	RegionID     int     `db:"region_id" json:"region_id"`
	TypeID       int     `db:"type_id" json:"type_id"`
	BottleSizeID int     `db:"bottle_size_id" json:"bottle_size_id"`
	Vintage      int     `db:"vintage" json:"vintage"`
	Quantity     int     `db:"quantity" json:"quantity"`
	BuyPrice     float64 `db:"buy_price" json:"buy_price"`
	Description  string  `db:"description" json:"description"`
	Image        string  `db:"image" json:"image"`
	AccountID    int     `db:"account_id" json:"account_id"`
}

type WineTransaction struct {
	DB_NAME string `db:"wine_transactions" json:"-"`

	ID       int       `db:"id" json:"id"`
	WineID   int       `db:"wine_id" json:"wine_id"`
	Quantity int       `db:"quantity" json:"quantity"`
	Type     string    `db:"type" json:"type"`
	Date     time.Time `db:"date" json:"date"`
}
