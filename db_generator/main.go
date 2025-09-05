package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// Glyph struct to match the JSON structure
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}

func main() {
	// Open and read the JSON file
	jsonFile, err := os.Open("../frontend/src/glyphs.json")
	if err != nil {
		log.Fatal("Error opening JSON file:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var glyphs []Glyph
	json.Unmarshal(byteValue, &glyphs)

	// Create a new SQLite database file
	// Use "sqlite" instead of "sqlite3" for modernc.org/sqlite
	db, err := sql.Open("sqlite", "./gylte.db")
	if err != nil {
		log.Fatal("Error creating database:", err)
	}
	defer db.Close()

	// Create the glyphs table
	sqlStmt := `
	CREATE TABLE glyphs (id INTEGER NOT NULL PRIMARY KEY, name TEXT, glyph TEXT);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// Insert unique glyphs into the database
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO glyphs(name, glyph) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	uniqueGlyphs := make(map[string]bool)
	for _, glyph := range glyphs {
		if _, ok := uniqueGlyphs[glyph.Name]; !ok {
			_, err = stmt.Exec(glyph.Name, glyph.Glyph)
			if err != nil {
				log.Fatal(err)
			}
			uniqueGlyphs[glyph.Name] = true
		}
	}
	tx.Commit()

	log.Println("Database 'gylte.db' created and populated successfully.")
}
