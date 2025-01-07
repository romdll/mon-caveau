package main

import (
	"moncaveau/database"
	"moncaveau/database/injector"
	"moncaveau/server"
	"moncaveau/utils"
)

var (
	logger = utils.CreateLogger("main")
)

func main() {
	logger.Info("Main called - Starting everything")

	if err := database.InitDB(); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	if err := database.ApplyMigrations(); err != nil {
		logger.Fatalf("Failed to apply migrations: %v", err)
	}

	if err := injector.SetupAndInjectAll(); err != nil {
		logger.Fatalf("Failed to inject basic data: %v", err)
	}

	serverEngine := server.CreateServer()
	utils.RunWithQuitNotification(serverEngine)
}
