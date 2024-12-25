package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/joy095/backend/api/routes"
	database "github.com/joy095/backend/config"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	db := database.InitDB()
	defer db.Close()

	r := gin.Default()

	routes.UserRoutes(r, db)

	port := os.Getenv("PORT")

	log.Println("Server running on http://localhost:" + port)
	r.Run(":" + port)
}
