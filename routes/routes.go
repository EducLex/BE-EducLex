package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/EducLex/BE-EducLex/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// -------- AUTH MANUAL --------
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// -------- AUTH GOOGLE --------
	r.GET("/auth/google/login", controllers.GoogleLogin)
	r.GET("/auth/google/callback", controllers.GoogleCallback)

	// -------- PROTECTED --------
	r.GET("/profile", middleware.AuthMiddleware(), controllers.ProfileHandler)
}
