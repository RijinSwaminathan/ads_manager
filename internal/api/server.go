package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rijin/ads_manager/internal/campaign"
	"github.com/rijin/ads_manager/internal/decision"
	"github.com/rijin/ads_manager/internal/platform"
	"github.com/rijin/ads_manager/internal/types"
)

type Server struct {
	router          *mux.Router
	httpServer      *http.Server
	campaignManager *campaign.Manager
	engine          *decision.Engine
	googleAdsClient *platform.GoogleAdsClient
}

func NewServer(campaignManager *campaign.Manager, engine *decision.Engine) *Server {
	router := mux.NewRouter()
	s := &Server{
		router:          router,
		campaignManager: campaignManager,
		engine:          engine,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/api/campaigns", s.handleCreateCampaign).Methods("POST")
	s.router.HandleFunc("/api/campaigns/{id}", s.handleGetCampaign).Methods("GET")
	s.router.HandleFunc("/api/campaigns/{id}/start", s.handleStartCampaign).Methods("POST")
	s.router.HandleFunc("/api/campaigns/{id}/stop", s.handleStopCampaign).Methods("POST")
	s.router.HandleFunc("/api/campaigns/{id}/metrics", s.handleGetMetrics).Methods("GET")
}

func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

type CreateCampaignRequest struct {
	Budget      float64          `json:"budget"`
	TargetReach int              `json:"target_reach"`
	TargetCPA   float64          `json:"target_cpa"`
	Platforms   []types.Platform `json:"platforms"`
	StartTime   time.Time        `json:"start_time"`
	EndTime     time.Time        `json:"end_time"`
}

func (s *Server) handleCreateCampaign(w http.ResponseWriter, r *http.Request) {
	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	campaign, err := s.campaignManager.CreateCampaign(
		r.Context(),
		req.Budget,
		req.TargetReach,
		req.TargetCPA,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create campaign: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

func (s *Server) handleGetCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignID := vars["id"]

	campaign, err := s.campaignManager.GetCampaign(campaignID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

func (s *Server) handleStartCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignID := vars["id"]

	if err := s.campaignManager.UpdateCampaignStatus(campaignID, types.StatusActive); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleStopCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignID := vars["id"]

	if err := s.campaignManager.UpdateCampaignStatus(campaignID, types.StatusStopped); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignID := vars["id"]

	campaign, err := s.campaignManager.GetCampaign(campaignID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get performance metrics from Google Ads
	ctx := r.Context()
	for _, platformConfig := range campaign.Platforms {
		if platformConfig != types.PlatformGoogle {
			continue
		}

		performance, err := s.googleAdsClient.GetPerformanceMetrics(ctx, campaignID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get metrics from Google Ads: %v", err), http.StatusInternalServerError)
			return
		}

		if performance != nil {
			totalImpressions := performance.Impressions
			totalClicks := performance.Clicks
			totalConversions := performance.Conversions
			totalSpend := performance.Spend
			lastUpdated := performance.LastUpdated

			// Calculate derived metrics
			var ctr float64
			var cvr float64
			var cpc float64
			var cpa float64

			if totalImpressions > 0 {
				ctr = float64(totalClicks) / float64(totalImpressions)
			}
			if totalClicks > 0 {
				cvr = float64(totalConversions) / float64(totalClicks)
				cpc = totalSpend / float64(totalClicks)
			}
			if totalConversions > 0 {
				cpa = totalSpend / float64(totalConversions)
			}

			metrics := map[string]interface{}{
				"campaign_id":      campaignID,
				"status":           campaign.Status,
				"impressions":      totalImpressions,
				"clicks":           totalClicks,
				"conversions":      totalConversions,
				"spend":            totalSpend,
				"ctr":              ctr,
				"cvr":              cvr,
				"cpc":              cpc,
				"cpa":              cpa,
				"target_cpa":       campaign.TargetCPA,
				"budget":           campaign.Budget,
				"remaining_budget": campaign.RemainingBudget,
				"last_updated":     lastUpdated,
			}

			json.NewEncoder(w).Encode(metrics)
		}
	}
}
