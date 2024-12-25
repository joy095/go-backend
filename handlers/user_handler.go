package handlers

import (
	"net/http"

	"github.com/joy095/backend/models"
	"github.com/joy095/backend/repositories"
	"github.com/joy095/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Validate required fields
		if user.Email == "" || user.Password == "" || user.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email, password, and username are required"})
			return
		}

		// Call the service
		err := services.CreateUser(db, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

// GetUserHandler fetches a user by email.
func GetUserHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// email := c.Param("email")
		username := c.Param("username")

		// user, err := repositories.GetUserByEmail(db, email)
		user, err := repositories.GetUserByUsername(db, username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// VerifyUserHandler updates the verification status of a user.
func VerifyUserHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")

		err := repositories.UpdateUserVerification(db, email, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
	}
}

// DeleteUserHandler deletes a user by username.
func DeleteUserHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		err := repositories.DeleteUser(db, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
