package campaign

import (
	"container/heap"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rijin/ads_manager/internal/types"
)

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInsufficientBudget = errors.New("insufficient budget")
)

// Campaign represents an advertising campaign
type Campaign struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Budget          float64            `json:"budget"`
	RemainingBudget float64            `json:"remaining_budget"`
	TargetReach     int                `json:"target_reach"`
	TargetCPA       float64            `json:"target_cpa"`
	StartTime       time.Time          `json:"start_time"`
	EndTime         time.Time          `json:"end_time"`
	Platforms       []types.Platform   `json:"platforms"`
	Status          types.Status       `json:"status"`
	Performance     *types.Performance `json:"performance,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

// NewCampaign creates a new campaign
func NewCampaign(name string, budget float64, targetReach int, targetCPA float64, startTime, endTime time.Time, platforms []types.Platform) *Campaign {
	return &Campaign{
		ID:              uuid.New().String(),
		Name:            name,
		Budget:          budget,
		RemainingBudget: budget,
		TargetReach:     targetReach,
		TargetCPA:       targetCPA,
		StartTime:       startTime,
		EndTime:         endTime,
		Platforms:       []types.Platform{},
		Status:          types.StatusDraft,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// BidItem represents a bid in the queue
type BidItem struct {
	CampaignID string
	BidAmount  float64
	Priority   float64
	Index      int
}

// BidQueue represents a priority queue for bids
type BidQueue struct {
	queue []*BidItem
}

// CreateBidQueue creates a new bid queue
func CreateBidQueue() *BidQueue {
	pq := &BidQueue{
		queue: make([]*BidItem, 0),
	}
	heap.Init((*bidQueueImpl)(&pq.queue))
	return pq
}

// Push adds a bid to the queue
func (pq *BidQueue) Push(item *BidItem) {
	heap.Push((*bidQueueImpl)(&pq.queue), item)
}

// Pop removes and returns the highest priority bid
func (pq *BidQueue) Pop() *BidItem {
	if len(pq.queue) == 0 {
		return nil
	}
	return heap.Pop((*bidQueueImpl)(&pq.queue)).(*BidItem)
}

// Len returns the number of items in the queue
func (pq *BidQueue) Len() int {
	return len(pq.queue)
}

// bidQueueImpl implements heap.Interface for priority queue
type bidQueueImpl []*BidItem

func (pq bidQueueImpl) Len() int {
	return len(pq)
}

func (pq bidQueueImpl) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq bidQueueImpl) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *bidQueueImpl) Push(x interface{}) {
	n := len(*pq)
	item := x.(*BidItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *bidQueueImpl) Pop() interface{} {
	old := *pq
	n := len(old)
	if n == 0 {
		return nil
	}
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

// BidRequest represents a bid request for a campaign
type BidRequest struct {
	CampaignID  string                 `json:"campaign_id"`
	Platform    types.Platform         `json:"platform"`
	CurrentCPC  float64                `json:"current_cpc"`
	Timestamp   time.Time              `json:"timestamp"`
	AuctionData map[string]interface{} `json:"auction_data"`
}

// BidResponse represents a bid response for a campaign
type BidResponse struct {
	Success      bool                  `json:"success"`
	CampaignID   string                `json:"campaign_id"`
	BidAmount    float64               `json:"bid_amount"`
	Strategy     types.BiddingStrategy `json:"strategy"`
	Confidence   float64               `json:"confidence"`
	Timestamp    time.Time             `json:"timestamp"`
	ErrorMessage string                `json:"error_message,omitempty"`
}

// Performance represents campaign performance metrics
type Performance struct {
	Impressions     int64                 `json:"impressions"`
	Clicks          int64                 `json:"clicks"`
	Conversions     int64                 `json:"conversions"`
	Spend           float64               `json:"spend"`
	Revenue         float64               `json:"revenue"`
	AverageCPC      float64               `json:"average_cpc"`
	AverageCPA      float64               `json:"average_cpa"`
	CTR             float64               `json:"ctr"`
	CVR             float64               `json:"cvr"`
	ROAS            float64               `json:"roas"`
	LastUpdated     time.Time             `json:"last_updated"`
	MarketInsights  *types.MarketInsights `json:"market_insights,omitempty"`
	MarketCondition types.MarketCondition `json:"market_condition"`
}

// Strategy constants for bidding
const (
	StrategyBalanced     = types.StrategyBalanced
	StrategyAggressive   = types.StrategyAggressive
	StrategyConservative = types.StrategyConservative
)

// PredictionModel represents ML model outputs for bid optimization
type PredictionModel struct {
	PredictedCPC      float64                `json:"predicted_cpc"`
	PredictedCVR      float64                `json:"predicted_cvr"`
	PredictedCPA      float64                `json:"predicted_cpa"`
	PredictedReach    int                    `json:"predicted_reach"`
	PredictedBudget   float64                `json:"predicted_budget"`
	ConfidenceScore   float64                `json:"confidence_score"`
	ModelVersion      string                 `json:"model_version"`
	UpdatedAt         time.Time              `json:"updated_at"`
	AdditionalMetrics map[string]interface{} `json:"additional_metrics,omitempty"`
}
