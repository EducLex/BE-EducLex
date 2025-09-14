package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		// manual
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)

		// google
		auth.GET("/google/login", controllers.GoogleLogin)
		auth.GET("/google/register", controllers.GoogleRegister)
		auth.GET("/google/callback", controllers.GoogleCallback)

		auth.GET("/user", controllers.GetUser)
	}

	return r
}
