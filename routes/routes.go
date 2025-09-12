package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Manual Auth
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Google Auth
	r.GET("/auth/google/login", controllers.GoogleLogin)       // redirect ke Google
	r.GET("/auth/google/callback", controllers.GoogleCallback) // login pakai Google
	r.GET("/auth/google/register", controllers.GoogleRegister) // register pakai Google

	return r
}
