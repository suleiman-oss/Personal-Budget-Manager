package services

import (
	"github.com/jinzhu/gorm"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
	"golang.org/x/crypto/bcrypt"
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
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil // Return nil if the user is not found
		}
		return nil, err
	}
	return &user, nil
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

	// Save the updated user to the database
	return s.DB.Save(existingUser).Error
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(userID uint) error {
	return s.DB.Delete(&models.User{}, userID).Error
}

func (s *UserService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
