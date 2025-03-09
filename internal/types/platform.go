package types

import (
	"sync"
	"time"
)

// Platform represents an advertising platform
type Platform string

const (
	PlatformGoogle Platform = "google"
	PlatformMeta   Platform = "meta"
)

// Status represents the status of a campaign
type Status string

const (
	StatusDraft   Status = "draft"
	StatusActive  Status = "active"
	StatusPaused  Status = "paused"
	StatusStopped Status = "stopped"
)

// BidRequest represents a bid request from a platform
type BidRequest struct {
	CampaignID  string
	Platform    Platform
	CurrentCPC  float64
	Timestamp   time.Time
	AuctionData map[string]interface{}
}

// BidResponse represents a bid response to a platform
type BidResponse struct {
	Success      bool
	CampaignID   string
	BidAmount    float64
	Strategy     BiddingStrategy
	Confidence   float64
	Timestamp    time.Time
	ErrorMessage string
}

// MarketTrend represents market trend direction
type MarketTrend string

const (
	TrendUp     MarketTrend = "up"
	TrendDown   MarketTrend = "down"
	TrendStable MarketTrend = "stable"
)

// CompetitionLevel represents market competition level
type CompetitionLevel string

const (
	CompetitionLow    CompetitionLevel = "low"
	CompetitionMedium CompetitionLevel = "medium"
	CompetitionHigh   CompetitionLevel = "high"
)

// MarketCondition represents the current state of the advertising market
type MarketCondition int

const (
	MarketNormal MarketCondition = iota
	MarketVolatile
	MarketCompetitive
	MarketQuiet
)

// BiddingStrategy represents different bidding strategies
type BiddingStrategy int

const (
	StrategyAggressive BiddingStrategy = iota
	StrategyBalanced
	StrategyConservative
)

// Performance represents campaign performance metrics
type Performance struct {
	Impressions     int64
	Clicks          int64
	Conversions     int64
	Spend           float64
	Revenue         float64
	AverageCPC      float64
	AverageCPA      float64
	CTR             float64
	CVR             float64
	ROAS            float64
	LastUpdated     time.Time
	MarketInsights  *MarketInsights
	MarketCondition MarketCondition
}

// MarketInsights represents market-related insights
type MarketInsights struct {
	AverageCPC       float64
	AverageCVR       float64
	MarketTrend      MarketTrend
	CompetitionLevel CompetitionLevel
	LastUpdated      time.Time
}

// Demographics represents demographic information
type Demographics struct {
	AgeGroups map[string]float64
	Gender    map[string]float64
}

// Interest represents an audience interest
type Interest struct {
	Name   string
	Weight float64
}

// Behavior represents an audience behavior
type Behavior struct {
	Name  string
	Score float64
}

// AudienceInsights represents audience-related insights
type AudienceInsights struct {
	Demographics Demographics
	Interests    []Interest
	Behaviors    []Behavior
	LastUpdated  time.Time
}

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	tokens     int
	maxTokens  int
	interval   time.Duration
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		interval:   interval,
		lastUpdate: time.Now(),
	}
}

// AcquireToken attempts to acquire a token from the rate limiter
func (r *RateLimiter) AcquireToken() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastUpdate)

	// Replenish tokens based on elapsed time
	tokensToAdd := int(elapsed.Seconds()) * r.maxTokens / int(r.interval.Seconds())
	r.tokens = min(r.tokens+tokensToAdd, r.maxTokens)
	r.lastUpdate = now

	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
