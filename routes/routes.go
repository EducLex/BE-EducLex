package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:8080"}, // alamat FE yg benar
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
}))

	// auth group
	auth := r.Group("/auth")
	{
		// manual login/register
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/register-admin", controllers.RegisterAdmin)

		// google login/register
		auth.GET("/google/login", controllers.GoogleLogin)
		auth.GET("/google/register", controllers.GoogleRegister)
		auth.GET("/google/callback", controllers.GoogleCallback)

		//qustion
		r.POST("/questions", controllers.CreateQuestion)
		r.GET("/questions", controllers.GetQuestions)

	}

	return r
}
