package services

import (
	"errors"
	"expense-tracker/database"
	"expense-tracker/models"
)

type CategoryService struct{}

func (s *CategoryService) CreateCategory(input *models.CreateCategoryInput, userID uint) (*models.Category, error) {
	category := models.Category{Name: input.Name, UserID: &userID}
	if err := database.DB.Create(&category).Error; err != nil {
		return nil, errors.New("failed to create category")
	}
	return &category, nil
}

func (s *CategoryService) ListCategories(userID uint) ([]models.Category, error) {
	var categories []models.Category

	// Get categories where user_id IS NULL (defaults) OR user_id matches the logged-in user
	if err := database.DB.Where("user_id IS NULL OR user_id = ?", userID).Order("id asc").Find(&categories).Error; err != nil {
		return nil, errors.New("failed to fetch categories")
	}

	if categories == nil {
		categories = []models.Category{}
	}

	return categories, nil
}
