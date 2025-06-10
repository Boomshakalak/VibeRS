# Parallel Recall → Dedup → Three‑Stage Ranking Prototype

### Go (back‑end) · SQLite (storage) · Python (model training)

A **laptop‑friendly playground** to test an e‑commerce search/recommend stack that mixes
parallel recall, in‑memory de‑duplication and three‑stage ranking.
Everything runs on **SQLite + pure Go**, with Python only needed for offline model training.

---

## 0 · Quick Start (no Makefile required)

```bash
# 0. Prereqs (macOS / Linux)
brew install go sqlite3               # or apt‑get install sqlite3
python ‑m pip install ‑U pip virtualenv

# 1. Clone & boot DB
git clone git@github.com:Boomshakalak/VibeRS.git
cd VibeRS
scripts/db_init.sh                    # executes data/ddl.sql + import sample.csv

# 2. Launch API (hot‑reload)
go run ./cmd/api                      # listens on :8080

# 3. Query
curl -X POST localhost:8080/search -d '{"q":"lv bag","page":1}' | jq
```

> **tip:** prefer `go test ./...` & `go vet ./...` for code checks – no Make needed.

---

## 1 · Repository Layout

```text
.
├── cmd/                # Go entrypoints (api, batch, cli)
├── internal/           # domain logic – each pkg ≈ 200 LoC max
│   ├── recall/         # text.go, attr.go, ann.go, hot.go, exp.go
│   ├── dedup/          # min‑heap & bloom filters
│   ├── rank/
│   │   ├── coarse/     # rule engine (pure Go template)
│   │   ├── ltr/        # ONNX runtime wrapper
│   │   └── final/      # greedy / LP re‑rank
│   ├── store/          # SQLite DAO + UDF (cosine)
│   └── util/
├── data/               # ddl.sql + sample.csv (10 K rows)
├── model‑training/     # Python notebooks + tools
│   └── requirements.txt
├── scripts/            # db_init.sh, lint.sh, bench.sh
└── README.md
```

> **No brain‑split:** each sub‑folder can be developed & unit‑tested in isolation.

---

## 2 · SQLite Schema (MVP)

```sql
CREATE TABLE items (
  item_id       INTEGER PRIMARY KEY,
  title         TEXT,
  brand         TEXT,
  price_cents   INTEGER,
  discount      REAL,   -- 0‑1
  rating        REAL,
  stock         INTEGER,
  launched_at   DATETIME,
  click_7d      INTEGER,
  buy_7d        INTEGER,
  gmv_30d       INTEGER,
  embedding     BLOB     -- []float32  serialized
);
CREATE INDEX idx_items_brand_price ON items(brand, price_cents);
-- Full‑text index for fuzzy recall
CREATE VIRTUAL TABLE items_fts USING fts5(title, brand, content='items', content_rowid='item_id');
-- Optional spellfix1 for typo‑tolerant search suggestions
-- CREATE VIRTUAL TABLE spellfix1 USING spellfix1;
-- INSERT INTO spellfix1(word) SELECT DISTINCT title FROM items;
```

The Go layer registers a **Cosine(embedding, queryVec)** UDF so that vector recall can be done directly in SQL.

---

## 3 · Parallel Recall Layer

| File    | Strategy       | SQL / Logic example                                                                  | Batch size |
| ------- | -------------- | ------------------------------------------------------------------------------------ | ---------- |
| text.go | Text & fuzzy   | `title MATCH ?` via **FTS5**  (+ optional `spellfix1` for typo‑tolerant suggestions) | 1 K        |
| attr.go | Filter rules   | `brand=? AND price_cents<?`                                                          | 1‑2 K      |
| ann.go  | ANN similarity | `ORDER BY Cosine(embedding,?) DESC`                                                  | 1 K        |
| hot.go  | Hot‑pool       | yesterday GMV Top‑1 K (pre‑query)                                                    | ≤1 K       |
| exp.go  | Exploration    | `ORDER BY RANDOM() LIMIT 500`                                                        | 0.5 K      |

Each returns `(items, nextCursor)`; cursors are local JSON tokens `{src, lastID, score}`.

---

## 4 · In‑memory Dedup + Merge

