package models

type Notification struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Message string `json:"message"`
	Type    string `json:"type"`
	UserID  uint   `json:"user_id"`
}
