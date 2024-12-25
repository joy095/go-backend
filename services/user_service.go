package services

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joy095/backend/models"
	"github.com/joy095/backend/repositories"
	"golang.org/x/crypto/argon2"
)

// HashPassword hashes a password using Argon2.
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Hash the password with Argon2
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Return the hashed password as a base64-encoded string
	encoded := base64.StdEncoding.EncodeToString(append(salt, hash...))
	return encoded, nil
}

// CreateUser hashes the password and saves the user.
func CreateUser(db *pgxpool.Pool, user models.User) error {
	// Hash the user's password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return err
	}

	// Set the hashed password and pass the user to the repository
	user.Password = hashedPassword
	return repositories.CreateUser(db, user)
}
