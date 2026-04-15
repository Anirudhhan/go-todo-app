package main

import "todo-app/database"

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
}
