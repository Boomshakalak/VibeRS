package recall

import (
	"strings"
	"sync"

	"github.com/Boomshakalak/VibeRS/internal/store"
)

// Service handles parallel recall operations using specialized recallers
type Service struct {
	store        *store.Service
	textRecaller *TextRecaller
	attrRecaller *AttrRecaller
	hotRecaller  *HotRecaller
	expRecaller  *ExpRecaller
	annRecaller  *ANNRecaller
}

// NewService creates a new recall service with all specialized recallers
func NewService(storeService *store.Service) *Service {
	return &Service{
		store:        storeService,
		textRecaller: NewTextRecaller(storeService),
		attrRecaller: NewAttrRecaller(storeService),
		hotRecaller:  NewHotRecaller(storeService),
		expRecaller:  NewExpRecaller(storeService),
		annRecaller:  NewANNRecaller(storeService),
	}
}

// RecallResult represents the result from a single recall strategy
type RecallResult struct {
	Items  []store.Item
	Source string
	Score  float64
}

// ParallelRecall executes multiple recall strategies in parallel
func (s *Service) ParallelRecall(query string) ([]store.Item, error) {
	query = strings.TrimSpace(query)

	// If query is empty, return hot items only
	if query == "" {
		return s.hotRecaller.HotRecall(100)
	}

	// First, try text search
	textItems, err := s.textRecaller.MultiStrategyTextRecall(query, 1000)
	if err != nil {
		textItems = []store.Item{}
	}

	// If text search found good results, prioritize them
	if len(textItems) >= 1 {
		// We have good text results, add minimal diversity
		seen := make(map[int]bool)
		var allItems []store.Item

		// Add all text results
		for _, item := range textItems {
			allItems = append(allItems, item)
			seen[item.ItemID] = true
		}

		// Add a small amount of hot items for diversity (only if not already included)
		hotItems, err := s.hotRecaller.HotRecall(20)
		if err == nil {
			diversityCount := 0
			maxDiversity := 2 // Reduce diversity items
			if len(textItems) == 1 {
				maxDiversity = 0 // No diversity for single exact matches
			}
			for _, item := range hotItems {
				if !seen[item.ItemID] && diversityCount < maxDiversity {
					allItems = append(allItems, item)
					seen[item.ItemID] = true
					diversityCount++
				}
			}
		}

		return allItems, nil
	}

	// If text search results are insufficient, use parallel strategies
	var wg sync.WaitGroup
	results := make(chan RecallResult, 4)

	// Text recall
	wg.Add(1)
	go func() {
		defer wg.Done()
		results <- RecallResult{Items: textItems, Source: "text", Score: 1.0}
	}()

	// Hot items recall
	wg.Add(1)
	go func() {
		defer wg.Done()
		items, err := s.hotRecaller.HotRecall(300)
		if err != nil {
			results <- RecallResult{Items: []store.Item{}, Source: "hot", Score: 0.4}
		} else {
			results <- RecallResult{Items: items, Source: "hot", Score: 0.4}
		}
	}()

	// Exploration recall
	wg.Add(1)
	go func() {
		defer wg.Done()
		items, err := s.expRecaller.RandomRecall(200)
		if err != nil {
			results <- RecallResult{Items: []store.Item{}, Source: "explore", Score: 0.2}
		} else {
			results <- RecallResult{Items: items, Source: "explore", Score: 0.2}
		}
	}()

	// Attribute recall only if query might contain brand/filter info
	if s.attrRecaller.mightContainBrand(query) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			items, err := s.attrRecaller.SmartAttrRecall(query, 300)
			if err != nil {
				results <- RecallResult{Items: []store.Item{}, Source: "attr", Score: 0.6}
			} else {
				results <- RecallResult{Items: items, Source: "attr", Score: 0.6}
			}
		}()
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results with deduplication
	seen := make(map[int]bool)
	var allItems []store.Item

	for result := range results {
		for _, item := range result.Items {
			if !seen[item.ItemID] {
				allItems = append(allItems, item)
				seen[item.ItemID] = true
			}
		}
	}

	return allItems, nil
}

// GetTextRecaller returns the text recaller for direct access
func (s *Service) GetTextRecaller() *TextRecaller {
	return s.textRecaller
}

// GetAttrRecaller returns the attribute recaller for direct access
func (s *Service) GetAttrRecaller() *AttrRecaller {
	return s.attrRecaller
}

// GetHotRecaller returns the hot recaller for direct access
func (s *Service) GetHotRecaller() *HotRecaller {
	return s.hotRecaller
}

// GetExpRecaller returns the exploration recaller for direct access
func (s *Service) GetExpRecaller() *ExpRecaller {
	return s.expRecaller
}

// GetANNRecaller returns the ANN recaller for direct access
func (s *Service) GetANNRecaller() *ANNRecaller {
	return s.annRecaller
}
