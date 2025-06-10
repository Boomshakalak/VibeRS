package dedup

import (
	"container/heap"

	"github.com/Boomshakalak/VibeRS/internal/store"
)

// Service handles deduplication using min-heap and bloom filters
type Service struct {
	seen map[int]bool // Simple dedup map (could be replaced with bloom filter)
}

// NewService creates a new deduplication service
func NewService() *Service {
	return &Service{
		seen: make(map[int]bool),
	}
}

// ItemHeap implements a min-heap of items ordered by score
type ItemHeap []ScoredItem

type ScoredItem struct {
	Item  store.Item
	Score float64
}

func (h ItemHeap) Len() int           { return len(h) }
func (h ItemHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h ItemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ItemHeap) Push(x interface{}) {
	*h = append(*h, x.(ScoredItem))
}

func (h *ItemHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// Deduplicate removes duplicate items and maintains top-N by score
func (s *Service) Deduplicate(items []store.Item) []store.Item {
	// Initialize heap
	h := &ItemHeap{}
	heap.Init(h)

	// Reset seen map for new deduplication session
	s.seen = make(map[int]bool)

	// Process items
	for _, item := range items {
		// Skip if already seen
		if s.seen[item.ItemID] {
			continue
		}
		s.seen[item.ItemID] = true

		// Calculate basic score (can be improved)
		score := s.calculateScore(item)

		// Add to heap
		heap.Push(h, ScoredItem{Item: item, Score: score})
	}

	// Extract all items from heap (sorted by score)
	var result []store.Item
	for h.Len() > 0 {
		scored := heap.Pop(h).(ScoredItem)
		result = append([]store.Item{scored.Item}, result...) // Prepend for desc order
	}

	return result
}

// calculateScore calculates a basic relevance score for an item
func (s *Service) calculateScore(item store.Item) float64 {
	// Simple scoring based on rating, GMV, and stock
	score := item.Rating * 0.3
	score += float64(item.GMV30d) / 10000.0 * 0.4
	score += float64(item.Click7d) / 100.0 * 0.2

	// Stock availability boost
	if item.Stock > 0 {
		score += 0.1
	}

	return score
}
