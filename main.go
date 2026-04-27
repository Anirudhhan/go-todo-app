package main

import (
	"os"
	"todo-app/database"
	"todo-app/routes"
)

func main() {

	err := database.ConnectAndMigrate(
		"localhost",
		"5432",
		"todo",
		"local",
		"local",
		database.SSLMode(database.SSLModeDisable),
	)

	if err != nil {
		panic(err)
	}

	secret, exists := os.LookupEnv("ACCESS_SECRET")
	if !exists || secret == "" {
		panic("ACCESS_SECRET is not set")
	}

	srv := routes.SetupRoutes()

	if err := srv.Run(":8080"); err != nil {
		panic(err)
	}
}
