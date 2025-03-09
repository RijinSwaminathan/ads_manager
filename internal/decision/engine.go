package decision

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/rijin/ads_manager/internal/campaign"
	"github.com/rijin/ads_manager/internal/predictive"
	"github.com/rijin/ads_manager/internal/types"
)

// Engine handles bid decisions and optimization strategies
type Engine struct {
	analyzer        *predictive.Analyzer
	cache           *decisionCache
	campaignManager *campaign.Manager
}

// decisionCache implements a thread-safe cache for recent decisions
type decisionCache struct {
	decisions map[string]*cachedDecision
	mu        sync.RWMutex
}

type cachedDecision struct {
	response   *types.BidResponse
	expiration time.Time
}

// NewEngine creates a new decision engine instance
func NewEngine(analyzer *predictive.Analyzer, campaignManager *campaign.Manager) *Engine {
	return &Engine{
		analyzer: analyzer,
		cache: &decisionCache{
			decisions: make(map[string]*cachedDecision),
		},
		campaignManager: campaignManager,
	}
}

// ProcessBidRequest processes a bid request and returns a bid response
func (e *Engine) ProcessBidRequest(ctx context.Context, req *types.BidRequest) (*types.BidResponse, error) {
	campaign, err := e.campaignManager.GetCampaign(req.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	// Get predictive metrics
	predictionReq := &types.PredictionRequest{
		CampaignID: campaign.ID,
		Platform:   req.Platform,
		Timestamp:  time.Now(),
	}
	prediction, err := e.analyzer.PredictPerformance(ctx, predictionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to predict metrics: %w", err)
	}

	// Calculate optimal bid
	bidAmount := e.calculateOptimalBid(campaign, prediction)
	if bidAmount <= 0 {
		return nil, nil
	}

	// Generate response
	response := &types.BidResponse{
		CampaignID: campaign.ID,
		BidAmount:  bidAmount,
		Confidence: prediction.Confidence,
		Strategy:   e.determineStrategy(prediction),
		Timestamp:  time.Now(),
	}

	// Cache the decision
	e.cacheDecision(req, response)

	return response, nil
}

// calculateOptimalBid determines the optimal bid amount using multiple factors
func (e *Engine) calculateOptimalBid(c *campaign.Campaign, pred *types.PredictionResponse) float64 {
	// Get maximum allowed bid based on remaining budget
	maxBid := e.calculateMaxBid(c)

	// Base bid calculation using predicted metrics
	baseBid := pred.PredictedCPC * (1 + pred.PredictedCVR)

	// Apply market condition multiplier
	marketMultiplier := e.getMarketMultiplier(pred.MarketCondition)
	adjustedBid := baseBid * marketMultiplier

	// Apply time-based multiplier
	timeMultiplier := e.getTimeBasedMultiplier(c)
	adjustedBid *= timeMultiplier

	// Apply confidence score
	adjustedBid *= (0.5 + pred.Confidence)

	// Ensure bid doesn't exceed max bid
	return math.Min(adjustedBid, maxBid)
}

// determineStrategy selects the optimal bidding strategy
func (e *Engine) determineStrategy(pred *types.PredictionResponse) types.BiddingStrategy {
	// Use market condition to determine base strategy
	switch pred.MarketCondition {
	case types.MarketVolatile:
		return types.StrategyConservative
	case types.MarketCompetitive:
		return types.StrategyAggressive
	default:
		return types.StrategyBalanced
	}
}

// getMarketMultiplier returns bid multiplier based on market conditions
func (e *Engine) getMarketMultiplier(condition types.MarketCondition) float64 {
	switch condition {
	case types.MarketVolatile:
		return 0.8 // Reduce bids in volatile markets
	case types.MarketCompetitive:
		return 1.2 // Increase bids in competitive markets
	case types.MarketQuiet:
		return 1.1 // Slightly increase bids in quiet markets
	default:
		return 1.0 // Normal market conditions
	}
}

// calculateMaxBid determines maximum allowed bid based on budget
func (e *Engine) calculateMaxBid(c *campaign.Campaign) float64 {
	// Never use more than 10% of remaining budget in a single bid
	return c.RemainingBudget * 0.1
}

// getTimeBasedMultiplier adjusts bids based on campaign timeline
func (e *Engine) getTimeBasedMultiplier(c *campaign.Campaign) float64 {
	totalDuration := c.EndTime.Sub(c.StartTime)
	elapsed := time.Since(c.StartTime)
	remaining := time.Until(time.Now())

	// Adjust bidding based on campaign progress
	if remaining < totalDuration/4 && c.RemainingBudget > (c.Budget*0.3) {
		return 1.2 // Accelerate spending if behind schedule
	}
	if elapsed < totalDuration/4 && c.RemainingBudget < (c.Budget*0.7) {
		return 0.8 // Slow down if spending too fast
	}
	return 1.0
}

// getPlatformConfig retrieves platform-specific configuration

// cacheDecision stores a decision in cache
func (e *Engine) cacheDecision(req *types.BidRequest, response *types.BidResponse) {
	e.cache.mu.Lock()
	defer e.cache.mu.Unlock()

	key := req.CampaignID + string(req.Platform)
	e.cache.decisions[key] = &cachedDecision{
		response:   response,
		expiration: time.Now().Add(5 * time.Second), // Cache for 5 seconds
	}
}
