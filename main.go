package main

import (
	"log"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/routes"
	"github.com/gin-contrib/cors"
)

func main() {
	// koneksi DB
	config.ConnectDB()

	// setup router
	r := routes.SetupRouter()
	log.Println("Server running on :8080")

	// CORS HARUS DIPASANG DI SINI SEBELUM ROUTES DIGUNAKAN
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// run server
	r.Run(":8080")

}
