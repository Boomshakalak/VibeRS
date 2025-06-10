package recall

import (
	"os"
	"testing"

	"github.com/Boomshakalak/VibeRS/internal/store"
)

func TestANNBuildAndRecall(t *testing.T) {
	dbPath := "test_ann.db"
	os.Remove(dbPath)
	db, err := store.InitDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()
	schema := `CREATE TABLE items (
               item_id INTEGER PRIMARY KEY,
               title TEXT,
               brand TEXT,
               price_cents INTEGER,
               discount REAL,
               rating REAL,
               stock INTEGER,
               launched_at DATETIME,
               click_7d INTEGER,
               buy_7d INTEGER,
               gmv_30d INTEGER,
               embedding BLOB
       );`
	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}
	s := store.NewService(db)
	emb := []byte{0, 0, 128, 63, 0, 0, 0, 64, 205, 204, 76, 63} // 1,2,0.95 in float32
	_, err = db.Exec(`INSERT INTO items (item_id, title, brand, price_cents, discount, rating, stock, click_7d, buy_7d, gmv_30d, embedding) VALUES (1,'a','b',100,0,5,1,1,1,100, ?)`, emb)
	if err != nil {
		t.Fatal(err)
	}
	rec := NewANNRecaller(s)
	if err := rec.Build(); err != nil {
		t.Fatal(err)
	}
	items, err := rec.VectorSimilarityRecall([]float32{1, 2, 0.95}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}
