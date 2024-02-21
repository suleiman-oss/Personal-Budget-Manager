package models

import (
	"time"
)

// Goal model
type Goal struct {
	Name         string    `json:"name" gorm:"primaryKey"`
	TargetAmount float64   `json:"target_amount"`
	UserID       uint      `json:"user_id"`
	Month        time.Time `json:"month"`
}
