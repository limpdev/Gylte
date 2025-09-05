## `app.go`

```go
package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	a.db, err = sql.Open("sqlite3", "./gylte.db")
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
```

## `db_generator\main.go`

```go
package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
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
	db, err := sql.Open("sqlite3", "./gylte.db")
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
```

## `main.go`

```go
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed gylte.db
var db embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Wlyph",
		Width:  500, // Taller than wide
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Handler: &assetserver.Handler{
				Mux: &assetserver.Mux{
					"./gylte.db": db.ReadFile("gylte.db"),
				},
			},
		},
		Frameless:        true, // Key for custom title bar
		BackgroundColour: &options.RGBA{R: 19, G: 19, B: 19, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		// --- CSS properties for dragging the window ---
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
```

