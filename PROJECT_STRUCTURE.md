# VibeRS Project Structure

## 📁 Directory Layout

```
VibeRS/
├── cmd/                    # Application entry points
│   ├── api/               # REST API server
│   │   └── main.go        # API server main entry
│   ├── batch/             # Batch processing jobs (placeholder)
│   └── cli/               # Command line tools (placeholder)
│
├── internal/              # Private application packages
│   ├── store/             # Data access layer
│   │   └── store.go       # SQLite operations & Item model
│   ├── recall/            # Parallel recall strategies
│   │   └── service.go     # Multi-strategy recall service
│   ├── dedup/             # Deduplication logic
│   │   └── service.go     # Min-heap based deduplication
│   ├── rank/              # Three-stage ranking
│   │   ├── coarse/        # Hard rules & basic scoring
│   │   │   └── ranker.go  # Business rule engine
│   │   ├── ltr/           # Learning-to-Rank
│   │   │   └── ranker.go  # ML-based ranking (ONNX placeholder)
│   │   └── final/         # Final business optimization
│   │       └── ranker.go  # Diversity & GMV optimization
│   └── util/              # Shared utilities (placeholder)
│
├── data/                  # Database & sample data
│   ├── ddl.sql           # Database schema definition
│   ├── sample.csv        # 20 luxury handbag sample items
│   └── vibers.db         # SQLite database (created by init script)
│
├── model-training/        # Python ML pipeline
│   └── requirements.txt  # Python dependencies for model training
│
├── scripts/              # Automation scripts
│   ├── db_init.sh        # Database initialization script
│   └── lint.sh           # Code quality checks
│
├── go.mod                # Go module definition
├── Makefile              # Optional build automation
└── README.md             # Main project documentation
```

## 🔧 Key Components

### API Layer (`cmd/api/`)
- **main.go**: Gin-based REST API server
- Handles `/search` endpoint with JSON request/response
- Orchestrates the full search pipeline

### Data Layer (`internal/store/`)
- **store.go**: SQLite database operations
- Item model with all product attributes
- Full-text search, filtering, and vector similarity (placeholder)
- Custom SQLite UDF for cosine similarity

### Recall Layer (`internal/recall/`)
- **service.go**: Parallel recall execution
- 5 concurrent strategies: text, attribute, hot, explore, ANN
- Goroutine-based parallel processing

### Deduplication (`internal/dedup/`)
- **service.go**: Min-heap based deduplication
- Removes duplicate items by ID
- Maintains top-N items by relevance score

### Ranking Pipeline (`internal/rank/`)
- **coarse/ranker.go**: Hard business rules & basic scoring
- **ltr/ranker.go**: Machine learning based ranking (placeholder)
- **final/ranker.go**: Business optimization with diversity constraints

### Database (`data/`)
- **ddl.sql**: Complete SQLite schema with FTS5 and indexes
- **sample.csv**: 20 sample luxury handbag items
- Supports full-text search, filtering, and vector operations

### Scripts (`scripts/`)
- **db_init.sh**: Automated database setup and data import
- **lint.sh**: Code quality checks (go vet, fmt, mod tidy)

## 🚀 Quick Start

1. **Initialize Database**:
   ```bash
   ./scripts/db_init.sh
   ```

2. **Run API Server**:
   ```bash
   go run ./cmd/api
   ```

3. **Test Search**:
   ```bash
   curl -X POST localhost:8080/search -d '{"q":"lv bag","page":1}' | jq
   ```

## 📊 Current Status

- ✅ **Complete project structure** with all packages
- ✅ **Working API server** with search endpoint
- ✅ **Database schema** with sample data
- ✅ **Parallel recall** framework
- ✅ **Three-stage ranking** pipeline
- 🚧 **Vector similarity** (placeholder implementation)
- 🚧 **ONNX model integration** (placeholder)
- 📋 **Unit tests** (to be added)

## 🎯 Next Steps

1. Run database initialization
2. Test API functionality
3. Implement vector similarity UDF
4. Add ONNX model integration
5. Create comprehensive unit tests
6. Performance benchmarking 