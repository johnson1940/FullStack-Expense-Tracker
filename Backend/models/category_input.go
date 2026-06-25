package models

// CreateCategoryInput is used for validating create category requests.
type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}
