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
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"}, // izinkan FE
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},        // OPTIONS penting buat preflight
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
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
	}

	// question routes (dipisah dari auth biar rapi)
	r.POST("/questions", controllers.CreateQuestion)
	r.GET("/questions", controllers.GetQuestions)

	return r
}