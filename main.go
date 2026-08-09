package main

import (
	"log"
	"os"
	"post-article-be/database"
	"post-article-be/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: file .env tidak ditemukan")
	}

	// Inisialisasi: Buka atau Buat Database
	database.InitDatabase()

	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	// Inisialisasi router Gin
	router := routes.SetupRouter()

	// Menjalankan server di port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
