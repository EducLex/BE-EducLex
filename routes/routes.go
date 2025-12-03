package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Konfigurasi CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5501"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
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

	r.POST("/auth/verify-email", controllers.VerifyEmail)

	r.POST("/auth/register-jaksa", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.CreateJaksa)

	// Pertanyaan dan Pengaduan
	r.GET("/questions", controllers.GetQuestions)
	r.POST("/questions", controllers.CreateQuestion)
	r.PUT("/questions/:id", middleware.AuthMiddleware(), controllers.UpdateQuestion)
	r.DELETE("/questions/:id", middleware.AuthMiddleware(), controllers.DeleteQuestion)
	r.POST("/questions/:id/diskusi", controllers.TambahDiskusi)
	r.GET("/:id/diskusi", controllers.GetDiskusiByQuestionID)

	// Artikel Routes
	r.GET("/articles", controllers.GetArticles)                                                               // Bisa diakses semua user
	r.POST("/articles", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.CreateArticle) // Hanya admin
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

	// Kategori routes
	r.POST("/categories", controllers.CreateCategory)
	r.GET("/categories", controllers.GetCategories)
	r.GET("/categories/:id", controllers.GetCategoryByID)
	r.PUT("/categories/:id", controllers.UpdateCategory)
	r.DELETE("/categories/:id", controllers.DeleteCategory)

	// Rute untuk Bidang
	r.POST("/bidang", controllers.CreateBidang)
	r.GET("/bidang", controllers.GetBidangs)
	r.GET("/bidang/:id", controllers.GetBidangByID)
	r.PUT("/bidang/:id", controllers.UpdateBidang)
	r.DELETE("/bidang/:id", controllers.DeleteBidang)

	// Tulisan Jaksa
	tulisan := r.Group("/tulisan")
	{
		tulisan.GET("/", controllers.GetAllTulisan)                                                             // Bisa diakses semua user
		tulisan.POST("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.CreateTulisan) // Hanya admin
		tulisan.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.UpdateTulisan)
		tulisan.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.DeleteTulisan)
	}

	r.GET("/tulisan/download/:id", controllers.DownloadFile)

	// Peraturan
	peraturan := r.Group("/peraturan")
	{
		peraturan.GET("", controllers.GetPeraturan)
		peraturan.GET(":id", controllers.GetPeraturanByID)

		peraturan.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.CreatePeraturan)
		peraturan.PUT(":id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.UpdatePeraturan)
		peraturan.DELETE(":id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controllers.DeletePeraturan)
	}

	//logout
	r.POST("/auth/logout", controllers.Logout)

	return r
}
