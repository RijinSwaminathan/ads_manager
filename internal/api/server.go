package api

import (
	"net/http"

	"github.com/rijin/ads_manager/internal/campaign"
	"github.com/rijin/ads_manager/internal/engine"
)

// Server represents the HTTP server
type Server struct {
	campaignManager *campaign.CampaignManager
	adEngine       *engine.AdEngine
}

// NewServer creates a new HTTP server instance
func NewServer(cm *campaign.CampaignManager, ae *engine.AdEngine) *Server {
	return &Server{
		campaignManager: cm,
		adEngine:       ae,
	}
}

// Start begins listening for HTTP requests
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, nil)
}
