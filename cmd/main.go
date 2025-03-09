package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rijin/ads_manager/internal/api"
	"github.com/rijin/ads_manager/internal/campaign"
	"github.com/rijin/ads_manager/internal/decision"
	"github.com/rijin/ads_manager/internal/messaging"
	"github.com/rijin/ads_manager/internal/predictive"
	"github.com/rijin/ads_manager/internal/types"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize components
	analyzer := predictive.NewAnalyzer()
	campaignManager := campaign.NewManager()
	engine := decision.NewEngine(analyzer, campaignManager)
	msgClient := messaging.NewClient()

	// Initialize HTTP server
	server := api.NewServer(campaignManager, engine)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on :8090")
		if err := server.Start(":8090"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Subscribe to bid responses
	bidResponseCh, err := msgClient.Subscribe("bid_responses")
	if err != nil {
		log.Fatalf("Failed to subscribe to bid responses: %v", err)
	}
	defer msgClient.Close()

	// Start processing bid responses in background
	go func() {
		for msg := range bidResponseCh {
			if msg.Type == messaging.BidResponseMessage {
				resp := msg.Payload.(*types.BidResponse)
				log.Printf("Received bid response for campaign %s: $%.2f", resp.CampaignID, resp.BidAmount)
			}
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
