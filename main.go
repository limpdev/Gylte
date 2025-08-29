package main

import (
	"context"
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
		OnDomReady:       app.onDomReady,
		OnBeforeClose:    app.onBeforeClose,
		OnShutdown:       app.onShutdown,
		Bind: []interface{}{
			app, // This automatically binds all public methods of App
		},
		// --- CSS properties for dragging the window ---
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// Optional lifecycle methods for the App
func (a *App) onDomReady(ctx context.Context) {
	// Called when the frontend Dom is ready
}

func (a *App) onBeforeClose(ctx context.Context) (prevent bool) {
	// Called when the application is about to quit,
	// either by clicking the window close button or calling runtime.Quit.
	// Returning true will cause the application to continue running.
	return false
}

func (a *App) onShutdown(ctx context.Context) {
	// Called during shutdown after OnBeforeClose
}
