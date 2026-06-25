package handlers

import (
	"expense-tracker/models"
	"expense-tracker/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var categoryService = services.CategoryService{}

// CreateCategory allows a user to create a custom category
func CreateCategory(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input models.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := categoryService.CreateCategory(&input, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Custom category created", "data": category})
}

// ListCategories returns the default categories + the user's custom categories
func ListCategories(c *gin.Context) {
	userID, _ := c.Get("userID")

	categories, err := categoryService.ListCategories(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categories retrieved successfully", "data": categories})
}
