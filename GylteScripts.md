## `app.go`

```go
package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// CopyToClipboard takes a string and copies it to the user's clipboard.
func (a *App) CopyToClipboard(text string) {
	runtime.ClipboardSetText(a.ctx, text)
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

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Gylte",
		Width:  500, // Taller than wide
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
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

## `App.svelte`

```svelte
<script lang="ts">
    import { onMount } from "svelte";
    import glyphData from "./glyphs.json";
    import AniToast from "./comps/aniToast.svelte";
    import { CopyToClipboard } from "../wailsjs/go/main/App";
    import { WindowMinimise, Quit } from "../wailsjs/runtime";

    // Define the type for a single glyph object
    interface Glyph {
        name: string;
        glyph: string;
    }

    let searchTerm = "";
    let filteredGlyphs: Glyph[] = [];
    let toastVisible = false;

    // Load all glyphs on initial render and remove duplicates
    onMount(() => {
        // Remove duplicates by name, keeping the first occurrence
        const uniqueGlyphs = glyphData.filter((item, index, self) => index === self.findIndex((g) => g.name === item.name));
        filteredGlyphs = uniqueGlyphs;
    });

    // Reactive statement to filter glyphs when search term changes
    $: {
        const sourceData = glyphData.filter((item, index, self) => index === self.findIndex((g) => g.name === item.name));

        if (searchTerm.trim() === "") {
            filteredGlyphs = sourceData;
        } else {
            const lowercasedFilter = searchTerm.toLowerCase();
            filteredGlyphs = sourceData.filter((item) => item.name.toLowerCase().includes(lowercasedFilter));
        }
    }

    const handleGlyphClick = (glyph: string) => {
        CopyToClipboard(glyph);
        // Show "Copied!" toast message
        toastVisible = true;
        setTimeout(() => {
            toastVisible = false;
        }, 1900); // Hide after the animation is almost done
    };
</script>

<div id="app">
    <!-- Custom Title Bar -->
    <div class="title-bar draggable">
        <div class="title"></div>
        <div class="spacer"></div>
        <div class="window-controls">
            <button on:click={WindowMinimise}>-</button>
            <button on:click={Quit}>×</button>
        </div>
    </div>

    <!-- Search Input -->
    <div class="search-container draggable">
        <input type="text" autocomplete="on" class="search-input" placeholder="" bind:value={searchTerm} />
    </div>

    <!-- Glyph Grid -->
    <div class="glyph-grid draggable">
        {#each filteredGlyphs as item, index (index)}
            <div
                class="glyph-card"
                title={`Click to copy "${item.glyph}"`}
                on:click={() => handleGlyphClick(item.glyph)}
                on:keydown={(e) => e.key === "Enter" && handleGlyphClick(item.glyph)}
                role="button"
                tabindex="0"
            >
                <span class="glyph-icon">{item.glyph}</span>
                <span class="glyph-name">{item.name}</span>
            </div>
        {/each}
    </div>

    <!-- "Copied!" Toast Notification -->
    {#if toastVisible}
        <div class="toast">
            <!-- IMPORT THE SVG FOR THE TOAST HERE AS A COMPONENT -->
            <AniToast width="12" height="12" fill="#45a847" class="animatedToast" />
        </div>
    {/if}
</div>

<style>
    /* You'll need to move your existing CSS here or import it */
    /* This is where your current style.css content should go */
</style>
```

## `main.ts`

```typescript
import "./style.css"; // Keep your existing CSS file name
import App from "./App.svelte";

const app = new App({
    target: document.getElementById("app")!,
});

export default app;
```

## `vite-env.d.ts`

```typescript
/// <reference types="svelte" />
/// <reference types="vite/client" />
```

