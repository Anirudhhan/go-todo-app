package main

import (
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

	srv := routes.SetupRoutes()

	if err := srv.Run(":8080"); err != nil {
		panic(err)
	}
}
