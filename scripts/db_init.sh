#!/bin/bash

# VibeRS Database Initialization Script
# Creates SQLite database and imports sample data

set -e

DB_PATH="./data/vibers.db"
DDL_PATH="./data/ddl.sql"
SAMPLE_PATH="./data/sample.csv"

echo "🚀 Initializing VibeRS database..."

# Remove existing database if it exists
if [ -f "$DB_PATH" ]; then
    echo "📝 Removing existing database..."
    rm "$DB_PATH"
fi

# Create database and run DDL
echo "📋 Creating database schema..."
sqlite3 "$DB_PATH" < "$DDL_PATH"

# Import sample data using INSERT statements
echo "📊 Importing sample data..."
cat > /tmp/import_data.sql << 'EOF'
INSERT INTO items (item_id, title, brand, price_cents, discount, rating, stock, click_7d, buy_7d, gmv_30d) VALUES
(1, 'Louis Vuitton Neverfull MM Tote Bag', 'Louis Vuitton', 150000, 0.0, 4.8, 5, 245, 12, 1800000),
(2, 'Gucci GG Marmont Small Shoulder Bag', 'Gucci', 180000, 0.1, 4.7, 3, 189, 8, 1440000),
(3, 'Chanel Classic Flap Bag Medium', 'Chanel', 650000, 0.0, 4.9, 2, 456, 23, 14950000),
(4, 'Hermès Birkin 30cm Togo Leather', 'Hermès', 1200000, 0.0, 4.9, 1, 89, 5, 6000000),
(5, 'Prada Re-Edition 2005 Nylon Bag', 'Prada', 120000, 0.15, 4.6, 8, 134, 15, 1800000),
(6, 'Saint Laurent Loulou Small Bag', 'Saint Laurent', 165000, 0.05, 4.5, 6, 98, 7, 1155000),
(7, 'Bottega Veneta Jodie Small Bag', 'Bottega Veneta', 280000, 0.0, 4.8, 4, 167, 11, 3080000),
(8, 'Fendi Baguette Bag', 'Fendi', 320000, 0.0, 4.7, 3, 234, 18, 5760000),
(9, 'Dior Saddle Bag', 'Dior', 380000, 0.0, 4.8, 2, 345, 21, 7980000),
(10, 'Balenciaga City Bag Classic', 'Balenciaga', 195000, 0.12, 4.4, 7, 76, 6, 1170000),
(11, 'Celine Luggage Nano Tote', 'Celine', 290000, 0.0, 4.6, 5, 123, 9, 2610000),
(12, 'Givenchy Antigona Small Bag', 'Givenchy', 225000, 0.08, 4.5, 6, 87, 8, 1800000),
(13, 'Valentino Rockstud Spike Bag', 'Valentino', 175000, 0.1, 4.3, 9, 65, 5, 875000),
(14, 'Loewe Puzzle Bag Small', 'Loewe', 245000, 0.0, 4.7, 4, 102, 12, 2940000),
(15, 'Jacquemus Le Chiquito Mini Bag', 'Jacquemus', 65000, 0.2, 4.2, 15, 234, 45, 2925000),
(16, 'Staud Shirley Bag', 'Staud', 25000, 0.25, 4.1, 22, 189, 38, 950000),
(17, 'Mansur Gavriel Bucket Bag', 'Mansur Gavriel', 45000, 0.15, 4.3, 18, 156, 28, 1260000),
(18, 'Cult Gaia Ark Bag Small', 'Cult Gaia', 35000, 0.1, 4.0, 25, 298, 67, 2345000),
(19, 'Polene Numero Un Mini', 'Polene', 38000, 0.0, 4.4, 20, 145, 32, 1216000),
(20, 'Wandler Hortensia Mini Bag', 'Wandler', 42000, 0.05, 4.2, 16, 87, 19, 798000);
EOF

sqlite3 "$DB_PATH" < /tmp/import_data.sql
rm /tmp/import_data.sql

# Verify import
echo "✅ Verifying data import..."
COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM items;")
echo "   Imported $COUNT items"

echo "🎉 Database initialization complete!"
echo "   Database: $DB_PATH"
echo "   Items: $COUNT"
echo ""
echo "Quick test:"
echo "  sqlite3 $DB_PATH \"SELECT title, brand, price_cents FROM items LIMIT 3;\""
echo "  sqlite3 $DB_PATH \"SELECT title FROM items WHERE title LIKE '%bag%' LIMIT 3;\"" 