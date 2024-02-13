package models

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates tables based on the provided GORM DB instance
func Migrate(db *gorm.DB) {
	// Auto-migrate creates tables based on the struct definitions
	db.AutoMigrate(&User{}, &Expense{}, &Budget{}, &Notification{})
}
