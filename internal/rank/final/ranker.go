package final

import (
	"math"

	"github.com/Boomshakalak/VibeRS/internal/store"
)

// Ranker implements final ranking with business constraints
type Ranker struct {
	// Business policy parameters
	maxSameBrand    int     // Max items from same brand in top results
	diversityWeight float64 // Weight for diversity vs relevance
	newItemBoost    float64 // Boost for newly launched items
}

// NewRanker creates a new final ranker
func NewRanker() *Ranker {
	return &Ranker{
		maxSameBrand:    3,    // Max 3 items from same brand in top 20
		diversityWeight: 0.15, // 15% weight for diversity
		newItemBoost:    1.1,  // 10% boost for new items
	}
}

// Rank applies final business-aware ranking
func (r *Ranker) Rank(items []store.Item) []store.Item {
	if len(items) == 0 {
		return items
	}

	// Apply diversity-aware greedy selection
	return r.greedyDiversityRanking(items)
}

// greedyDiversityRanking implements a greedy algorithm for diversity-aware ranking
func (r *Ranker) greedyDiversityRanking(items []store.Item) []store.Item {
	if len(items) <= 1 {
		return items
	}

	result := make([]store.Item, 0, len(items))
	remaining := make([]store.Item, len(items))
	copy(remaining, items)

	brandCount := make(map[string]int)

	// Greedy selection with brand diversity constraint
	for len(remaining) > 0 {
		bestIdx := r.selectBestItem(remaining, brandCount)

		if bestIdx >= 0 {
			selected := remaining[bestIdx]
			result = append(result, selected)
			brandCount[selected.Brand]++

			// Remove selected item from remaining
			remaining = append(remaining[:bestIdx], remaining[bestIdx+1:]...)
		} else {
			// No valid item found (all brands exceeded limit), add remaining items
			result = append(result, remaining...)
			break
		}
	}

	return result
}

// selectBestItem selects the best item considering diversity constraints
func (r *Ranker) selectBestItem(items []store.Item, brandCount map[string]int) int {
	bestIdx := -1
	bestScore := -1.0

	for i, item := range items {
		// Check brand diversity constraint
		if brandCount[item.Brand] >= r.maxSameBrand {
			continue // Skip if brand limit exceeded
		}

		score := r.calculateFinalScore(item, brandCount)

		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	return bestIdx
}

// calculateFinalScore calculates the final ranking score with business adjustments
func (r *Ranker) calculateFinalScore(item store.Item, brandCount map[string]int) float64 {
	// Base score from LTR stage (simulated here)
	baseScore := r.simulateLTRScore(item)

	// Apply business adjustments
	finalScore := baseScore

	// New item boost (items launched in last 30 days)
	// In production, you'd calculate days since launch
	// For now, use a simple heuristic based on click count
	if item.Click7d > 0 && item.Click7d < 50 { // Assume new items have low clicks
		finalScore *= r.newItemBoost
	}

	// Brand diversity penalty
	currentBrandCount := brandCount[item.Brand]
	if currentBrandCount > 0 {
		diversityPenalty := 1.0 - (r.diversityWeight * float64(currentBrandCount))
		finalScore *= math.Max(diversityPenalty, 0.5) // Min 50% of original score
	}

	// GMV optimization (prioritize high-value items)
	gmvBoost := 1.0 + (float64(item.GMV30d) / 10000000.0 * 0.1)
	finalScore *= gmvBoost

	// Stock urgency (slightly prioritize low stock items)
	if item.Stock <= 3 && item.Stock > 0 {
		finalScore *= 1.05 // 5% boost for low stock
	}

	return finalScore
}

// simulateLTRScore simulates the output from LTR stage
func (r *Ranker) simulateLTRScore(item store.Item) float64 {
	// This would normally come from the LTR stage
	// For simulation, use a combination of signals
	score := item.Rating / 5.0 * 40.0                // Rating component
	score += float64(item.GMV30d) / 1000000.0 * 30.0 // GMV component
	score += float64(item.Click7d) / 100.0 * 20.0    // CTR component
	score += (1.0 - item.Discount) * 10.0            // Discount component

	return score
}

// TODO: Implement LP-based re-ranking for more sophisticated optimization
// func (r *Ranker) lpReranking(items []store.Item) []store.Item {
//     // Linear programming based re-ranking for complex business constraints
//     // This could optimize for multiple objectives simultaneously:
//     // - Maximize relevance
//     // - Ensure brand diversity
//     // - Optimize for GMV
//     // - Satisfy business rules
//     return items
// }
