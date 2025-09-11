package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/EducLex/BE-EducLex/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// --- Auth Manual ---
	r.POST("/auth/register", controllers.Register)
	r.POST("/auth/login", controllers.Login)

	// --- Google OAuth ---
	r.GET("/auth/google/login", controllers.GoogleLogin)
	r.GET("/auth/google/callback", controllers.GoogleCallback)

	// --- Protected Routes ---
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controllers.ProfileHandler)
	}

	return r
}
