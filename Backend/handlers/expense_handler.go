package handlers

import (
	"expense-tracker/models"
	"expense-tracker/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var expenseService = services.ExpenseService{}

// CreateExpense handles POST /expenses
func CreateExpense(c *gin.Context) {
	// Get the user's ID from the Gin context. The `AuthRequired` middleware
	// verified the JWT and placed the user's ID here for us to use.
	userID, exists := c.Get("userID")
	// If, for some reason, the userID isn't here, it means the user is not
	// authenticated. We should stop immediately.
	if !exists {
		// Return a 401 Unauthorized error to the client.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// Stop further execution of this function.
		return
	}

	// Create an empty variable that will hold the JSON data from the request.
	// We use the `CreateExpenseInput` struct which has validation tags.
	var input models.CreateExpenseInput
	// `ShouldBindJSON` attempts to parse the request's JSON body into our `input`
	// struct. It also automatically validates the fields based on the `binding` tags.
	if err := c.ShouldBindJSON(&input); err != nil {
		// If parsing or validation fails, return a 400 Bad Request error
		// with the specific validation error message.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the `CreateExpense` method from our `expenseService`. This is where the
	// actual business logic lives. We pass the validated input and the user's ID.
	// We use `userID.(uint)` to convert the `userID` from an `interface{}` to a `uint`.
	expense, err := expenseService.CreateExpense(&input, userID.(uint))
	// Check if the service layer returned an error.
	if err != nil {
		// If the specific error was an invalid category ID, it's the client's fault.
		if err.Error() == "invalid category ID provided" {
			// So, we return a 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			// For any other kind of error (like a database failure), it's a server-side
			// problem, so we return a 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		// Stop execution.
		return
	}

	// If everything was successful, return a 201 Created status code.
	// We also return a success message and the newly created expense object.
	c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully", "data": expense})
}

// ListExpenses handles GET /expenses
func ListExpenses(c *gin.Context) {
	userID, _ := c.Get("userID")

	expenses, err := expenseService.ListExpenses(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expenses retrieved successfully", "data": expenses})
}

// UpdateExpense handles PUT /expenses/:id
func UpdateExpense(c *gin.Context) {
	expenseID := c.Param("id")
	userID, _ := c.Get("userID")

	var input models.CreateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense, err := expenseService.UpdateExpense(expenseID, userID.(uint), &input)
	if err != nil {
		switch err.Error() {
		case "expense not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "invalid category ID provided":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully", "data": expense})
}

// DeleteExpense handles DELETE /expenses/:id
func DeleteExpense(c *gin.Context) {
	expenseID := c.Param("id")
	userID, _ := c.Get("userID")

	err := expenseService.DeleteExpense(expenseID, userID.(uint))
	if err != nil {
		if err.Error() == "expense not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
