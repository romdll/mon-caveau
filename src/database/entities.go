package database

type Account struct {
	DB_NAME string `db:"accounts"`

	ID         int    `db:"id"`
	AccountKey string `db:"account_key" json:"account_key,omitempty"`
	Email      string `db:"email" json:"email,omitempty"`
	Password   string `db:"password" json:"password,omitempty"`
	Name       string `db:"name"`
	Surname    string `db:"surname"`
	CreatedAt  string `db:"created_at"`
}
