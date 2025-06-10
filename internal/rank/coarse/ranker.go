package coarse

import (
	"github.com/Boomshakalak/VibeRS/internal/store"
)

// Ranker implements coarse ranking with hard rules
type Ranker struct {
	// Configuration could be loaded from config file
	minStock      int
	maxPriceCents int
	minRating     float64
}

// NewRanker creates a new coarse ranker
func NewRanker() *Ranker {
	return &Ranker{
		minStock:      1,       // At least 1 in stock
		maxPriceCents: 2000000, // Max 20k USD
		minRating:     3.0,     // Min 3.0 rating
	}
}

// Rank applies hard filtering rules and basic scoring
func (r *Ranker) Rank(items []store.Item) []store.Item {
	var filtered []store.Item

	// Apply hard rules
	for _, item := range items {
		if r.passesHardRules(item) {
			filtered = append(filtered, item)
		}
	}

	// Apply coarse scoring and sorting
	scored := r.scoreItems(filtered)

	return scored
}

// passesHardRules checks if an item passes all hard filtering rules
func (r *Ranker) passesHardRules(item store.Item) bool {
	// Stock check
	if item.Stock < r.minStock {
		return false
	}

	// Price check
	if item.PriceCents > r.maxPriceCents {
		return false
	}

	// Rating check
	if item.Rating < r.minRating {
		return false
	}

	// Additional business rules can be added here
	return true
}

// scoreItems applies coarse scoring logic
func (r *Ranker) scoreItems(items []store.Item) []store.Item {
	// For now, just return items sorted by a simple formula
	// In production, this could use a more sophisticated rule engine

	// Simple bubble sort by composite score (for demonstration)
	// In production, use sort.Slice
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			score1 := r.calculateCoarseScore(items[i])
			score2 := r.calculateCoarseScore(items[j])

			if score1 < score2 {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	return items
}

// calculateCoarseScore calculates a coarse relevance score
func (r *Ranker) calculateCoarseScore(item store.Item) float64 {
	score := 0.0

	// Rating component (0-5 scale)
	score += item.Rating * 20.0

	// GMV component (normalized)
	score += float64(item.GMV30d) / 100000.0

	// Stock availability bonus
	if item.Stock > 5 {
		score += 10.0
	} else if item.Stock > 0 {
		score += 5.0
	}

	// Discount penalty (high discount might indicate poor quality)
	if item.Discount > 0.3 {
		score -= 5.0
	}

	// Click-through rate bonus
	score += float64(item.Click7d) / 10.0

	return score
}
