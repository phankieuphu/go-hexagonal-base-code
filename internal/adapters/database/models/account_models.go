package models

type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (a Account) TableName() string {
	return "account"
}
