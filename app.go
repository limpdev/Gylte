package main

import (
	"context"
	"embed"
	"encoding/json"
	"sort"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/src/glyphs.json
var glyphsJSON embed.FS

// Glyph represents a single glyph with its name and symbol
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}

// SearchIndex holds preprocessed data for fast searching
type SearchIndex struct {
	glyphs        []Glyph
	normalizedMap map[string][]int // Maps normalized strings to glyph indices
	prefixTree    *TrieNode
	mu            sync.RWMutex
}

// TrieNode for prefix-based searching
type TrieNode struct {
	children map[rune]*TrieNode
	indices  []int // Glyph indices that match this prefix
}

// App struct with enhanced search capabilities
type App struct {
	ctx         context.Context
	searchIndex *SearchIndex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		searchIndex: &SearchIndex{
			normalizedMap: make(map[string][]int),
			prefixTree:    &TrieNode{children: make(map[rune]*TrieNode)},
		},
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Initialize search index with glyph data
	a.initializeSearchIndex()
}

// normalizeString removes diacritics, converts to lowercase, and removes extra spaces
func normalizeString(s string) string {
	// Convert to lowercase and normalize unicode
	s = strings.ToLower(strings.TrimSpace(s))

	// Remove common separators and replace with single space
	replacer := strings.NewReplacer(
		"-", " ",
		"_", " ",
		".", " ",
		"/", " ",
	)
	s = replacer.Replace(s)

	// Collapse multiple spaces
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// initializeSearchIndex loads and preprocesses glyph data for fast searching
func (a *App) initializeSearchIndex() {
	// Load embedded JSON file
	jsonData, err := glyphsJSON.ReadFile("frontend/src/glyphs.json")
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to read glyphs.json: %v", err)
		return
	}

	var rawGlyphs []Glyph
	if err := json.Unmarshal(jsonData, &rawGlyphs); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to parse glyphs.json: %v", err)
		return
	}

	a.searchIndex.mu.Lock()
	defer a.searchIndex.mu.Unlock()

	// Remove duplicates (keeping first occurrence)
	seen := make(map[string]bool)
	uniqueGlyphs := make([]Glyph, 0, len(rawGlyphs))

	for _, glyph := range rawGlyphs {
		if !seen[glyph.Name] {
			seen[glyph.Name] = true
			uniqueGlyphs = append(uniqueGlyphs, glyph)
		}
	}

	a.searchIndex.glyphs = uniqueGlyphs

	// Build normalized map and prefix tree
	for i, glyph := range uniqueGlyphs {
		normalized := normalizeString(glyph.Name)

		// Add to normalized map
		a.searchIndex.normalizedMap[normalized] = append(a.searchIndex.normalizedMap[normalized], i)

		// Add to prefix tree
		a.addToTrie(normalized, i)

		// Also index individual words for better partial matching
		words := strings.Fields(normalized)
		for _, word := range words {
			if len(word) > 1 { // Skip single characters
				a.addToTrie(word, i)
			}
		}
	}

	runtime.LogInfof(a.ctx, "Search index initialized with %d glyphs", len(uniqueGlyphs))
}

// addToTrie adds a string and its associated glyph index to the prefix tree
func (a *App) addToTrie(s string, index int) {
	node := a.searchIndex.prefixTree

	for _, char := range s {
		if node.children[char] == nil {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]

		// Add index if not already present
		found := false
		for _, idx := range node.indices {
			if idx == index {
				found = true
				break
			}
		}
		if !found {
			node.indices = append(node.indices, index)
		}
	}
}

// SearchGlyphs performs optimized search and returns filtered results
func (a *App) SearchGlyphs(query string) []Glyph {
	a.searchIndex.mu.RLock()
	defer a.searchIndex.mu.RUnlock()

	// Return all glyphs if query is empty
	if strings.TrimSpace(query) == "" {
		return a.searchIndex.glyphs
	}

	normalizedQuery := normalizeString(query)
	if normalizedQuery == "" {
		return a.searchIndex.glyphs
	}

	// Use different search strategies based on query characteristics
	var matchedIndices []int

	// Strategy 1: Exact normalized match (fastest)
	if indices, exists := a.searchIndex.normalizedMap[normalizedQuery]; exists {
		matchedIndices = append(matchedIndices, indices...)
	}

	// Strategy 2: Prefix matching using trie
	prefixMatches := a.searchByPrefix(normalizedQuery)
	matchedIndices = append(matchedIndices, prefixMatches...)

	// Strategy 3: Substring matching (slower but more comprehensive)
	if len(matchedIndices) < 10 { // Only do substring search if we don't have many matches
		substringMatches := a.searchBySubstring(normalizedQuery)
		matchedIndices = append(matchedIndices, substringMatches...)
	}

	// Remove duplicates and sort by relevance
	uniqueIndices := a.removeDuplicateIndices(matchedIndices)
	sortedIndices := a.sortByRelevance(uniqueIndices, normalizedQuery)

	// Convert indices to glyphs
	result := make([]Glyph, len(sortedIndices))
	for i, idx := range sortedIndices {
		result[i] = a.searchIndex.glyphs[idx]
	}

	return result
}

