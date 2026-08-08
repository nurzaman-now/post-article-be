package main

import (
	"backend/database"
	"backend/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi: Buka atau Buat Database
	database.InitDatabase()

	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	// Inisialisasi router Gin
	router := routes.SetupRouter()

	// Menjalankan server di port 8080
	router.Run(":8080")
}
