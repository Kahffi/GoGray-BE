package models

type Image struct {
	ImageBase
	ImgurLink string `json:"imgurLink"`
}

type ImageBase struct {
	Id            string `gorm:"primaryKey" json:"id"`
	TransactionID string `gorm:"index" json:"transactionID"`
	Status        string `json:"status"`
}
