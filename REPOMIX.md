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
- Only files matching these patterns are included: **/*.go, frontend/*.js, frontend/**/*.ts, **/*.svelte
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Content has been compressed - code blocks are separated by ⋮---- delimiter
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
db_generator/
  main.go
frontend/
  src/
    comps/
      aniToast.svelte
    App.svelte
    main.ts
    vite-env.d.ts
  wailsjs/
    go/
      main/
        App.d.ts
      models.ts
    runtime/
      runtime.d.ts
  svelte.config.js
  vite.config.js
app.go
main.go
```

# Files

## File: db_generator/main.go
```go
package main
⋮----
import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)
⋮----
"database/sql"
"encoding/json"
"io/ioutil"
"log"
"os"
⋮----
_ "modernc.org/sqlite" // Pure Go SQLite driver
⋮----
// Glyph struct to match the JSON structure
type Glyph struct {
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
}
⋮----
func main()
⋮----
// Open and read the JSON file
⋮----
var glyphs []Glyph
⋮----
// Create a new SQLite database file
// Use "sqlite" instead of "sqlite3" for modernc.org/sqlite
⋮----
// Create the glyphs table
⋮----
// Insert unique glyphs into the database
```

## File: frontend/src/comps/aniToast.svelte
```
<script>
    export let width = "24";
    export let height = "24";
    export let fill = "{fill}";
    export let customClass = "animatedToast"; // Use a different name to avoid conflict with the class attribute
</script>

<svg
    xmlns="http://www.w3.org/2000/svg"
    width={width}
    height={height}
    viewBox="0 0 24 24"
    {...$$props}
>
    <rect width="7.33" height="7.33" x="1" y="1" {fill}>
        <animate
            id="SVGzjrPLenI"
            attributeName="x"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="y"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="width"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="0;SVGXAURnSRI.end+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="8.33" y="1" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="1" y="8.33" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.125s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="15.66" y="1" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="8.33" y="8.33" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="1" y="15.66" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="1;4;1"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.25s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="15.66" y="8.33" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="8.33" y="15.66" {fill}>
        <animate
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="8.33;11.33;8.33"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.375s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
    <rect width="7.33" height="7.33" x="15.66" y="15.66" {fill}>
        <animate
            id="SVGXAURnSRI"
            attributeName="x"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="y"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="15.66;18.66;15.66"
        />
        <animate
            attributeName="width"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
        <animate
            attributeName="height"
            begin="SVGzjrPLenI.begin+0.5s"
            dur="0.75s"
            values="7.33;1.33;7.33"
        />
    </rect>
</svg>
```

## File: frontend/src/main.ts
```typescript
import "./style.css"; // Keep your existing CSS file name
import App from "./App.svelte";
```

## File: frontend/src/vite-env.d.ts
```typescript
/// <reference types="svelte" />
/// <reference types="vite/client" />
```

## File: frontend/wailsjs/runtime/runtime.d.ts
```typescript
/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
⋮----
export interface Position {
    x: number;
    y: number;
}
⋮----
export interface Size {
    w: number;
    h: number;
}
⋮----
export interface Screen {
    isCurrent: boolean;
    isPrimary: boolean;
    width : number
    height : number
}
⋮----
// Environment information such as platform, buildtype, ...
export interface EnvironmentInfo {
    buildType: string;
    platform: string;
    arch: string;
}
⋮----
// [EventsEmit](https://wails.io/docs/reference/runtime/events#eventsemit)
// emits the given event. Optional data may be passed with the event.
// This will trigger any event listeners.
export function EventsEmit(eventName: string, ...data: any): void;
⋮----
// [EventsOn](https://wails.io/docs/reference/runtime/events#eventson) sets up a listener for the given event name.
export function EventsOn(eventName: string, callback: (...data: any)
⋮----
// [EventsOnMultiple](https://wails.io/docs/reference/runtime/events#eventsonmultiple)
// sets up a listener for the given event name, but will only trigger a given number times.
export function EventsOnMultiple(eventName: string, callback: (...data: any)
⋮----
// [EventsOnce](https://wails.io/docs/reference/runtime/events#eventsonce)
// sets up a listener for the given event name, but will only trigger once.
export function EventsOnce(eventName: string, callback: (...data: any)
⋮----
// [EventsOff](https://wails.io/docs/reference/runtime/events#eventsoff)
// unregisters the listener for the given event name.
export function EventsOff(eventName: string, ...additionalEventNames: string[]): void;
⋮----
// [EventsOffAll](https://wails.io/docs/reference/runtime/events#eventsoffall)
// unregisters all listeners.
export function EventsOffAll(): void;
⋮----
// [LogPrint](https://wails.io/docs/reference/runtime/log#logprint)
// logs the given message as a raw message
export function LogPrint(message: string): void;
⋮----
// [LogTrace](https://wails.io/docs/reference/runtime/log#logtrace)
// logs the given message at the `trace` log level.
export function LogTrace(message: string): void;
⋮----
// [LogDebug](https://wails.io/docs/reference/runtime/log#logdebug)
// logs the given message at the `debug` log level.
export function LogDebug(message: string): void;
⋮----
// [LogError](https://wails.io/docs/reference/runtime/log#logerror)
// logs the given message at the `error` log level.
export function LogError(message: string): void;
⋮----
// [LogFatal](https://wails.io/docs/reference/runtime/log#logfatal)
// logs the given message at the `fatal` log level.
// The application will quit after calling this method.
export function LogFatal(message: string): void;
⋮----
// [LogInfo](https://wails.io/docs/reference/runtime/log#loginfo)
// logs the given message at the `info` log level.
export function LogInfo(message: string): void;
⋮----
// [LogWarning](https://wails.io/docs/reference/runtime/log#logwarning)
// logs the given message at the `warning` log level.
export function LogWarning(message: string): void;
⋮----
// [WindowReload](https://wails.io/docs/reference/runtime/window#windowreload)
// Forces a reload by the main application as well as connected browsers.
export function WindowReload(): void;
⋮----
// [WindowReloadApp](https://wails.io/docs/reference/runtime/window#windowreloadapp)
// Reloads the application frontend.
export function WindowReloadApp(): void;
⋮----
// [WindowSetAlwaysOnTop](https://wails.io/docs/reference/runtime/window#windowsetalwaysontop)
// Sets the window AlwaysOnTop or not on top.
export function WindowSetAlwaysOnTop(b: boolean): void;
⋮----
// [WindowSetSystemDefaultTheme](https://wails.io/docs/next/reference/runtime/window#windowsetsystemdefaulttheme)
// *Windows only*
// Sets window theme to system default (dark/light).
export function WindowSetSystemDefaultTheme(): void;
⋮----
// [WindowSetLightTheme](https://wails.io/docs/next/reference/runtime/window#windowsetlighttheme)
// *Windows only*
// Sets window to light theme.
export function WindowSetLightTheme(): void;
⋮----
// [WindowSetDarkTheme](https://wails.io/docs/next/reference/runtime/window#windowsetdarktheme)
// *Windows only*
// Sets window to dark theme.
export function WindowSetDarkTheme(): void;
⋮----
// [WindowCenter](https://wails.io/docs/reference/runtime/window#windowcenter)
// Centers the window on the monitor the window is currently on.
export function WindowCenter(): void;
⋮----
// [WindowSetTitle](https://wails.io/docs/reference/runtime/window#windowsettitle)
// Sets the text in the window title bar.
export function WindowSetTitle(title: string): void;
⋮----
// [WindowFullscreen](https://wails.io/docs/reference/runtime/window#windowfullscreen)
// Makes the window full screen.
export function WindowFullscreen(): void;
⋮----
// [WindowUnfullscreen](https://wails.io/docs/reference/runtime/window#windowunfullscreen)
// Restores the previous window dimensions and position prior to full screen.
export function WindowUnfullscreen(): void;
⋮----
// [WindowIsFullscreen](https://wails.io/docs/reference/runtime/window#windowisfullscreen)
// Returns the state of the window, i.e. whether the window is in full screen mode or not.
export function WindowIsFullscreen(): Promise<boolean>;
⋮----
// [WindowSetSize](https://wails.io/docs/reference/runtime/window#windowsetsize)
// Sets the width and height of the window.
export function WindowSetSize(width: number, height: number): void;
⋮----
// [WindowGetSize](https://wails.io/docs/reference/runtime/window#windowgetsize)
// Gets the width and height of the window.
export function WindowGetSize(): Promise<Size>;
⋮----
// [WindowSetMaxSize](https://wails.io/docs/reference/runtime/window#windowsetmaxsize)
// Sets the maximum window size. Will resize the window if the window is currently larger than the given dimensions.
// Setting a size of 0,0 will disable this constraint.
export function WindowSetMaxSize(width: number, height: number): void;
⋮----
// [WindowSetMinSize](https://wails.io/docs/reference/runtime/window#windowsetminsize)
// Sets the minimum window size. Will resize the window if the window is currently smaller than the given dimensions.
// Setting a size of 0,0 will disable this constraint.
export function WindowSetMinSize(width: number, height: number): void;
⋮----
// [WindowSetPosition](https://wails.io/docs/reference/runtime/window#windowsetposition)
// Sets the window position relative to the monitor the window is currently on.
export function WindowSetPosition(x: number, y: number): void;
⋮----
// [WindowGetPosition](https://wails.io/docs/reference/runtime/window#windowgetposition)
// Gets the window position relative to the monitor the window is currently on.
export function WindowGetPosition(): Promise<Position>;
⋮----
// [WindowHide](https://wails.io/docs/reference/runtime/window#windowhide)
// Hides the window.
export function WindowHide(): void;
⋮----
// [WindowShow](https://wails.io/docs/reference/runtime/window#windowshow)
// Shows the window, if it is currently hidden.
export function WindowShow(): void;
⋮----
// [WindowMaximise](https://wails.io/docs/reference/runtime/window#windowmaximise)
// Maximises the window to fill the screen.
export function WindowMaximise(): void;
⋮----
// [WindowToggleMaximise](https://wails.io/docs/reference/runtime/window#windowtogglemaximise)
// Toggles between Maximised and UnMaximised.
export function WindowToggleMaximise(): void;
⋮----
// [WindowUnmaximise](https://wails.io/docs/reference/runtime/window#windowunmaximise)
// Restores the window to the dimensions and position prior to maximising.
export function WindowUnmaximise(): void;
⋮----
// [WindowIsMaximised](https://wails.io/docs/reference/runtime/window#windowismaximised)
// Returns the state of the window, i.e. whether the window is maximised or not.
export function WindowIsMaximised(): Promise<boolean>;
⋮----
// [WindowMinimise](https://wails.io/docs/reference/runtime/window#windowminimise)
// Minimises the window.
export function WindowMinimise(): void;
⋮----
// [WindowUnminimise](https://wails.io/docs/reference/runtime/window#windowunminimise)
// Restores the window to the dimensions and position prior to minimising.
export function WindowUnminimise(): void;
⋮----
// [WindowIsMinimised](https://wails.io/docs/reference/runtime/window#windowisminimised)
// Returns the state of the window, i.e. whether the window is minimised or not.
export function WindowIsMinimised(): Promise<boolean>;
⋮----
// [WindowIsNormal](https://wails.io/docs/reference/runtime/window#windowisnormal)
// Returns the state of the window, i.e. whether the window is normal or not.
export function WindowIsNormal(): Promise<boolean>;
⋮----
// [WindowSetBackgroundColour](https://wails.io/docs/reference/runtime/window#windowsetbackgroundcolour)
// Sets the background colour of the window to the given RGBA colour definition. This colour will show through for all transparent pixels.
export function WindowSetBackgroundColour(R: number, G: number, B: number, A: number): void;
⋮----
// [ScreenGetAll](https://wails.io/docs/reference/runtime/window#screengetall)
// Gets the all screens. Call this anew each time you want to refresh data from the underlying windowing system.
export function ScreenGetAll(): Promise<Screen[]>;
⋮----
// [BrowserOpenURL](https://wails.io/docs/reference/runtime/browser#browseropenurl)
// Opens the given URL in the system browser.
export function BrowserOpenURL(url: string): void;
⋮----
// [Environment](https://wails.io/docs/reference/runtime/intro#environment)
// Returns information about the environment
export function Environment(): Promise<EnvironmentInfo>;
⋮----
// [Quit](https://wails.io/docs/reference/runtime/intro#quit)
// Quits the application.
export function Quit(): void;
⋮----
// [Hide](https://wails.io/docs/reference/runtime/intro#hide)
// Hides the application.
export function Hide(): void;
⋮----
// [Show](https://wails.io/docs/reference/runtime/intro#show)
// Shows the application.
export function Show(): void;
⋮----
// [ClipboardGetText](https://wails.io/docs/reference/runtime/clipboard#clipboardgettext)
// Returns the current text stored on clipboard
export function ClipboardGetText(): Promise<string>;
⋮----
// [ClipboardSetText](https://wails.io/docs/reference/runtime/clipboard#clipboardsettext)
// Sets a text on the clipboard
export function ClipboardSetText(text: string): Promise<boolean>;
⋮----
// [OnFileDrop](https://wails.io/docs/reference/runtime/draganddrop#onfiledrop)
// OnFileDrop listens to drag and drop events and calls the callback with the coordinates of the drop and an array of path strings.
export function OnFileDrop(callback: (x: number, y: number ,paths: string[])
⋮----
// [OnFileDropOff](https://wails.io/docs/reference/runtime/draganddrop#dragandddropoff)
// OnFileDropOff removes the drag and drop listeners and handlers.
export function OnFileDropOff() :void
⋮----
// Check if the file path resolver is available
export function CanResolveFilePaths(): boolean;
⋮----
// Resolves file paths for an array of files
export function ResolveFilePaths(files: File[]): void
```

## File: frontend/svelte.config.js
```javascript
// Consult https://svelte.dev/docs#compile-time-svelte-preprocess
// for more information about preprocessors
preprocess: vitePreprocess(),
```

## File: frontend/vite.config.js
```javascript
// https://vitejs.dev/config/
export default defineConfig({
plugins: [svelte()]
```

## File: frontend/wailsjs/go/models.ts
```typescript
export class Glyph
⋮----
static createFrom(source: any =
⋮----
constructor(source: any =
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

## File: frontend/src/App.svelte
```
<script lang="ts">
    import { onMount } from "svelte";
    import AniToast from "./comps/aniToast.svelte";
    import MonoNF from "./comps/mono-nf.webp";
    import { CopyToClipboard, GetGlyphs } from "../wailsjs/go/main/App";
    import { WindowMinimise, Quit } from "../wailsjs/runtime";

    // Define the type for a single glyph object
    interface Glyph {
        name: string;
        glyph: string;
    }

    let searchTerm = "";
    let filteredGlyphs: Glyph[] = [];
    let toastVisible = false;
    let searchTimeout: NodeJS.Timeout;
    let isLoading = true;

    // Load all glyphs on initial render
    onMount(async () => {
        try {
            // Get all glyphs with empty search term
            filteredGlyphs = await GetGlyphs("");
            isLoading = false;
        } catch (error) {
            console.error("Failed to load glyphs from Go backend:", error);
            isLoading = false;
        }
    });

    // Debounced search function for better performance
    const performSearch = async (query: string) => {
        try {
            filteredGlyphs = await GetGlyphs(query);
        } catch (error) {
            console.error("Search failed:", error);
        }
    };

    // Reactive statement with debouncing for search
    $: {
        // Clear previous timeout
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        // Debounce search by 150ms to avoid excessive backend calls
        if (!isLoading) {
            searchTimeout = setTimeout(() => {
                performSearch(searchTerm);
            }, 150);
        }
    }

    const handleGlyphClick = (glyph: string) => {
        CopyToClipboard(glyph);
        // Show "Copied!" toast message
        toastVisible = true;
        setTimeout(() => {
            toastVisible = false;
        }, 1900);
    };
</script>

<div id="app">
    <!-- Custom Title Bar -->
    <div class="title-bar">
        <div class="title draggable">
            <img 
                src={MonoNF}
                id="nf-icon" 
                alt=""
                width="31"
                height="31"
            />
        </div>
        <div class="spacer draggable"></div>
        <div class="window-controls">
            <button on:click={WindowMinimise}>-</button>
            <button on:click={Quit}>×</button>
        </div>
    </div>

    <!-- Search Input -->
    <div class="search-container draggable">
        <input 
            type="text" 
            autocomplete="on" 
            class="search-input" 
            placeholder="" 
            bind:value={searchTerm}
            disabled={isLoading}
        />
    </div>

    <!-- Loading State -->
    {#if isLoading}
        <div class="loading">Loading glyphs...</div>
    {:else}
        <!-- Glyph Grid -->
        <div class="glyph-grid draggable">
            {#each filteredGlyphs as item (item.name)}
                <div
                    class="glyph-card"
                    title={`${item.name}`}
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
    {/if}

    <!-- "Copied!" Toast Notification -->
    {#if toastVisible}
        <div class="toast">
            <AniToast width="12" height="12" fill="#45a847" class="animatedToast" />
        </div>
    {/if}
</div>

<style>
    .loading {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 200px;
        color: #888;
        font-size: 14px;
        background-color: #12121290;
    }
    
    /* Add any additional styles here */
</style>
```

## File: frontend/wailsjs/go/main/App.d.ts
```typescript
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
⋮----
export function CopyToClipboard(arg1:string):Promise<void>;
⋮----
export function GetGlyphs(arg1:string):Promise<Array<main.Glyph>>;
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
