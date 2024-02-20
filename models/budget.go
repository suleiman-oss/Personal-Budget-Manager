package models

type Budget struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	UserID   uint    `json:"user_id"`
	Month    string  `json:"month"`
}
