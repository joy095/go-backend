package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	Username   string    `json:"username" binding:"required"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type PublicUser struct {
	ID         int       `json:"id"`
	Email      string    `json:"email" binding:"required"`
	Username   string    `json:"username" binding:"required"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}