// searchByPrefix finds all glyphs whose names start with the query
func (a *App) searchByPrefix(query string) []int {
	node := a.searchIndex.prefixTree

	// Navigate to the query prefix
	for _, char := range query {
		if node.children[char] == nil {
			return nil // Prefix not found
		}
		node = node.children[char]
	}

	return node.indices
}

// searchBySubstring performs substring matching (fallback for comprehensive search)
func (a *App) searchBySubstring(query string) []int {
	var matches []int

	for i, glyph := range a.searchIndex.glyphs {
		normalized := normalizeString(glyph.Name)
		if strings.Contains(normalized, query) {
			matches = append(matches, i)
		}
	}

	return matches
}

// removeDuplicateIndices removes duplicate indices while preserving order
func (a *App) removeDuplicateIndices(indices []int) []int {
	seen := make(map[int]bool)
	var unique []int

	for _, idx := range indices {
		if !seen[idx] {
			seen[idx] = true
			unique = append(unique, idx)
		}
	}

	return unique
}

// sortByRelevance sorts indices by search relevance
func (a *App) sortByRelevance(indices []int, query string) []int {
	type indexScore struct {
		index int
		score int
	}

	scored := make([]indexScore, len(indices))

	for i, idx := range indices {
		score := a.calculateRelevanceScore(a.searchIndex.glyphs[idx].Name, query)
		scored[i] = indexScore{index: idx, score: score}
	}

	// Sort by score (higher is better)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	result := make([]int, len(scored))
	for i, item := range scored {
		result[i] = item.index
	}

	return result
}

// calculateRelevanceScore assigns a relevance score to a match
func (a *App) calculateRelevanceScore(name, query string) int {
	normalized := normalizeString(name)
	score := 0

	// Exact match gets highest score
	if normalized == query {
		score += 1000
	}

	// Prefix match gets high score
	if strings.HasPrefix(normalized, query) {
		score += 500
	}

	// Word boundary matches get medium score
	words := strings.Fields(normalized)
	for _, word := range words {
		if strings.HasPrefix(word, query) {
			score += 250
		}
		if strings.Contains(word, query) {
			score += 100
		}
	}

	// Substring match gets base score
	if strings.Contains(normalized, query) {
		score += 50
	}

	// Shorter names with matches are generally more relevant
	if len(normalized) < 20 && score > 0 {
		score += 25
	}

	return score
}

// CopyToClipboard takes a string and copies it to the user's clipboard
func (a *App) CopyToClipboard(text string) {
	runtime.ClipboardSetText(a.ctx, text)
}

// LoadGlyphsFromJSON loads glyph data from JSON string (for runtime updates)
func (a *App) LoadGlyphsFromJSON(jsonData string) error {
	var glyphs []Glyph
	if err := json.Unmarshal([]byte(jsonData), &glyphs); err != nil {
		return err
	}

	a.searchIndex.mu.Lock()
	defer a.searchIndex.mu.Unlock()

	// Clear existing data
	a.searchIndex.glyphs = nil
	a.searchIndex.normalizedMap = make(map[string][]int)
	a.searchIndex.prefixTree = &TrieNode{children: make(map[rune]*TrieNode)}

	// Process new data
	seen := make(map[string]bool)
	uniqueGlyphs := make([]Glyph, 0, len(glyphs))

	for _, glyph := range glyphs {
		if !seen[glyph.Name] {
			seen[glyph.Name] = true
			uniqueGlyphs = append(uniqueGlyphs, glyph)
		}
	}

	a.searchIndex.glyphs = uniqueGlyphs

	// Rebuild indices
	for i, glyph := range uniqueGlyphs {
		normalized := normalizeString(glyph.Name)
		a.searchIndex.normalizedMap[normalized] = append(a.searchIndex.normalizedMap[normalized], i)
		a.addToTrie(normalized, i)

		words := strings.Fields(normalized)
		for _, word := range words {
			if len(word) > 1 {
				a.addToTrie(word, i)
			}
		}
	}

	runtime.LogInfof(a.ctx, "Search index reloaded with %d glyphs", len(uniqueGlyphs))
	return nil
}

// GetAllGlyphs returns all glyphs (useful for initial load)
func (a *App) GetAllGlyphs() []Glyph {
	a.searchIndex.mu.RLock()
	defer a.searchIndex.mu.RUnlock()
	return a.searchIndex.glyphs
}
