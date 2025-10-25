package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task_2/config"
	"task_2/initializers"
	"task_2/routes"

	"github.com/gin-gonic/gin"
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
	
	// HTTP server start up stuff...
	router := gin.Default()
	routes.SetupRoutes(router, db)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router)
	if err != nil {
		log.Println("Failed to start HTTP server because ", err.Error())
	}
	// Block until interrupt signal (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down")
}
