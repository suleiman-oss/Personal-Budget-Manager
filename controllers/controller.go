package controllers

import "github.com/jinzhu/gorm"

type Controller struct {
	DB *gorm.DB
}
