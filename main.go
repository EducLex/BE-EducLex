package main

import (
	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/routes"
)

func main() {
	// koneksi DB
	config.ConnectDB()

	// setup router
	r := routes.SetupRouter()

	// run server
	r.Run(":8080")
}
