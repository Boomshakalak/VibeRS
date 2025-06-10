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

# Import sample data
echo "📊 Importing sample data..."
sqlite3 "$DB_PATH" <<EOF
.mode csv
.headers on
.import $SAMPLE_PATH items
EOF

# Verify import
echo "✅ Verifying data import..."
COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM items;")
echo "   Imported $COUNT items"

# Build FTS index
echo "🔍 Building full-text search index..."
sqlite3 "$DB_PATH" "INSERT INTO items_fts(items_fts) VALUES('rebuild');"

echo "🎉 Database initialization complete!"
echo "   Database: $DB_PATH"
echo "   Items: $COUNT"
echo ""
echo "Quick test:"
echo "  sqlite3 $DB_PATH \"SELECT title, brand, price_cents FROM items LIMIT 3;\"" 