package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// auth group
	auth := r.Group("/auth")
	{
		// manual login/register
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)

		// google login/register
		auth.GET("/google/login", controllers.GoogleLogin)
		auth.GET("/google/register", controllers.GoogleRegister)
		auth.GET("/google/callback", controllers.GoogleCallback)
	}

	// protected routes (butuh token JWT)
	protected := r.Group("/auth")
	protected.Use(middleware.AuthMiddleware()) // middleware dipasang disini
	{
		protected.GET("/user", controllers.GetUser)
	}

	return r
}
