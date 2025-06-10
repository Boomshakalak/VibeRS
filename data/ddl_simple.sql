-- VibeRS Database Schema (Simplified)
-- SQLite-based e-commerce search system

CREATE TABLE IF NOT EXISTS items (
  item_id       INTEGER PRIMARY KEY,
  title         TEXT NOT NULL,
  brand         TEXT,
  price_cents   INTEGER,
  discount      REAL DEFAULT 0,   -- 0-1
  rating        REAL DEFAULT 0,
  stock         INTEGER DEFAULT 0,
  launched_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
  click_7d      INTEGER DEFAULT 0,
  buy_7d        INTEGER DEFAULT 0,
  gmv_30d       INTEGER DEFAULT 0,
  embedding     BLOB              -- []float32 serialized
);

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_items_brand_price ON items(brand, price_cents);
CREATE INDEX IF NOT EXISTS idx_items_rating ON items(rating DESC);
CREATE INDEX IF NOT EXISTS idx_items_stock ON items(stock);
CREATE INDEX IF NOT EXISTS idx_items_launched ON items(launched_at DESC);
CREATE INDEX IF NOT EXISTS idx_items_gmv ON items(gmv_30d DESC);
CREATE INDEX IF NOT EXISTS idx_items_title ON items(title);
CREATE INDEX IF NOT EXISTS idx_items_brand ON items(brand);

-- Session cache for pagination
CREATE TABLE IF NOT EXISTS session_cache (
  session_id    TEXT PRIMARY KEY,
  query_hash    TEXT NOT NULL,
  results       BLOB NOT NULL,  -- JSON serialized results
  created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
  expires_at    DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_session_expires ON session_cache(expires_at);

-- User behavior tracking (optional)
CREATE TABLE IF NOT EXISTS user_actions (
  action_id     INTEGER PRIMARY KEY,
  user_id       TEXT,
  item_id       INTEGER,
  action_type   TEXT,  -- 'view', 'click', 'add_to_cart', 'buy'
  query         TEXT,
  timestamp     DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (item_id) REFERENCES items(item_id)
);

CREATE INDEX IF NOT EXISTS idx_user_actions_item ON user_actions(item_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_user ON user_actions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_time ON user_actions(timestamp DESC); 