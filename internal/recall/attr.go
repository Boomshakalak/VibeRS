package recall

import (
	"strings"

	"github.com/Boomshakalak/VibeRS/internal/store"
)

// AttrRecaller handles attribute-based recall strategies
type AttrRecaller struct {
	store *store.Service
}

// NewAttrRecaller creates a new attribute recall handler
func NewAttrRecaller(storeService *store.Service) *AttrRecaller {
	return &AttrRecaller{store: storeService}
}

// FilterRecall performs attribute-based filtering
// SQL: brand=? AND price_cents<? AND rating>=?
func (ar *AttrRecaller) FilterRecall(brand string, maxPrice int, minRating float64, limit int) ([]store.Item, error) {
	return ar.store.GetItemsByFilter(brand, maxPrice, minRating, limit)
}

// BrandRecall performs brand-specific recall
func (ar *AttrRecaller) BrandRecall(brand string, limit int) ([]store.Item, error) {
	return ar.store.GetItemsByFilter(brand, 0, 0.0, limit)
}

// PriceRangeRecall performs price-based filtering
func (ar *AttrRecaller) PriceRangeRecall(minPrice, maxPrice int, limit int) ([]store.Item, error) {
	return ar.store.GetItemsByFilter("", maxPrice, 0.0, limit)
}

// RatingRecall performs rating-based filtering
func (ar *AttrRecaller) RatingRecall(minRating float64, limit int) ([]store.Item, error) {
	return ar.store.GetItemsByFilter("", 0, minRating, limit)
}

// SmartAttrRecall extracts attributes from query and performs smart filtering
func (ar *AttrRecaller) SmartAttrRecall(query string, limit int) ([]store.Item, error) {
	// Extract brand from query
	brand := ar.extractBrandFromQuery(query)

	// Extract price hints (future enhancement)
	// Extract rating hints (future enhancement)

	if brand != "" {
		return ar.BrandRecall(brand, limit)
	}

	// Default to no filter if no attributes found
	return []store.Item{}, nil
}

// extractBrandFromQuery attempts to extract brand name from query
func (ar *AttrRecaller) extractBrandFromQuery(query string) string {
	brands := map[string]string{
		"gucci":          "Gucci",
		"louis vuitton":  "Louis Vuitton",
		"chanel":         "Chanel",
		"hermes":         "Hermès",
		"prada":          "Prada",
		"saint laurent":  "Saint Laurent",
		"bottega veneta": "Bottega Veneta",
		"fendi":          "Fendi",
		"dior":           "Dior",
		"balenciaga":     "Balenciaga",
		"celine":         "Celine",
		"givenchy":       "Givenchy",
		"valentino":      "Valentino",
		"loewe":          "Loewe",
		"jacquemus":      "Jacquemus",
		"staud":          "Staud",
		"mansur gavriel": "Mansur Gavriel",
		"cult gaia":      "Cult Gaia",
		"polene":         "Polene",
		"wandler":        "Wandler",
	}

	queryLower := strings.ToLower(query)
	for key, brand := range brands {
		if strings.Contains(queryLower, key) {
			return brand
		}
	}
	return ""
}

// mightContainBrand checks if query might contain brand information
func (ar *AttrRecaller) mightContainBrand(query string) bool {
	brands := []string{"gucci", "louis vuitton", "chanel", "hermes", "prada",
		"saint laurent", "bottega veneta", "fendi", "dior", "balenciaga",
		"celine", "givenchy", "valentino", "loewe", "jacquemus", "staud",
		"mansur gavriel", "cult gaia", "polene", "wandler"}

	queryLower := strings.ToLower(query)
	for _, brand := range brands {
		if strings.Contains(queryLower, brand) {
			return true
		}
	}
	return false
}
