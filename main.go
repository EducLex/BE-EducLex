package main

import (
	"log"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/routes"
)

func main() {
	// Koneksi database
	config.ConnectDB()

	// Setup router (CORS sudah ada di dalam SetupRouter)
	r := routes.SetupRouter()

	log.Println("Server running on :8080")
	r.Run(":8080")
}
