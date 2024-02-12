package services

import (
	"github.com/jinzhu/gorm"
	"github.com/suleiman/Personal-Budget-Manager/models"
)

// UserService provides CRUD operations for the User model
type UserService struct {
	DB *gorm.DB
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := s.DB.First(&user, userID).Error
	return &user, err
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(userID uint, updatedUser *models.User) error {
	// Fetch the user from the database by ID
	existingUser, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	existingUser.Username = updatedUser.Username
	existingUser.Email = updatedUser.Email
	existingUser.Password = updatedUser.Password
	existingUser.Role = updatedUser.Role

	// Save the updated user to the database
	return s.DB.Save(existingUser).Error
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(userID uint) error {
	return s.DB.Delete(&models.User{}, userID).Error
}
