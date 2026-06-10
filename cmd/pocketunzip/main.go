package main

import (
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"pocketunzip"
	"pocketunzip/internal/app"
)

func main() {
	app := app.NewApp()

	err := wails.Run(&options.App{
		Title:  "PocketUnzip",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: pocketunzip.Assets,
		},
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
