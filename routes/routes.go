package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/google/login", controllers.GoogleLoginRedirect)
		auth.GET("/google/callback", controllers.GoogleCallback)
	}

	return r
}