1. insert unseen items into a **min‑heap** keyed by coarse‑rank score
2. keep `seen` map to deduplicate by `item_id`
3. pop top‑N for requested page.

CPU cost: **O(N log R)** where *N ≈ total recall* (≤3 K) – few ms.

---

## 5 · Three‑Stage Ranking

| Stage  | Goal                             | Implementation (Go)          | Latency |
| ------ | -------------------------------- | ---------------------------- | ------- |
| Coarse | Hard rules (stock, price band …) | `/rank/coarse/rules.go`      | <5 ms   |
| LTR    | Buy‑probability score            | ONNX runtime + XGBoost model | \~10 ms |
| Final  | GMV × New × Brand fairness       | `/rank/final/greedy.go`      | \~10 ms |

### Offline training pipeline

```bash
cd model‑training
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
python extract.py       # parquet features → train
jupyter notebook train.ipynb  # exports model.onnx
python eval.py          # prints AUC, NDCG@10
```

`model‑training/requirements.txt` minimal example:

```
onnxruntime>=1.18.0
xgboost>=2.0.0
pandas>=2.2
scikit‑learn>=1.5
```

---

## 6 · Pagination Strategy

* **Snapshot list (default):** first page builds the full ordered list → store in `sessionCache` (sync.Map / Redis), return `cursor=sessionID|offset`.
* **K‑Way merge (opt‑in):** for ultra‑hot queries switch to stateful merge cursors.

TTL default **10 min**; after expiry client gets a fresh snapshot.

---

## 7 · Common Dev Commands

| What                    | Command                                                |
| ----------------------- | ------------------------------------------------------ |
| Run API (dev)           | `go run ./cmd/api`                                     |
| Unit tests              | `go test ./...`                                        |
| Lint                    | `go vet ./...`                                         |
| Initialise DB           | `scripts/db_init.sh`                                   |
| Python model env        | `cd model‑training && pip install -r requirements.txt` |
| Benchmark 1 K QPS (WIP) | `scripts/bench.sh`                                     |

*You may still use the included **Makefile** as a thin wrapper, but it's optional.*

---

## 8 · Development Progress

### ✅ Completed
- [x] **Project Structure**: Complete directory layout with all core packages
- [x] **Database Schema**: SQLite schema with FTS5, indexes, and triggers  
- [x] **Sample Data**: 20 luxury handbag items for testing
- [x] **API Framework**: Gin-based REST API with search endpoint
- [x] **Store Layer**: SQLite operations with cosine similarity UDF placeholder
- [x] **Parallel Recall**: Multi-strategy recall service (text, attr, hot, explore, ANN placeholder)
- [x] **Deduplication**: Min-heap based dedup with basic scoring
- [x] **Three-Stage Ranking**: Coarse (rules) → LTR (placeholder) → Final (greedy diversity)
- [x] **Scripts**: Database initialization and code linting automation
- [x] **Build System**: Both Makefile and direct Go commands support

### 🚧 In Progress  
- [ ] **Vector Similarity**: Implement actual cosine similarity UDF in SQLite
- [ ] **ONNX Integration**: Load and run machine learning models for LTR
- [ ] **Feature Engineering**: Complete ML feature extraction pipeline
- [ ] **Unit Tests**: Comprehensive test coverage for all packages

### 📋 Next Sprint
- [ ] **Database Initialization**: Run `./scripts/db_init.sh` and test API
- [ ] **Model Training**: Python notebooks for XGBoost → ONNX export  
- [ ] **Performance Testing**: Benchmark with concurrent requests
- [ ] **Pagination**: Implement session-based result caching

---

## 9 · Future Roadmap

* Swap SQLite ANN with **Qdrant (HNSW)** once recall logic is solid
* Add **Typesense / ES** for large‑scale text & filter recall
* **Redis** snapshot cache + hot‑key limiter
* Benchmark target **P99 ≤ 100 ms @ 1 K QPS** on single VM
* YJS + React front‑end demo (GitHub Pages)

---

## 10 · Contributing (vibe‑style)

1. Create an Issue → *describe intent in natural language*.
2. Let the AI agent scaffold the PR ("vibe coding").
3. Add tests & README snippet that explains what changed.

> **Less control, more vibing** – keep modules small, readable, and AI‑friendly.

---

Happy hacking! ✨
