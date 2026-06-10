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

// find7Zip 查找 7z 可执行文件的路径
// 优先级：同目录 > third_party/7zip > 当前工作目录 > PATH
func find7Zip() string {
	// 优先查找同目录下的 7z.exe
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)
	local7z := filepath.Join(exeDir, "7z.exe")
	if _, err := os.Stat(local7z); err == nil {
		return local7z
	}

	// macOS/Linux: 查找同目录下的 7z（无扩展名）
	local7zUnix := filepath.Join(exeDir, "7z")
	if _, err := os.Stat(local7zUnix); err == nil {
		return local7zUnix
	}

	// 查找 third_party/7zip 目录（相对于可执行文件）
	thirdParty := filepath.Join(exeDir, "third_party", "7zip", "7z.exe")
	if _, err := os.Stat(thirdParty); err == nil {
		return thirdParty
	}

	// 查找当前工作目录下的 third_party/7zip
	cwd, _ := os.Getwd()
	cwd7z := filepath.Join(cwd, "third_party", "7zip", "7z.exe")
	if _, err := os.Stat(cwd7z); err == nil {
		return cwd7z
	}

	// macOS/Linux: 查找 third_party/7zip/7z（无扩展名）
	thirdPartyUnix := filepath.Join(exeDir, "third_party", "7zip", "7z")
	if _, err := os.Stat(thirdPartyUnix); err == nil {
		return thirdPartyUnix
	}

	cwd7zUnix := filepath.Join(cwd, "third_party", "7zip", "7z")
	if _, err := os.Stat(cwd7zUnix); err == nil {
		return cwd7zUnix
	}

	// 假设在 PATH 中
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

	sevenZipPath := find7Zip()

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
