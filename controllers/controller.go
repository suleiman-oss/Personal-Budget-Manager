package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/suleiman-oss/Personal-Budget-Manager/services"
)

type Controller struct {
	UserController    *UserController
	BudgetController  *BudgetController
	ExpenseController *ExpenseController
}

func NewController(db *gorm.DB) *Controller {
	cr := Controller{}

	ucr := UserController{}
	ucr.UserService = &services.UserService{}
	ucr.UserService.DB = db
	cr.UserController = &ucr

	bcr := BudgetController{}
	bcr.BudgetService = &services.BudgetService{}
	bcr.BudgetService.DB = db
	cr.BudgetController = &bcr

	ecr := ExpenseController{}
	ecr.ExpenseService = &services.ExpenseService{}
	ecr.ExpenseService.DB = db
	cr.ExpenseController = &ecr

	return &cr
}
