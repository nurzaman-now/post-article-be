package routes

import (
	"post-article-be/controllers"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// SetupRouter - Setup semua routes untuk aplikasi
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Mengaktifkan CORS middleware
	router.Use(CORSMiddleware())

	// Inisialisasi controller
	postController := controllers.InitPostController()

	// Setup routes untuk post
	postRoutes := router.Group("/article")
	{
		postRoutes.POST("", postController.Create)
		// Menggunakan RouteGet dengan catch-all wildcard "/*action" untuk menghindari panic wildcard conflict di Gin.
		// Gin tidak memperbolehkan dua endpoint dengan format parameter berbeda pada segmen yang sama (misal: /:limit/:offset vs /:id).
		// RouteGet akan mem-parsing secara dinamis URL path menjadi parameter limit/offset (untuk paging) atau id (untuk detail).
		postRoutes.GET("/*action", postController.RouteGet)
		postRoutes.PUT("/:id", postController.Update)
		postRoutes.DELETE("/:id", postController.Delete)
	}

	return router
}
