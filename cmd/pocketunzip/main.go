package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"pocketunzip"
	"pocketunzip/internal/app"
	"pocketunzip/internal/db"
)

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	appDir := filepath.Join(configDir, "PocketUnzip")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Fatal(err)
	}

	database, err := db.Init(db.DefaultDBPath(appDir))
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	sevenZipPath := "7z"

	a := app.NewApp(sevenZipPath, database)

	err = wails.Run(&options.App{
		Title:  "PocketUnzip",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: pocketunzip.Assets,
		},
		OnStartup: a.Startup,
		Bind: []interface{}{
			a,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
