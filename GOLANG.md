This file is a merged representation of a subset of the codebase, containing specifically included files, combined into a single document by Repomix.
The content has been processed where content has been compressed (code blocks are separated by ⋮---- delimiter).

# File Summary

## Purpose
This file contains a packed representation of a subset of the repository's contents that is considered the most important context.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Only files matching these patterns are included: **/*.go
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Content has been compressed - code blocks are separated by ⋮---- delimiter
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
 app.go
 appicon.png
 build
├──  appicon.png
├──  bin
│   └──  Gylte.app
│       └──  Contents
│           ├──  Info.plist
│           ├──  MacOS
│           └──  Resources
│               └──  iconfile.icns
├──  darwin
│   ├──  Info.dev.plist
│   └──  Info.plist
└──  windows
    ├──  icon.ico
    ├──  info.json
    ├──  installer
    │   ├──  project.nsi
    │   └──  wails_tools.nsh
    └──  wails.exe.manifest
 darwin
├──  Info.dev.plist
└──  Info.plist
 db_generator
├──  glyphs.json
├──  go.mod
├──  go.sum
└──  main.go
 frontend
├──  dist
│   ├──  assets
│   │   ├──  index-D1mb1ySM.css
│   │   ├──  index-Dh1o15Ng.js
│   │   └──  mono-nf-ChJ0ab12.webp
│   └──  index.html
├──  index.html
├──  jsconfig.json
├──  package.json
├── 󰕥 package.json.md5
├── 󰂺 README.md
├──  src
│   ├──  App.css
│   ├──  App.svelte
│   ├──  assets
│   │   ├──  fonts
│   │   │   ├──  nunito-v16-latin-regular.woff2
│   │   │   └──  OFL.txt
│   │   └──  images
│   │       └──  logo-universal.png
│   ├──  comps
│   │   ├──  aniToast.svelte
│   │   ├──  mono-nf.webp
│   │   ├──  monogram-nf.jpg
│   │   └── 󰕙 monogram-nf.svg
│   ├──  glyphs.json
│   ├──  main.ts
│   ├──  style.css
│   └──  vite-env.d.ts
├──  svelte.config.js
├──  tsconfig.node.json
├──  vite.config.js
├──  wailsjs
│   ├──  go
│   │   ├──  main
│   │   │   ├──  App.d.ts
│   │   │   └──  App.js
│   │   └──  models.ts
│   └──  runtime
│       ├──  package.json
│       ├──  runtime.d.ts
│       └──  runtime.js
└──  yarn.lock
 go.mod
 go.sum
 GOLANG.md
 gylte.db
 gylte.db.bak
 Icon\r
 LICENSE
 main.go
󰂺 README.md
 REPOMIX.md
 wails.json
```

# Files

## File: db_generator/main.go
```go
package main
⋮----
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)
⋮----
"database/sql"
"encoding/json"
"fmt"
"log"
"os"
⋮----
_ "modernc.org/sqlite"
⋮----
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}
⋮----
func main()
⋮----
func run() error
⋮----
// Read and parse JSON
⋮----
// Initialize database
⋮----
// Populate database
⋮----
// Example queries
⋮----
func loadGlyphs(filename string) ([]Glyph, error)
⋮----
data, err := os.ReadFile(filename) // ioutil.ReadAll is deprecated
⋮----
var glyphs []Glyph
⋮----
func initDB(filename string) (*sql.DB, error)
⋮----
// Create table with proper constraints and indexing
⋮----
func populateDB(db *sql.DB, glyphs []Glyph) error
⋮----
defer tx.Rollback() // Rollback if not committed
⋮----
func runExampleQueries(db *sql.DB) error
⋮----
// 1. Get a specific glyph by name
⋮----
// 2. Search for glyphs by pattern
⋮----
// 3. Count total glyphs
⋮----
func getGlyphByName(db *sql.DB, name string) (*Glyph, error)
⋮----
var g Glyph
⋮----
func searchGlyphs(db *sql.DB, pattern string) ([]Glyph, error)
⋮----
var g Glyph
⋮----
func countGlyphs(db *sql.DB) (int, error)
⋮----
var count int
```

## File: main.go
```go
package main
⋮----
import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)
⋮----
"embed"
⋮----
"github.com/wailsapp/wails/v2"
"github.com/wailsapp/wails/v2/pkg/options"
"github.com/wailsapp/wails/v2/pkg/options/assetserver"
⋮----
//go:embed all:frontend/dist
var assets embed.FS
⋮----
func main()
⋮----
// Create an instance of the app structure
⋮----
// Create application with options
```

## File: app.go
```go
package main
⋮----
import (
	"context"
	"database/sql"
	"log"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)
⋮----
"context"
"database/sql"
"log"
"sort"
"strings"
⋮----
"github.com/wailsapp/wails/v2/pkg/runtime"
_ "modernc.org/sqlite"
⋮----
// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}
⋮----
// Glyph struct for database results
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}
⋮----
// GlyphMatch represents a glyph with its fuzzy match score
type GlyphMatch struct {
	Glyph
	Score int
}
⋮----
// NewApp creates a new App application struct
func NewApp() *App
⋮----
// startup is called when the app starts. The context is saved
// and the database connection is initialized.
func (a *App) startup(ctx context.Context)
⋮----
var err error
⋮----
// fuzzyMatch implements fzf-style fuzzy matching
// Returns a score (higher is better) and whether it matches
func fuzzyMatch(pattern, text string) (int, bool)
⋮----
// Exact substring match gets highest priority
⋮----
// Bonus for matches at the beginning
⋮----
// Bonus for shorter strings (more relevant)
⋮----
// Fuzzy matching: all characters of pattern must appear in order in text
⋮----
// Character matches
⋮----
// Bonus for consecutive matches
⋮----
// Bonus for matches at word boundaries
⋮----
// All pattern characters must be matched
⋮----
// Penalty for longer strings
⋮----
// GetGlyphs retrieves all glyphs or filters them by a search term with fuzzy matching
func (a *App) GetGlyphs(searchTerm string) ([]Glyph, error)
⋮----
// Always get all glyphs from database
⋮----
var allGlyphs []Glyph
⋮----
var glyph Glyph
⋮----
// If no search term, return all glyphs
⋮----
// Apply fuzzy matching
var matches []GlyphMatch
⋮----
// Sort by score (highest first)
⋮----
// Convert back to []Glyph
var result []Glyph
⋮----
// CopyToClipboard takes a string and copies it to the user's clipboard.
func (a *App) CopyToClipboard(text string)
```
