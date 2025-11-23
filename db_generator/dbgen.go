package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("Starting database generation...")

	// Read and parse JSON
	glyphs, err := loadGlyphs("glyphs.json")
	if err != nil {
		return fmt.Errorf("loading glyphs: %w", err)
	}
	log.Printf("Loaded %d glyphs from JSON", len(glyphs))

	// Initialize database
	db, err := initDB("../gylte.db")
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}
	defer db.Close()

	// Populate database
	if err := populateDB(db, glyphs); err != nil {
		return fmt.Errorf("populating database: %w", err)
	}

	// Generate statistics
	stats, err := generateStats(db)
	if err != nil {
		return fmt.Errorf("generating stats: %w", err)
	}

	log.Println("\n=== Database Generation Complete ===")
	log.Printf("Total glyphs: %d", stats["total"])
	log.Printf("Unique categories: %d", stats["categories"])
	log.Printf("Database file: ../gylte.db")

	// Print top categories
	log.Println("\nTop 10 categories:")
	topCats, _ := getTopCategories(db, 10)
	for i, cat := range topCats {
		log.Printf("  %d. %s: %d glyphs", i+1, cat.Name, cat.Count)
	}

	return nil
}

func loadGlyphs(filename string) ([]Glyph, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var glyphs []Glyph
	if err := json.Unmarshal(data, &glyphs); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return glyphs, nil
}

func initDB(filename string) (*sql.DB, error) {
	// Remove old database if exists
	os.Remove(filename)

	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	// Create schema with enhanced metadata
	schema := `
	PRAGMA journal_mode = WAL;
	PRAGMA synchronous = NORMAL;
	PRAGMA cache_size = 10000;
	PRAGMA temp_store = MEMORY;

	CREATE TABLE IF NOT EXISTS glyphs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		glyph TEXT NOT NULL,
		category TEXT,
		prefix TEXT,
		normalized_name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Indexes for fast searching
	CREATE INDEX IF NOT EXISTS idx_name ON glyphs(name);
	CREATE INDEX IF NOT EXISTS idx_category ON glyphs(category);
	CREATE INDEX IF NOT EXISTS idx_prefix ON glyphs(prefix);
	CREATE INDEX IF NOT EXISTS idx_normalized ON glyphs(normalized_name);

	-- Full-text search support
	CREATE VIRTUAL TABLE IF NOT EXISTS glyphs_fts USING fts5(
		name, 
		category,
		content='glyphs',
		content_rowid='id'
	);

	-- Triggers to keep FTS in sync
	CREATE TRIGGER IF NOT EXISTS glyphs_ai AFTER INSERT ON glyphs BEGIN
		INSERT INTO glyphs_fts(rowid, name, category)
		VALUES (new.id, new.name, new.category);
	END;

	CREATE TRIGGER IF NOT EXISTS glyphs_ad AFTER DELETE ON glyphs BEGIN
		DELETE FROM glyphs_fts WHERE rowid = old.id;
	END;

	CREATE TRIGGER IF NOT EXISTS glyphs_au AFTER UPDATE ON glyphs BEGIN
		UPDATE glyphs_fts SET name = new.name, category = new.category
		WHERE rowid = new.id;
	END;

	-- Metadata table for app info
	CREATE TABLE IF NOT EXISTS metadata (
		key TEXT PRIMARY KEY,
		value TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("creating schema: %w", err)
	}

	return db, nil
}

func extractMetadata(name string) (category, prefix, normalized string) {
	// Split on hyphens: "nf-cod-account" -> ["nf", "cod", "account"]
	parts := strings.Split(name, "-")

	if len(parts) >= 1 {
		prefix = parts[0] // "nf"
	}

	if len(parts) >= 2 {
		category = parts[1] // "cod"
	}

	// Normalized name without prefix (for better searching)
	if len(parts) > 1 {
		normalized = strings.Join(parts[1:], " ")
	} else {
		normalized = name
	}

	return
}

func populateDB(db *sql.DB, glyphs []Glyph) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO glyphs(name, glyph, category, prefix, normalized_name) 
		VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	duplicates := 0
	for i, glyph := range glyphs {
		category, prefix, normalized := extractMetadata(glyph.Name)

		_, err := stmt.Exec(
			glyph.Name,
			glyph.Glyph,
			category,
			prefix,
			normalized,
		)

		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				duplicates++
				continue
			}
			return fmt.Errorf("inserting glyph %s: %w", glyph.Name, err)
		}

		// Progress indicator
		if (i+1)%1000 == 0 {
			log.Printf("Processed %d/%d glyphs...", i+1, len(glyphs))
		}
	}

	if duplicates > 0 {
		log.Printf("Skipped %d duplicate glyphs", duplicates)
	}

	// Store metadata
	_, err = tx.Exec(`
		INSERT OR REPLACE INTO metadata(key, value) 
		VALUES('last_updated', datetime('now'))
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT OR REPLACE INTO metadata(key, value) 
		VALUES('version', '1.0')
	`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func generateStats(db *sql.DB) (map[string]int, error) {
	stats := make(map[string]int)

	var total, categories, prefixes int

	// Total glyphs
	if err := db.QueryRow("SELECT COUNT(*) FROM glyphs").Scan(&total); err != nil {
		return nil, err
	}

	// Unique categories
	if err := db.QueryRow("SELECT COUNT(DISTINCT category) FROM glyphs WHERE category != ''").Scan(&categories); err != nil {
		return nil, err
	}

	// Unique prefixes
	if err := db.QueryRow("SELECT COUNT(DISTINCT prefix) FROM glyphs WHERE prefix != ''").Scan(&prefixes); err != nil {
		return nil, err
	}

	stats["total"] = total
	stats["categories"] = categories
	stats["prefixes"] = prefixes

	return stats, nil
}

type CategoryStats struct {
	Name  string
	Count int
}

func getTopCategories(db *sql.DB, limit int) ([]CategoryStats, error) {
	rows, err := db.Query(`
		SELECT category, COUNT(*) as count 
		FROM glyphs 
		WHERE category != '' 
		GROUP BY category 
		ORDER BY count DESC 
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryStats
	for rows.Next() {
		var cat CategoryStats
		if err := rows.Scan(&cat.Name, &cat.Count); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, rows.Err()
}
