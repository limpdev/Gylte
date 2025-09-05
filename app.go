package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// Glyph struct for database results
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// and the database connection is initialized.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	var err error
	// Use "sqlite" instead of "sqlite3" for modernc.org/sqlite
	a.db, err = sql.Open("sqlite", "./gylte.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
}

// GetGlyphs retrieves all glyphs or filters them by a search term.
func (a *App) GetGlyphs(searchTerm string) ([]Glyph, error) {
	var rows *sql.Rows
	var err error

	if searchTerm == "" {
		rows, err = a.db.Query("SELECT name, glyph FROM glyphs")
	} else {
		rows, err = a.db.Query("SELECT name, glyph FROM glyphs WHERE name LIKE ?", "%"+searchTerm+"%")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var glyphs []Glyph
	for rows.Next() {
		var glyph Glyph
		if err := rows.Scan(&glyph.Name, &glyph.Glyph); err != nil {
			return nil, err
		}
		glyphs = append(glyphs, glyph)
	}

	return glyphs, nil
}

// CopyToClipboard takes a string and copies it to the user's clipboard.
func (a *App) CopyToClipboard(text string) {
	runtime.ClipboardSetText(a.ctx, text)
}
