package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rijin/ads_manager/internal/api"
	"github.com/rijin/ads_manager/internal/campaign"
	"github.com/rijin/ads_manager/internal/engine"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading environment variables")
	}
	// Initialize components
	campaignManager := campaign.NewCampaignManager(10)
	adEngine := engine.NewAdEngine()

	// Initialize HTTP server
	server := api.NewServer(campaignManager, adEngine)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on :8080")
		if err := server.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")
}
