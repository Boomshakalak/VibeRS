#!/usr/bin/env python3
"""Generate a synthetic SQLite dataset with random items and embeddings.

Usage:
  python gen_mock_data.py --items 100000 --db data/vibers_big.db
"""
import argparse
import random
import sqlite3
import struct
from pathlib import Path

BRANDS = [
    "Gucci", "Louis Vuitton", "Chanel", "Hermes", "Prada",
    "Saint Laurent", "Bottega Veneta", "Fendi", "Dior", "Balenciaga",
]

ADJECTIVES = [
    "Classic", "Vintage", "Elegant", "Modern", "Chic",
    "Luxury", "Sport", "Urban", "Retro", "Bold",
]

ITEMS = ["Bag", "Wallet", "Belt", "Shoe", "Backpack"]


def random_embedding(dim: int) -> bytes:
    values = [random.uniform(-1, 1) for _ in range(dim)]
    return struct.pack('<%df' % dim, *values)


def generate(db_path: Path, num_items: int, dim: int = 8) -> None:
    if db_path.exists():
        db_path.unlink()
    conn = sqlite3.connect(str(db_path))
    ddl = Path("data/ddl.sql").read_text()
    conn.executescript(ddl)
    for i in range(1, num_items + 1):
        brand = random.choice(BRANDS)
        title = f"{brand} {random.choice(ADJECTIVES)} {random.choice(ITEMS)} {i}"
        price = random.randint(20000, 500000)
        discount = round(random.uniform(0, 0.3), 2)
        rating = round(random.uniform(3.5, 5.0), 1)
        stock = random.randint(1, 50)
        click = random.randint(0, 500)
        buy = random.randint(0, 50)
        gmv = price * buy
        embedding = random_embedding(dim)
        conn.execute(
            "INSERT INTO items (item_id, title, brand, price_cents, discount, rating, stock, click_7d, buy_7d, gmv_30d, embedding) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
            (i, title, brand, price, discount, rating, stock, click, buy, gmv, embedding),
        )
    conn.commit()
    conn.close()


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate synthetic dataset")
    parser.add_argument("--items", type=int, default=100000, help="number of items")
    parser.add_argument("--db", type=Path, default=Path("data/vibers_big.db"), help="output database path")
    parser.add_argument("--dim", type=int, default=8, help="embedding dimension")
    args = parser.parse_args()
    args.db.parent.mkdir(parents=True, exist_ok=True)
    generate(args.db, args.items, args.dim)
    print(f"Generated {args.items} items in {args.db}")


if __name__ == "__main__":
    main()
