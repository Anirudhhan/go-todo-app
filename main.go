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

	srv := routes.SetUpRoutes()

	srv.Run()

	if err != nil {
		panic(err)
	}
}
