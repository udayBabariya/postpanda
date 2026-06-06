package main

import (
	"fmt"
	"log"
	"os"
	"postpanda/backend-go/internal/api"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/database"

	// Register all social providers
	_ "postpanda/backend-go/internal/integrations/social"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redis, err := database.ConnectRedis()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		defer redis.Close()
	}

	server := api.NewServer()
	server.SetupRoutes()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)

	if err := server.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
