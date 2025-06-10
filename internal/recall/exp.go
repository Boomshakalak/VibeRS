package recall

import (
	"github.com/Boomshakalak/VibeRS/internal/store"
)

// ExpRecaller handles exploration recall strategies
type ExpRecaller struct {
	store *store.Service
}

// NewExpRecaller creates a new exploration recall handler
func NewExpRecaller(storeService *store.Service) *ExpRecaller {
	return &ExpRecaller{store: storeService}
}

// RandomRecall returns random items for exploration
// SQL: ORDER BY RANDOM() LIMIT ?
func (er *ExpRecaller) RandomRecall(limit int) ([]store.Item, error) {
	return er.store.GetRandomItems(limit)
}

// DiversityRecall returns diverse items to increase exploration
func (er *ExpRecaller) DiversityRecall(limit int) ([]store.Item, error) {
	// This could implement brand diversity, price diversity, etc.
	// For now, using random selection which provides natural diversity
	return er.store.GetRandomItems(limit)
}

// LongTailRecall returns less popular items for discovery
func (er *ExpRecaller) LongTailRecall(limit int) ([]store.Item, error) {
	// This would return items with lower popularity scores
	// Could be inverse of hot items - items with lower clicks/GMV
	// For now, using random which includes long-tail items
	return er.store.GetRandomItems(limit)
}

// SerendipityRecall returns unexpected but potentially interesting items
func (er *ExpRecaller) SerendipityRecall(limit int) ([]store.Item, error) {
	// This could use collaborative filtering or content-based surprises
	// For now, using random selection
	return er.store.GetRandomItems(limit)
}

// NewItemsRecall returns newly added items for discovery
func (er *ExpRecaller) NewItemsRecall(limit int) ([]store.Item, error) {
	// This would sort by recently added items (launched_at DESC)
	// For now, using random selection
	return er.store.GetRandomItems(limit)
}

// BudgetFriendlyRecall returns lower-priced items for exploration
func (er *ExpRecaller) BudgetFriendlyRecall(limit int) ([]store.Item, error) {
	// This would filter by lower price ranges
	// For now, using random selection
	return er.store.GetRandomItems(limit)
}

// UnderTheRadarRecall returns items that might be overlooked
func (er *ExpRecaller) UnderTheRadarRecall(limit int) ([]store.Item, error) {
	// This could find items with good ratings but low clicks
	// Items that are quality but not well-discovered
	// For now, using random selection
	return er.store.GetRandomItems(limit)
}
