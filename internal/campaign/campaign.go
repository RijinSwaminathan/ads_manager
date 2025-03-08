package campaign

// CampaignManager handles ad campaign operations
type CampaignManager struct {
	maxCampaigns int
}

// NewCampaignManager creates a new campaign manager
func NewCampaignManager(maxCampaigns int) *CampaignManager {
	return &CampaignManager{
		maxCampaigns: maxCampaigns,
	}
}
