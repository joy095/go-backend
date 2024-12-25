package repositories

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joy095/backend/models"
)

func CreateUser(db *pgxpool.Pool, user models.User) error {
	// Explicitly set is_verified to false
	isVerified := false

	log.Printf("Inserting user: email=%s, password=%s, is_verified=%t, username=%s",
		user.Email, user.Password, isVerified, user.Username)

	_, err := db.Exec(context.Background(),
		"INSERT INTO users (email, password, is_verified, user_name) VALUES ($1, $2, $3, $4)",
		user.Email, user.Password, isVerified, user.Username)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		os.Exit(1)
	}

	return nil
}

func GetUserByEmail(db *pgxpool.Pool, email string) (*models.User, error) {
	var user models.User
	err := db.QueryRow(context.Background(),
		"SELECT id, email, is_verified, user_name, created_at FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.IsVerified, &user.Username, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(db *pgxpool.Pool, username string) (*models.PublicUser, error) {
	query := `
		SELECT id, email, is_verified, user_name, created_at
		FROM users
		WHERE user_name = $1
	`

	var user models.PublicUser

	err := db.QueryRow(context.Background(), query, username).Scan(
		&user.ID,
		&user.Email,
		&user.IsVerified,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserVerification(db *pgxpool.Pool, email string, isVerified bool) error {
	_, err := db.Exec(context.Background(),
		"UPDATE users SET is_verified=$1 WHERE email=$2", isVerified, email)
	return err
}

func DeleteUser(db *pgxpool.Pool, email string) error {
	_, err := db.Exec(context.Background(), "DELETE FROM users WHERE email=$1", email)
	return err
}

type UserRepository interface {
	SaveOTP(username, otp string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveOTP(username, otp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO email_verifications (username, otp, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (username) 
		DO UPDATE SET otp = EXCLUDED.otp, created_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, username, otp)
	return err
}
