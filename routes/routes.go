package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Manual auth
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Google auth
	r.GET("/auth/google/login", controllers.GoogleLogin)
	r.GET("/auth/google/callback", controllers.GoogleCallback)

	// Protected route
	r.GET("/profile", middleware.AuthMiddleware(), controllers.ProfileHandler)
}
