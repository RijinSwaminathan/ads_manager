package campaign

import (
	"sync"
)

// PriorityQueue implements a thread-safe priority queue for bid requests
type PriorityQueue struct {
	queue bidHeap
	mu    sync.RWMutex
}

// bidHeap implements heap.Interface
type bidHeap []*BidRequest
