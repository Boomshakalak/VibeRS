package recall

import (
	"strings"

	"github.com/Boomshakalak/VibeRS/internal/store"
)

// TextRecaller handles text-based recall strategies
type TextRecaller struct {
	store *store.Service
}

// NewTextRecaller creates a new text recall handler
func NewTextRecaller(storeService *store.Service) *TextRecaller {
	return &TextRecaller{store: storeService}
}

// FuzzyTextSearch performs advanced fuzzy text search
func (tr *TextRecaller) FuzzyTextSearch(query string, limit int) ([]store.Item, error) {
	// Use the improved text search from store
	return tr.store.GetItemsByTextSearch(query, limit)
}

// ExactSearch performs exact phrase matching
func (tr *TextRecaller) ExactSearch(query string, limit int) ([]store.Item, error) {
	// For exact search, we'll use the original query as-is
	return tr.store.GetItemsByTextSearch(query, limit)
}

// PrefixSearch performs prefix-based search (useful for autocomplete)
func (tr *TextRecaller) PrefixSearch(query string, limit int) ([]store.Item, error) {
	// This would be better implemented with FTS5, but for now use LIKE with prefix
	query = strings.TrimSpace(query)
	if query == "" {
		return []store.Item{}, nil
	}

	// Add prefix pattern
	prefixQuery := query + "%"
	return tr.store.GetItemsByPrefixSearch(prefixQuery, limit)
}

// BrandSearch performs brand-specific search
func (tr *TextRecaller) BrandSearch(brand string, limit int) ([]store.Item, error) {
	return tr.store.GetItemsByFilter(brand, 0, 0.0, limit)
}

// MultiStrategyTextRecall combines multiple text search strategies
func (tr *TextRecaller) MultiStrategyTextRecall(query string, limit int) ([]store.Item, error) {
	var allItems []store.Item
	seen := make(map[int]bool)

	// Strategy 1: Exact/fuzzy search (primary)
	fuzzyItems, err := tr.FuzzyTextSearch(query, limit)
	if err == nil {
		for _, item := range fuzzyItems {
			if !seen[item.ItemID] {
				allItems = append(allItems, item)
				seen[item.ItemID] = true
			}
		}
	}

	// If we found good results from fuzzy search, don't dilute with prefix search
	// Only use prefix search if we have very few results
	if len(allItems) < 3 {
		prefixItems, err := tr.PrefixSearch(query, limit-len(allItems))
		if err == nil {
			for _, item := range prefixItems {
				if !seen[item.ItemID] {
					allItems = append(allItems, item)
					seen[item.ItemID] = true
				}
			}
		}
	}

	// Limit results
	if len(allItems) > limit {
		allItems = allItems[:limit]
	}

	return allItems, nil
}
