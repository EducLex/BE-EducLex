package routes

import (
	"time"

	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
		auth.GET("/google/callback", controllers.GoogleCallback)
	}

	auth.GET("/user", middleware.AuthMiddleware(), controllers.GetUser)

	// hanya admin
	auth.PUT("/update-role", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.UpdateRole)

	auth.GET("/profile", middleware.AuthMiddleware(), controllers.ProfileHandler)

		// Dashboard Admin
	r.GET("/dashboard", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.GetDashboardStats)

	// Data Pengguna
	r.GET("/users", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.GetAllUsers)


	// question routes (dipisah dari auth biar rapi)
	r.POST("/questions", controllers.CreateQuestion)
	r.GET("/questions", controllers.GetQuestions)

	//artikel
	r.GET("/articles", controllers.GetArticles)
	r.GET("/articles/:id", controllers.GetArticleByID)
	r.POST("/articles", controllers.CreateArticle)
	r.PUT("/articles/:id", controllers.UpdateArticle)
	r.DELETE("/articles/:id", controllers.DeleteArticle)

	// Tulisan Jaksa
	r.POST("/tulisan", controllers.CreateTulisan)
	r.GET("/tulisan", controllers.GetTulisans)

	//Peraturan
	r.POST("/peraturan", controllers.CreatePeraturan)
	r.GET("/peraturan", controllers.GetPeraturan)
	r.DELETE("/peraturan/:id", controllers.DeletePeraturan)

	//logout
	r.POST("/auth/logout", controllers.Logout)

	

	return r
}
