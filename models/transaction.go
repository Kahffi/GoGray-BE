package models

import "time"

type Transaction struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `json:"userID"`
	Images    []Image   `gorm:"foreignKey:TransactionID" json:"images"`
	CreatedAt time.Time `json:"createdAt"`
}

type TransactionCreateRequest struct {
	Id     string   `json:"id"`
	UserID string   `json:"userID"`
	Images []string `json:"images"`
	// CreatedAt time.Time `json:"createdAt"`
}
