package main

import (
	"log"
	"net/http"

	"github.com/Boomshakalak/VibeRS/internal/dedup"
	"github.com/Boomshakalak/VibeRS/internal/rank/coarse"
	"github.com/Boomshakalak/VibeRS/internal/rank/final"
	"github.com/Boomshakalak/VibeRS/internal/rank/ltr"
	"github.com/Boomshakalak/VibeRS/internal/recall"
	"github.com/Boomshakalak/VibeRS/internal/store"
	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Query string `json:"q"`
	Page  int    `json:"page"`
}

type SearchResponse struct {
	Items   []store.Item `json:"items"`
	Total   int          `json:"total"`
	Page    int          `json:"page"`
	HasNext bool         `json:"has_next"`
}

func main() {
	// Initialize database
	db, err := store.InitDB("./data/vibers.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize services
	storeService := store.NewService(db)
	recallService := recall.NewService(storeService)
	dedupService := dedup.NewService()
	coarseRanker := coarse.NewRanker()
	ltrRanker := ltr.NewRanker()
	finalRanker := final.NewRanker()

	r := gin.Default()

	r.POST("/search", func(c *gin.Context) {
		var req SearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("JSON binding error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Search request: query='%s', page=%d", req.Query, req.Page)

		// Implement parallel recall
		items, err := recallService.ParallelRecall(req.Query)
		if err != nil {
			log.Printf("Parallel recall error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Parallel recall returned %d items", len(items))

		// TODO: Implement deduplication
		dedupedItems := dedupService.Deduplicate(items)

		// TODO: Implement three-stage ranking
		coarseRanked := coarseRanker.Rank(dedupedItems)
		ltrRanked := ltrRanker.Rank(coarseRanked)
		finalRanked := finalRanker.Rank(ltrRanked)

		// Pagination
		pageSize := 20
		start := (req.Page - 1) * pageSize
		end := start + pageSize
		if end > len(finalRanked) {
			end = len(finalRanked)
		}
		if start >= len(finalRanked) {
			start = len(finalRanked)
		}

		response := SearchResponse{
			Items:   finalRanked[start:end],
			Total:   len(finalRanked),
			Page:    req.Page,
			HasNext: end < len(finalRanked),
		}

		c.JSON(http.StatusOK, response)
	})

	log.Println("API server starting on :8080")
	r.Run(":8080")
}
