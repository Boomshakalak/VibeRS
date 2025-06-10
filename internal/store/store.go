package store

import (
	"database/sql"
	"encoding/binary"
	"math"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Item represents a product item in the database
type Item struct {
	ItemID     int       `json:"item_id"`
	Title      string    `json:"title"`
	Brand      string    `json:"brand"`
	PriceCents int       `json:"price_cents"`
	Discount   float64   `json:"discount"`
	Rating     float64   `json:"rating"`
	Stock      int       `json:"stock"`
	LaunchedAt time.Time `json:"launched_at"`
	Click7d    int       `json:"click_7d"`
	Buy7d      int       `json:"buy_7d"`
	GMV30d     int       `json:"gmv_30d"`
	Embedding  []float32 `json:"-"` // Hidden from JSON
}

// Service handles database operations
type Service struct {
	db *sql.DB
}

// NewService creates a new store service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// InitDB initializes the SQLite database with custom functions
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Register cosine similarity function
	if err := registerCosineSimilarity(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// registerCosineSimilarity registers a custom SQLite function for cosine similarity
func registerCosineSimilarity(db *sql.DB) error {
	cosineFunc := func(embedding1, embedding2 []byte) float64 {
		if len(embedding1) != len(embedding2) {
			return 0.0
		}

		vec1 := bytesToFloat32Slice(embedding1)
		vec2 := bytesToFloat32Slice(embedding2)

		if len(vec1) != len(vec2) {
			return 0.0
		}

		return cosineSimilarity(vec1, vec2)
	}

	// Note: This is a simplified version. In a real implementation,
	// you'd need to use a Go SQLite driver that supports custom functions
	// like modernc.org/sqlite or a CGO-enabled version
	_ = cosineFunc // TODO: Actually register the function

	return nil
}

// bytesToFloat32Slice converts byte slice to float32 slice
func bytesToFloat32Slice(data []byte) []float32 {
	if len(data)%4 != 0 {
		return nil
	}

	result := make([]float32, len(data)/4)
	for i := 0; i < len(result); i++ {
		bits := binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
		result[i] = math.Float32frombits(bits)
	}
	return result
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0.0 || normB == 0.0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// GetItemsByTextSearch performs simple text search (without FTS for now)
func (s *Service) GetItemsByTextSearch(query string, limit int) ([]Item, error) {
	sqlQuery := `
		SELECT item_id, title, brand, price_cents, discount, 
		       rating, stock, launched_at, click_7d, buy_7d, gmv_30d
		FROM items
		WHERE title LIKE '%' || ? || '%' OR brand LIKE '%' || ? || '%'
		ORDER BY rating DESC, gmv_30d DESC
		LIMIT ?
	`

	rows, err := s.db.Query(sqlQuery, query, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanItems(rows)
}

// GetItemsByFilter performs attribute-based filtering
func (s *Service) GetItemsByFilter(brand string, maxPrice int, minRating float64, limit int) ([]Item, error) {
	sqlQuery := `
		SELECT item_id, title, brand, price_cents, discount, 
		       rating, stock, launched_at, click_7d, buy_7d, gmv_30d
		FROM items
		WHERE ($1 = '' OR brand = $1)
		  AND ($2 = 0 OR price_cents <= $2)
		  AND rating >= $3
		  AND stock > 0
		ORDER BY rating DESC, gmv_30d DESC
		LIMIT $4
	`

	rows, err := s.db.Query(sqlQuery, brand, maxPrice, minRating, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanItems(rows)
}

// GetHotItems returns trending items
func (s *Service) GetHotItems(limit int) ([]Item, error) {
	sqlQuery := `
		SELECT item_id, title, brand, price_cents, discount, 
		       rating, stock, launched_at, click_7d, buy_7d, gmv_30d
		FROM items
		WHERE stock > 0
		ORDER BY gmv_30d DESC, click_7d DESC
		LIMIT ?
	`

	rows, err := s.db.Query(sqlQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanItems(rows)
}

// GetRandomItems returns random items for exploration
func (s *Service) GetRandomItems(limit int) ([]Item, error) {
	sqlQuery := `
		SELECT item_id, title, brand, price_cents, discount, 
		       rating, stock, launched_at, click_7d, buy_7d, gmv_30d
		FROM items
		WHERE stock > 0
		ORDER BY RANDOM()
		LIMIT ?
	`

	rows, err := s.db.Query(sqlQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanItems(rows)
}

// scanItems scans database rows into Item structs
func (s *Service) scanItems(rows *sql.Rows) ([]Item, error) {
	var items []Item

	for rows.Next() {
		var item Item
		var launchedAt sql.NullTime
		var click7d, buy7d, gmv30d sql.NullInt64

		err := rows.Scan(
			&item.ItemID, &item.Title, &item.Brand, &item.PriceCents,
			&item.Discount, &item.Rating, &item.Stock, &launchedAt,
			&click7d, &buy7d, &gmv30d,
		)
		if err != nil {
			return nil, err
		}

		// Handle NULL values
		if launchedAt.Valid {
			item.LaunchedAt = launchedAt.Time
		}
		if click7d.Valid {
			item.Click7d = int(click7d.Int64)
		}
		if buy7d.Valid {
			item.Buy7d = int(buy7d.Int64)
		}
		if gmv30d.Valid {
			item.GMV30d = int(gmv30d.Int64)
		}

		items = append(items, item)
	}

	return items, rows.Err()
}
