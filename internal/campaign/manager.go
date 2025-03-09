package campaign

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rijin/ads_manager/internal/types"
)

// Manager handles campaign operations and bidding processes
type Manager struct {
	campaigns map[string]*Campaign
	mu        sync.RWMutex
}

// NewManager creates a new campaign manager
func NewManager() *Manager {
	return &Manager{
		campaigns: make(map[string]*Campaign),
	}
}

// CreateCampaign initializes a new campaign with the given parameters
func (m *Manager) CreateCampaign(ctx context.Context, budget float64, targetReach int, targetCPA float64, platform []types.Platform, startTime, endTime time.Time) (*Campaign, error) {
	if budget <= 0 || targetReach <= 0 || targetCPA <= 0 {
		return nil, errors.New("invalid campaign parameters")
	}

	campaign := &Campaign{
		ID:              uuid.New().String(),
		Name:            "Name of the campaign",
		Budget:          budget,
		RemainingBudget: budget,
		TargetReach:     targetReach,
		TargetCPA:       targetCPA,
		Platforms:       platform[:],
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          types.StatusDraft,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	m.mu.Lock()
	m.campaigns[campaign.ID] = campaign
	m.mu.Unlock()

	return campaign, nil
}

// GetCampaign retrieves a campaign by its ID
func (m *Manager) GetCampaign(id string) (*Campaign, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	campaign, exists := m.campaigns[id]
	if !exists {
		return nil, errors.New("campaign not found")
	}
	return campaign, nil
}

// UpdateCampaign updates an existing campaign
func (m *Manager) UpdateCampaign(ctx context.Context, campaign *Campaign) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.campaigns[campaign.ID]; !exists {
		return errors.New("campaign not found")
	}

	campaign.UpdatedAt = time.Now()
	m.campaigns[campaign.ID] = campaign
	return nil
}

// DeleteCampaign removes a campaign
func (m *Manager) DeleteCampaign(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.campaigns[id]; !exists {
		return errors.New("campaign not found")
	}

	delete(m.campaigns, id)
	return nil
}

// UpdateCampaignStatus updates the status of a campaign
func (m *Manager) UpdateCampaignStatus(id string, status types.Status) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	campaign, exists := m.campaigns[id]
	if !exists {
		return errors.New("campaign not found")
	}

	campaign.Status = status
	campaign.UpdatedAt = time.Now()
	return nil
}

// ProcessBid processes a bid for a campaign
func (m *Manager) ProcessBid(ctx context.Context, campaignID string, bidAmount float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	campaign, exists := m.campaigns[campaignID]
	if !exists {
		return errors.New("campaign not found")
	}

	if campaign.Status != types.StatusActive {
		return errors.New("campaign is not active")
	}

	if bidAmount > campaign.RemainingBudget {
		return errors.New("insufficient budget")
	}

	campaign.RemainingBudget -= bidAmount
	campaign.UpdatedAt = time.Now()
	return nil
}
