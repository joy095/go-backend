package routes

import (
	"github.com/joy095/backend/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(r *gin.Engine, db *pgxpool.Pool) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", handlers.CreateUserHandler(db))
		userGroup.GET("/:username", handlers.GetUserHandler(db))
		userGroup.PUT("/:email/verify", handlers.VerifyUserHandler(db))
		userGroup.DELETE("/:username", handlers.DeleteUserHandler(db))
	}
}
