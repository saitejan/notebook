package models

// User struct
type User struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	Number       string        `json:"number"`
	Amount       float32       `json:"amount"`
	Transactions []interface{} `json:"transactions"`
}
