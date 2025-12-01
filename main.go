package main

import (
	"log"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/gin-contrib/cors"
	"github.com/EducLex/BE-EducLex/routes"
	"github.com/EducLex/BE-EducLex/controllers"
)

func main() {
	// Koneksi database
	config.ConnectDB()

	// Setup router (CORS sudah ada di dalam SetupRouter)
	r := routes.SetupRouter()

	// Seed kategori
	controllers.SeedCategories()

	// Aktifkan CORS
    r.Use(cors.Default()) 

	log.Println("Server running on :8080")
	r.Run(":8080")
}
