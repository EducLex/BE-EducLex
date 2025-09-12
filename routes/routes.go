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
		auth.GET("/google/login", controllers.GoogleLogin)
		auth.GET("/google/callback", controllers.GoogleCallback)

	}

	// Tambahkan route google-success
	r.GET("/google-success", func(c *gin.Context) {
		token := c.Query("token")
		c.Header("Content-Type", "text/html")
		c.String(200, "<h1>Login Success!</h1><p>Token: "+token+"</p>")
	})

	return r
}
