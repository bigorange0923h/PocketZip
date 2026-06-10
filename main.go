package main

import (
	"log"
	"os"
	"path/filepath"

	"pocketunzip/internal/app"
	"pocketunzip/internal/db"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func find7Zip() string {
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)

	candidates := []string{
		filepath.Join(exeDir, "7z.exe"),
		filepath.Join(exeDir, "7z"),
		filepath.Join(exeDir, "third_party", "7zip", "7z.exe"),
		filepath.Join(exeDir, "third_party", "7zip", "7z"),
	}

	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, "third_party", "7zip", "7z.exe"),
			filepath.Join(cwd, "third_party", "7zip", "7z"),
		)
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return "7z"
}

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

	a := app.NewApp(find7Zip(), database)

	err = wails.Run(&options.App{
		Title:  "PocketUnzip",
		Width:  1024,
		Height: 768,
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},
		AssetServer: &assetserver.Options{
			Assets: Assets,
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
