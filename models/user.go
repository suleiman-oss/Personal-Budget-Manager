package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (user *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
