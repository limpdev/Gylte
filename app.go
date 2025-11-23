package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx        context.Context
	db         *sql.DB
	cache      *GlyphCache
	history    *SearchHistory
	favorites  *Favorites
	categories *CategoryManager
}

// Glyph struct for database results
type Glyph struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Glyph    string `json:"glyph"`
	Category string `json:"category,omitempty"`
	Tags     string `json:"tags,omitempty"`
}

// GlyphMatch represents a glyph with its fuzzy match score
type GlyphMatch struct {
	Glyph
	Score      int  `json:"score"`
	IsFavorite bool `json:"isFavorite"`
}

// GlyphCache provides in-memory caching for faster searches
type GlyphCache struct {
	mu     sync.RWMutex
	glyphs []Glyph
	loaded bool
}

// SearchHistory tracks recent searches
type SearchHistory struct {
	mu      sync.RWMutex
	history []string
	maxSize int
}

// Favorites manages user favorites
type Favorites struct {
	mu        sync.RWMutex
	favorites map[int]bool
	db        *sql.DB
}

// CategoryManager handles glyph categorization
type CategoryManager struct {
	mu         sync.RWMutex
	categories map[string][]int
}

// SearchResult wraps results with metadata
type SearchResult struct {
	Glyphs     []GlyphMatch `json:"glyphs"`
	Total      int          `json:"total"`
	SearchTime float64      `json:"searchTime"`
	HasMore    bool         `json:"hasMore"`
	Categories []string     `json:"categories,omitempty"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		cache:      &GlyphCache{},
		history:    &SearchHistory{maxSize: 20},
		favorites:  &Favorites{favorites: make(map[int]bool)},
		categories: &CategoryManager{categories: make(map[string][]int)},
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	var err error
	a.db, err = sql.Open("sqlite", "./gylte.db")
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return
	}

	// Initialize favorites table
	if err := a.initFavoritesTable(); err != nil {
		log.Printf("Failed to initialize favorites: %v", err)
	}

	a.favorites.db = a.db

	// Preload cache in background
	go a.preloadCache()

	// Load favorites
	go a.loadFavorites()

	log.Println("App started successfully")
}

// shutdown cleanup
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// initFavoritesTable creates the favorites table if it doesn't exist
func (a *App) initFavoritesTable() error {
	_, err := a.db.Exec(`
		CREATE TABLE IF NOT EXISTS favorites (
			glyph_id INTEGER PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_favorites_created ON favorites(created_at);
	`)
	return err
}

// preloadCache loads all glyphs into memory
func (a *App) preloadCache() {
	rows, err := a.db.Query("SELECT id, name, glyph FROM glyphs ORDER BY name")
	if err != nil {
		log.Printf("Failed to preload cache: %v", err)
		return
	}
	defer rows.Close()

	a.cache.mu.Lock()
	defer a.cache.mu.Unlock()

	a.cache.glyphs = nil
	for rows.Next() {
		var g Glyph
		if err := rows.Scan(&g.ID, &g.Name, &g.Glyph); err != nil {
			log.Printf("Error scanning glyph: %v", err)
			continue
		}
		a.cache.glyphs = append(a.cache.glyphs, g)

		// Extract category from name (e.g., "nf-cod-account" -> "cod")
		a.categorizeGlyph(&g)
	}

	a.cache.loaded = true
	log.Printf("Cache loaded: %d glyphs", len(a.cache.glyphs))
}

// categorizeGlyph extracts category from glyph name
func (a *App) categorizeGlyph(g *Glyph) {
	parts := strings.Split(g.Name, "-")
	if len(parts) >= 2 {
		category := parts[1] // e.g., "nf-cod-account" -> "cod"

		a.categories.mu.Lock()
		a.categories.categories[category] = append(a.categories.categories[category], g.ID)
		a.categories.mu.Unlock()
	}
}

// loadFavorites loads favorites from database
func (a *App) loadFavorites() {
	rows, err := a.db.Query("SELECT glyph_id FROM favorites")
	if err != nil {
		log.Printf("Failed to load favorites: %v", err)
		return
	}
	defer rows.Close()

	a.favorites.mu.Lock()
	defer a.favorites.mu.Unlock()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			continue
		}
		a.favorites.favorites[id] = true
	}

	log.Printf("Loaded %d favorites", len(a.favorites.favorites))
}

// fuzzyMatch implements fzf-style fuzzy matching
func fuzzyMatch(pattern, text string) (int, bool) {
	pattern = strings.ToLower(pattern)
	text = strings.ToLower(text)

	if pattern == "" {
		return 0, true
	}

	// Exact match gets highest score
	if pattern == text {
		return 10000, true
	}

	// Exact substring match
	if idx := strings.Index(text, pattern); idx != -1 {
		score := 5000
		if idx == 0 {
			score += 2000 // Bonus for prefix match
		}
		score -= len(text) * 2 // Penalty for length
		return score, true
	}

	// Fuzzy matching
	score := 0
	textIdx := 0
	consecutiveMatches := 0
	lastMatchIdx := -1

	for i := 0; i < len(pattern); i++ {
		found := false
		for textIdx < len(text) {
			if pattern[i] == text[textIdx] {
				found = true
				score += 100

				// Bonus for consecutive matches
				if textIdx == lastMatchIdx+1 {
					consecutiveMatches++
					score += consecutiveMatches * 50
				} else {
					consecutiveMatches = 0
				}

				// Bonus for word boundary matches
				if textIdx == 0 || text[textIdx-1] == '-' || text[textIdx-1] == '_' {
					score += 200
				}

				lastMatchIdx = textIdx
				textIdx++
				break
			}
			textIdx++
		}

		if !found {
			return 0, false
		}
	}

	// Penalty for length difference
	score -= (len(text) - len(pattern)) * 3

	return score, true
}

// GetGlyphs retrieves glyphs with advanced filtering
func (a *App) GetGlyphs(searchTerm string, category string, limit int, offset int) (*SearchResult, error) {
	startTime := time.Now()

	// Wait for cache to load if not ready
	for i := 0; i < 50 && !a.cache.loaded; i++ {
		time.Sleep(10 * time.Millisecond)
	}

	a.cache.mu.RLock()
	allGlyphs := a.cache.glyphs
	a.cache.mu.RUnlock()

	if len(allGlyphs) == 0 {
		return &SearchResult{Glyphs: []GlyphMatch{}, Total: 0}, nil
	}

	// Filter by category if specified
	var filtered []Glyph
	if category != "" {
		a.categories.mu.RLock()
		categoryIDs := a.categories.categories[category]
		a.categories.mu.RUnlock()

		idMap := make(map[int]bool)
		for _, id := range categoryIDs {
			idMap[id] = true
		}

		for _, g := range allGlyphs {
			if idMap[g.ID] {
				filtered = append(filtered, g)
			}
		}
	} else {
		filtered = allGlyphs
	}

	// Apply search term
	searchTerm = strings.TrimSpace(searchTerm)
	var matches []GlyphMatch

	if searchTerm == "" {
		// No search term - return all with favorites marked
		a.favorites.mu.RLock()
		for _, g := range filtered {
			matches = append(matches, GlyphMatch{
				Glyph:      g,
				Score:      0,
				IsFavorite: a.favorites.favorites[g.ID],
			})
		}
		a.favorites.mu.RUnlock()
	} else {
		// Apply fuzzy matching
		a.favorites.mu.RLock()
		for _, g := range filtered {
			score, ok := fuzzyMatch(searchTerm, g.Name)
			if ok {
				matches = append(matches, GlyphMatch{
					Glyph:      g,
					Score:      score,
					IsFavorite: a.favorites.favorites[g.ID],
				})
			}
		}
		a.favorites.mu.RUnlock()

		// Sort by favorites first, then by score
		sort.Slice(matches, func(i, j int) bool {
			if matches[i].IsFavorite != matches[j].IsFavorite {
				return matches[i].IsFavorite
			}
			return matches[i].Score > matches[j].Score
		})

		// Add to search history
		if searchTerm != "" {
			a.history.Add(searchTerm)
		}
	}

	// Apply pagination
	total := len(matches)
	if limit <= 0 {
		limit = 50 // Default limit
	}

	start := offset
	end := offset + limit
	if start > len(matches) {
		start = len(matches)
	}
	if end > len(matches) {
		end = len(matches)
	}

	result := &SearchResult{
		Glyphs:     matches[start:end],
		Total:      total,
		SearchTime: time.Since(startTime).Seconds(),
		HasMore:    end < total,
	}

	return result, nil
}

// GetCategories returns all available categories with counts
func (a *App) GetCategories() map[string]int {
	a.categories.mu.RLock()
	defer a.categories.mu.RUnlock()

	result := make(map[string]int)
	for cat, ids := range a.categories.categories {
		result[cat] = len(ids)
	}
	return result
}

// ToggleFavorite adds or removes a glyph from favorites
func (a *App) ToggleFavorite(glyphID int) error {
	a.favorites.mu.Lock()
	defer a.favorites.mu.Unlock()

	if a.favorites.favorites[glyphID] {
		// Remove from favorites
		_, err := a.db.Exec("DELETE FROM favorites WHERE glyph_id = ?", glyphID)
		if err != nil {
			return fmt.Errorf("failed to remove favorite: %w", err)
		}
		delete(a.favorites.favorites, glyphID)
	} else {
		// Add to favorites
		_, err := a.db.Exec("INSERT INTO favorites (glyph_id) VALUES (?)", glyphID)
		if err != nil {
			return fmt.Errorf("failed to add favorite: %w", err)
		}
		a.favorites.favorites[glyphID] = true
	}

	return nil
}

// GetFavorites returns all favorited glyphs
func (a *App) GetFavorites() ([]GlyphMatch, error) {
	a.favorites.mu.RLock()
	favoriteIDs := make([]int, 0, len(a.favorites.favorites))
	for id := range a.favorites.favorites {
		favoriteIDs = append(favoriteIDs, id)
	}
	a.favorites.mu.RUnlock()

	if len(favoriteIDs) == 0 {
		return []GlyphMatch{}, nil
	}

	a.cache.mu.RLock()
	defer a.cache.mu.RUnlock()

	var favorites []GlyphMatch
	idMap := make(map[int]bool)
	for _, id := range favoriteIDs {
		idMap[id] = true
	}

	for _, g := range a.cache.glyphs {
		if idMap[g.ID] {
			favorites = append(favorites, GlyphMatch{
				Glyph:      g,
				IsFavorite: true,
			})
		}
	}

	return favorites, nil
}

// GetSearchHistory returns recent searches
func (a *App) GetSearchHistory() []string {
	a.history.mu.RLock()
	defer a.history.mu.RUnlock()

	result := make([]string, len(a.history.history))
	copy(result, a.history.history)
	return result
}

// ClearSearchHistory clears the search history
func (a *App) ClearSearchHistory() {
	a.history.mu.Lock()
	defer a.history.mu.Unlock()
	a.history.history = nil
}

// CopyToClipboard copies text to clipboard
func (a *App) CopyToClipboard(text string) {
	runtime.ClipboardSetText(a.ctx, text)
}

// GetStats returns app statistics
func (a *App) GetStats() map[string]interface{} {
	a.cache.mu.RLock()
	totalGlyphs := len(a.cache.glyphs)
	a.cache.mu.RUnlock()

	a.favorites.mu.RLock()
	totalFavorites := len(a.favorites.favorites)
	a.favorites.mu.RUnlock()

	a.categories.mu.RLock()
	totalCategories := len(a.categories.categories)
	a.categories.mu.RUnlock()

	return map[string]interface{}{
		"totalGlyphs":     totalGlyphs,
		"totalFavorites":  totalFavorites,
		"totalCategories": totalCategories,
		"cacheLoaded":     a.cache.loaded,
	}
}

// Add method for SearchHistory
func (sh *SearchHistory) Add(term string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	// Remove if already exists
	for i, t := range sh.history {
		if t == term {
			sh.history = append(sh.history[:i], sh.history[i+1:]...)
			break
		}
	}

	// Add to front
	sh.history = append([]string{term}, sh.history...)

	// Trim to max size
	if len(sh.history) > sh.maxSize {
		sh.history = sh.history[:sh.maxSize]
	}
}
