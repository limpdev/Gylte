package main

import (
	"context"
	"database/sql"
	"log"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
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

// GlyphMatch represents a glyph with its fuzzy match score
type GlyphMatch struct {
	Glyph
	Score int
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
	a.db, err = sql.Open("sqlite", "./gylte.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
}

// fuzzyMatch implements fzf-style fuzzy matching
// Returns a score (higher is better) and whether it matches
func fuzzyMatch(pattern, text string) (int, bool) {
	pattern = strings.ToLower(pattern)
	text = strings.ToLower(text)

	if pattern == "" {
		return 0, true
	}

	// Exact substring match gets highest priority
	if strings.Contains(text, pattern) {
		score := 1000
		// Bonus for matches at the beginning
		if strings.HasPrefix(text, pattern) {
			score += 500
		}
		// Bonus for shorter strings (more relevant)
		score -= len(text) - len(pattern)
		return score, true
	}

	// Fuzzy matching: all characters of pattern must appear in order in text
	patternChars := []rune(pattern)
	textChars := []rune(text)

	patternIdx := 0
	score := 0
	consecutiveMatches := 0

	for textIdx, char := range textChars {
		if patternIdx < len(patternChars) && char == patternChars[patternIdx] {
			// Character matches
			patternIdx++
			consecutiveMatches++

			// Bonus for consecutive matches
			score += consecutiveMatches * 2

			// Bonus for matches at word boundaries
			if textIdx == 0 || textChars[textIdx-1] == ' ' || textChars[textIdx-1] == '-' || textChars[textIdx-1] == '_' {
				score += 10
			}
		} else {
			consecutiveMatches = 0
		}
	}

	// All pattern characters must be matched
	if patternIdx == len(patternChars) {
		// Penalty for longer strings
		score -= len(text) - len(pattern)
		return score, true
	}

	return 0, false
}

// GetGlyphs retrieves all glyphs or filters them by a search term with fuzzy matching
func (a *App) GetGlyphs(searchTerm string) ([]Glyph, error) {
	// Always get all glyphs from database
	rows, err := a.db.Query("SELECT name, glyph FROM glyphs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allGlyphs []Glyph
	for rows.Next() {
		var glyph Glyph
		if err := rows.Scan(&glyph.Name, &glyph.Glyph); err != nil {
			return nil, err
		}
		allGlyphs = append(allGlyphs, glyph)
	}

	// If no search term, return all glyphs
	if strings.TrimSpace(searchTerm) == "" {
		return allGlyphs, nil
	}

	// Apply fuzzy matching
	var matches []GlyphMatch
	for _, glyph := range allGlyphs {
		if score, isMatch := fuzzyMatch(searchTerm, glyph.Name); isMatch {
			matches = append(matches, GlyphMatch{
				Glyph: glyph,
				Score: score,
			})
		}
	}

	// Sort by score (highest first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	// Convert back to []Glyph
	var result []Glyph
	for _, match := range matches {
		result = append(result, match.Glyph)
	}

	return result, nil
}

// CopyToClipboard takes a string and copies it to the user's clipboard.
func (a *App) CopyToClipboard(text string) {
	runtime.ClipboardSetText(a.ctx, text)
}
