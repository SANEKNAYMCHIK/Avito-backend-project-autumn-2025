package main

import (
	"log"

	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/db"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	defer database.Close()
}
