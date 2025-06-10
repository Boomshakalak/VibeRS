package recall

import (
	"github.com/Boomshakalak/VibeRS/internal/store"
)

// HotRecaller handles hot item recall strategies
type HotRecaller struct {
	store *store.Service
}

// NewHotRecaller creates a new hot recall handler
func NewHotRecaller(storeService *store.Service) *HotRecaller {
	return &HotRecaller{store: storeService}
}

// HotRecall returns trending/hot items
// SQL: ORDER BY gmv_30d DESC, click_7d DESC
func (hr *HotRecaller) HotRecall(limit int) ([]store.Item, error) {
	return hr.store.GetHotItems(limit)
}

// GMVBasedRecall returns items sorted by GMV performance
func (hr *HotRecaller) GMVBasedRecall(limit int) ([]store.Item, error) {
	// Use the same logic as hot items for now
	// In production, this could have different sorting logic
	return hr.store.GetHotItems(limit)
}

// ClickBasedRecall returns items sorted by click performance
func (hr *HotRecaller) ClickBasedRecall(limit int) ([]store.Item, error) {
	// This would need a separate store method for click-based sorting
	// For now, using the existing hot items logic
	return hr.store.GetHotItems(limit)
}

// TrendingRecall returns items that are currently trending
func (hr *HotRecaller) TrendingRecall(limit int) ([]store.Item, error) {
	// This could combine multiple signals like:
	// - Recent click growth rate
	// - Recent purchase growth rate
	// - Stock movement velocity
	// For now, using the existing hot items logic
	return hr.store.GetHotItems(limit)
}

// RecentlyLaunchedRecall returns recently launched popular items
func (hr *HotRecaller) RecentlyLaunchedRecall(limit int) ([]store.Item, error) {
	// This would filter by recent launch date AND popularity
	// For now, using the existing hot items logic
	return hr.store.GetHotItems(limit)
}

// BrandPopularRecall returns popular items from specific brands
func (hr *HotRecaller) BrandPopularRecall(brands []string, limit int) ([]store.Item, error) {
	// This would filter hot items by brand
	// For now, using the existing hot items logic
	return hr.store.GetHotItems(limit)
}
