package app

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"

	"pocketunzip/internal/archive"
	"pocketunzip/internal/history"
	"pocketunzip/internal/password"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	db           *sql.DB
	sevenZipPath string
}

func NewApp(sevenZipPath string, db *sql.DB) *App {
	return &App{
		sevenZipPath: sevenZipPath,
		db:           db,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择压缩包",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "压缩包 (*.zip, *.7z, *.rar, *.tar, *.gz)",
				Pattern:     "*.zip;*.7z;*.rar;*.tar;*.gz;*.bz2;*.xz",
			},
		},
	})
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择解压目录",
	})
}

func (a *App) Extract(archivePath, outputDir string) error {
	if outputDir == "" {
		outputDir = defaultOutputDir(archivePath)
	}

	onLog := func(line string) {
		runtime.EventsEmit(a.ctx, "extract-log", line)
	}

	result := archive.Extract(a.ctx, archive.ExtractRequest{
		SevenZipPath: a.sevenZipPath,
		ArchivePath:  archivePath,
		OutputDir:    outputDir,
	}, onLog)

	h := history.ExtractHistory{
		ArchivePath: archivePath,
		OutputDir:   outputDir,
		Success:     result.Success,
	}
	if result.ExitErr != nil {
		h.ErrorMessage = result.ExitErr.Error()
	}
	history.Record(a.db, h)

	if !result.Success && archive.IsPasswordError(h.ErrorMessage) {
		return ErrPasswordRequired
	}

	if !result.Success {
		return result.ExitErr
	}

	return nil
}

func (a *App) GetPasswordCandidates(archivePath string) ([]string, error) {
	return password.Match(a.db, archivePath)
}

func (a *App) SavePassword(archivePath, passwordStr string) error {
	return password.Save(a.db, archivePath, passwordStr)
}

func (a *App) ExtractWithPassword(archivePath, outputDir, passwordStr string) error {
	if outputDir == "" {
		outputDir = defaultOutputDir(archivePath)
	}

	onLog := func(line string) {
		runtime.EventsEmit(a.ctx, "extract-log", line)
	}

	result := archive.Extract(a.ctx, archive.ExtractRequest{
		SevenZipPath: a.sevenZipPath,
		ArchivePath:  archivePath,
		OutputDir:    outputDir,
		Password:     passwordStr,
	}, onLog)

	h := history.ExtractHistory{
		ArchivePath:  archivePath,
		OutputDir:    outputDir,
		Success:      result.Success,
		UsedPassword: true,
	}
	if result.ExitErr != nil {
		h.ErrorMessage = result.ExitErr.Error()
	}
	history.Record(a.db, h)

	if result.Success {
		password.UpdateSuccess(a.db, archivePath, passwordStr)
	}

	if !result.Success {
		return result.ExitErr
	}

	return nil
}

func defaultOutputDir(archivePath string) string {
	dir := filepath.Dir(archivePath)
	name := filepath.Base(archivePath)
	ext := filepath.Ext(name)
	nameWithoutExt := strings.TrimSuffix(name, ext)
	return filepath.Join(dir, nameWithoutExt)
}

var ErrPasswordRequired = errors.New("password required")
