package campaign

import "sync"

// Manager handles campaign operations and bidding processes
type Manager struct {
	campaigns  map[string]*Campaign
	bidQueue   *PriorityQueue
	workerPool chan struct{}
	maxWorkers int
	mu         sync.RWMutex
}

// NewManager creates a new campaign manager instance
func NewManager(maxWorkers int) *Manager {
	return &Manager{
		campaigns:  make(map[string]*Campaign),
		bidQueue:   NewPriorityQueue(),
		workerPool: make(chan struct{}, maxWorkers),
		maxWorkers: maxWorkers,
	}
}

// NewPriorityQueue creates a new priority queue instance
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		queue: make(bidHeap, 0),
	}
	return pq
}
