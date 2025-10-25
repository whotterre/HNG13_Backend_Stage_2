package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"task_2/config"
	initializers "task_2/initializer"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Obtain MySQL connection
	db, err := initializers.ConnectToDB(cfg.DBString)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL db: %v", err)
	}

	// Perform automatic migration
	if err := initializers.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to perform database migrations: %v", err)
	}

	log.Println("Database connected and migrations applied")

	// TODO: start your HTTP server / router here.

	// Block until interrupt signal (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down")
}
