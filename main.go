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
		Width:  500,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 19, G: 19, B: 19, A: 50},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
