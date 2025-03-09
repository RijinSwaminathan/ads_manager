package types

import "time"

// PredictionRequest represents a request for performance prediction
type PredictionRequest struct {
	CampaignID string
	Platform   Platform
	Timestamp  time.Time
}

// PredictionResponse represents predicted performance metrics
type PredictionResponse struct {
	PredictedCPC     float64
	PredictedCVR     float64
	Confidence       float64
	MarketCondition  MarketCondition
	MarketInsights   *MarketInsights
	AudienceInsights *AudienceInsights
	LastUpdated      time.Time
}

// PredictionModel represents a trained model for campaign performance prediction
type PredictionModel struct {
	ModelID     string
	Platform    Platform
	LastTrained time.Time
	Accuracy    float64
}
