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

	// CORS BENAR
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
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

	// Pertanyaan dan Pengaduan
	r.GET("/questions", controllers.GetQuestions)
	r.POST("/questions", controllers.CreateQuestion)
	r.PUT("/questions/:id", middleware.AuthMiddleware(), controllers.UpdateQuestion)
	r.DELETE("/questions/:id", middleware.AuthMiddleware(), controllers.DeleteQuestion)
	r.POST("/questions/:id/diskusi", controllers.TambahDiskusi)
	r.GET("/:id/diskusi", controllers.GetDiskusiByQuestionID)

	// Artikel Routes
	r.GET("/articles", controllers.GetArticles)
	r.GET("/articles/:id", controllers.GetArticleByID)
	r.POST("/articles", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.CreateArticle)
	r.PUT("/articles/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.UpdateArticle)
	r.DELETE("/articles/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.DeleteArticle)

	// ✅ Endpoint jaksa
	r.POST("/jaksa", controllers.CreateJaksa)
	r.GET("/jaksa", controllers.GetAllJaksa)
	r.PUT("/jaksa/:id", controllers.UpdateJaksa)
	r.DELETE("/jaksa/:id", controllers.DeleteJaksa)

	// PROFILE JAKSA
	r.GET("/jaksa/profile/:id", controllers.GetJaksaProfile)
	r.PUT("/jaksa/profile/:id", controllers.UpdateJaksaProfile)
	r.POST("/jaksa/auth/forgot-password", controllers.ForgotPassword)
	r.POST("/jaksa/auth/reset-password-jaksa", controllers.ResetPassword)

	// Tulisan Jaksa
	tulisan := r.Group("/tulisan")
	{
		tulisan.GET("/", controllers.GetAllTulisan)                                                             // bisa diakses semua user
		tulisan.POST("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.CreateTulisan) // cuma admin
		tulisan.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.UpdateTulisan)
		tulisan.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.DeleteTulisan)
	}

	// Peraturan (CRUD)
	peraturan := r.Group("/peraturan")
	{
		peraturan.GET("/", controllers.GetPeraturan)
		peraturan.GET("/:id", controllers.GetPeraturanByID)

		peraturan.POST("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.CreatePeraturan)
		peraturan.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.UpdatePeraturan)
		peraturan.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.DeletePeraturan)
	}

	return r

	//logout
	r.POST("/auth/logout", controllers.Logout)

	return r
}
