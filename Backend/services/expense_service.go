package services

import (
	"errors"
	"expense-tracker/database"
	"expense-tracker/models"

	"gorm.io/gorm"
)

type ExpenseService struct{}

func (s *ExpenseService) CreateExpense(input *models.CreateExpenseInput, userID uint) (*models.Expense, error) {
	// This is the business logic for creating an expense, called by the handler.

	// First, we check if the user provided a CategoryID. The `input.CategoryID` is a
	// pointer (*uint), so if no ID was sent, it will be `nil`.
	if input.CategoryID != nil {
		// If a CategoryID was provided, we must verify it's a valid one.
		var category models.Category
		// This database query checks for a category that meets two conditions:
		// 1. The category's `id` must match the one the user provided.
		// 2. The category must either be a default one (`user_id IS NULL`) OR
		//    it must belong to the current user (`user_id = ?`).
		// This prevents a user from assigning an expense to another user's private category.
		// We use `*input.CategoryID` to get the actual uint value from the pointer.
		if err := database.DB.Where("id = ? AND (user_id IS NULL OR user_id = ?)", *input.CategoryID, userID).First(&category).Error; err != nil {
			// If the query fails (no matching category found), we return a specific error
			// that the handler will use to send a 400 Bad Request.
			return nil, errors.New("invalid category ID provided")
		}
	}

	// Once validation is done, we create an `Expense` model struct.
	// We populate it with the validated data from the `input` and the `userID` from the token.
	expense := models.Expense{
		Title:      input.Title,
		Amount:     input.Amount,
		Date:       input.Date,
		CategoryID: input.CategoryID,
		UserID:     userID,
	}

	// We ask GORM to create a new record in the "expenses" table with our data.
	if err := database.DB.Create(&expense).Error; err != nil {
		// If the database fails to insert the record, we return a generic error.
		return nil, errors.New("failed to create expense")
	}

	// The `expense` object now has an ID, but its `Category` field is empty.
	// We use `Preload` to run a second query that fetches the associated category
	// data and populates the `expense.Category` struct field.
	database.DB.Preload("Category").First(&expense, expense.ID)

	// Finally, we return the complete expense object (with its category) and no error.
	return &expense, nil
}

func (s *ExpenseService) ListExpenses(userID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	if err := database.DB.Preload("Category").Where("user_id = ?", userID).Order("date desc").Find(&expenses).Error; err != nil {
		return nil, errors.New("failed to fetch expenses")
	}
	return expenses, nil
}

func (s *ExpenseService) UpdateExpense(expenseID string, userID uint, input *models.CreateExpenseInput) (*models.Expense, error) {
	var expense models.Expense
	if err := database.DB.Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}

	// Verify that the new Category exists
	if input.CategoryID != nil {
		var category models.Category
		if err := database.DB.Where("id = ? AND (user_id IS NULL OR user_id = ?)", *input.CategoryID, userID).First(&category).Error; err != nil {
			return nil, errors.New("invalid category ID provided")
		}
	}

	expense.Title = input.Title
	expense.Amount = input.Amount
	expense.Date = input.Date
	expense.CategoryID = input.CategoryID

	if err := database.DB.Save(&expense).Error; err != nil {
		return nil, errors.New("failed to update expense")
	}

	database.DB.Preload("Category").First(&expense, expense.ID)
	return &expense, nil
}

func (s *ExpenseService) DeleteExpense(expenseID string, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", expenseID, userID).Delete(&models.Expense{})
	if result.Error != nil {
		return errors.New("failed to delete expense")
	}
	if result.RowsAffected == 0 {
		return errors.New("expense not found")
	}
	return nil
}
