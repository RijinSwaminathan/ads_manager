package campaign

import (
	"sync"
	"time"
)

// Platform represents an advertising platform (e.g., Meta, Google)
type Platform string

const (
	PlatformGoogle Platform = "GOOGLE"
	PlatformMeta   Platform = "META"
)

// Campaign represents an advertising campaign with its parameters and constraints
type Campaign struct {
	ID              string
	Budget          float64
	RemainingBudget float64
	TargetReach     int
	TargetCPA       float64 // Target Cost Per Acquisition
	StartTime       time.Time
	EndTime         time.Time
	Platforms       []PlatformConfig
	Status          CampaignStatus
	mu              sync.RWMutex
}

// BidRequest represents a real-time bidding request
type BidRequest struct {
	CampaignID      string
	Platform        Platform
	Timestamp       time.Time
	CurrentCPC      float64
	PredictedCVR    float64
	AuctionDeadline time.Duration
}

// PlatformConfig holds platform-specific configuration and performance metrics
type PlatformConfig struct {
	Platform    Platform
	CurrentCPC  float64 // Current Cost Per Click
	CurrentCVR  float64 // Current Conversion Rate
	LastUpdated time.Time
	HistorySize int
	CPCHistory  []float64 // Historical CPC data
	CVRHistory  []float64 // Historical CVR data
	Performance Performance
}

// CampaignStatus represents the current state of a campaign
type CampaignStatus int

// Performance tracks real-time campaign performance metrics
type Performance struct {
	Impressions int64
	Clicks      int64
	Conversions int64
	Spend       float64
	LastUpdated time.Time
}
