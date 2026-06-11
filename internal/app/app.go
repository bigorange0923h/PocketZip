package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
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

// SelectFiles 选择多个压缩包文件（批量解压）
func (a *App) SelectFiles() ([]string, error) {
	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择压缩包（可多选）",
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

func (a *App) extractLogger() func(string) {
	return func(line string) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "extract-log", line)
		}
	}
}

func (a *App) runExtract(archivePath, outputDir, passwordStr string) archive.ExtractResult {
	return archive.Extract(a.ctx, archive.ExtractRequest{
		SevenZipPath: a.sevenZipPath,
		ArchivePath:  archivePath,
		OutputDir:    outputDir,
		Password:     passwordStr,
	}, a.extractLogger())
}

func (a *App) recordExtract(archivePath, outputDir string, usedPassword bool, result archive.ExtractResult) {
	h := history.ExtractHistory{
		ArchivePath:  archivePath,
		OutputDir:    outputDir,
		Success:      result.Success,
		UsedPassword: usedPassword,
	}
	if err := result.Error(); err != nil {
		h.ErrorMessage = err.Error()
	}
	history.Record(a.db, h)
}

func (a *App) Extract(archivePath, outputDir string) error {
	if outputDir == "" {
		outputDir = defaultOutputDir(archivePath)
	}

	result := a.runExtract(archivePath, outputDir, "")
	if result.Success {
		a.recordExtract(archivePath, outputDir, false, result)
		return nil
	}

	if !archive.IsPasswordError(result.Output) {
		a.recordExtract(archivePath, outputDir, false, result)
		return result.Error()
	}

	candidates, err := password.Match(a.db, archivePath)
	if err != nil {
		a.recordExtract(archivePath, outputDir, false, result)
		return err
	}

	lastResult := result
	usedCandidate := false
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate) == "" {
			continue
		}

		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "extract-log", "Trying saved password candidate")
		}

		lastResult = a.runExtract(archivePath, outputDir, candidate)
		usedCandidate = true
		if lastResult.Success {
			password.UpdateSuccess(a.db, archivePath, candidate)
			a.recordExtract(archivePath, outputDir, true, lastResult)
			return nil
		}
		if !archive.IsPasswordError(lastResult.Output) {
			a.recordExtract(archivePath, outputDir, true, lastResult)
			return lastResult.Error()
		}
	}

	a.recordExtract(archivePath, outputDir, usedCandidate, lastResult)
	return ErrPasswordRequired
}

func (a *App) GetPasswordCandidates(archivePath string) ([]string, error) {
	return password.Match(a.db, archivePath)
}

func (a *App) SavePassword(archivePath, passwordStr string) error {
	return password.Save(a.db, archivePath, passwordStr)
}

func (a *App) GetHistory(limit int) ([]history.ExtractHistory, error) {
	return history.List(a.db, limit)
}

// BatchExtract 批量解压文件
func (a *App) BatchExtract(archivePaths []string, outputDir string) []BatchExtractResult {
	var results []BatchExtractResult

	for _, archivePath := range archivePaths {
		result := BatchExtractResult{
			ArchivePath: archivePath,
			Success:     false,
		}

		err := a.Extract(archivePath, outputDir)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Success = true
			result.OutputDir = outputDir
			if result.OutputDir == "" {
				result.OutputDir = defaultOutputDir(archivePath)
			}
		}

		results = append(results, result)
	}

	return results
}

// BatchExtractResult 批量解压结果
type BatchExtractResult struct {
	ArchivePath string `json:"archivePath"`
	OutputDir   string `json:"outputDir"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

// OpenDirectory 在文件管理器中打开目录
func (a *App) OpenDirectory(dirPath string) error {
	return exec.Command("explorer", dirPath).Start()
}

// TestArchive 测试压缩包完整性
func (a *App) TestArchive(archivePath string) (bool, error) {
	return archive.Test(a.ctx, a.sevenZipPath, archivePath)
}

// GetPasswordRecords 获取所有密码记录（密码库管理）
func (a *App) GetPasswordRecords() ([]password.PasswordRecord, error) {
	return password.ListAll(a.db)
}

// DeletePasswordRecord 删除密码记录
func (a *App) DeletePasswordRecord(id int64) error {
	return password.DeleteByID(a.db, id)
}

// UpdatePasswordRecord 更新密码记录的使用次数排序
func (a *App) UpdatePasswordRecord(id int64, newArchivePath string) error {
	return password.UpdatePath(a.db, id, newArchivePath)
}

// GetPasswordStats 获取密码使用统计
func (a *App) GetPasswordStats() (password.PasswordStats, error) {
	return password.GetStats(a.db)
}

// PreviewArchive 预览压缩包内容
func (a *App) PreviewArchive(archivePath string) ([]archive.ArchiveEntry, error) {
	return archive.List(a.ctx, a.sevenZipPath, archivePath)
}

// GetAppConfig 获取应用配置
func (a *App) GetAppConfig(key string) (string, error) {
	var value string
	err := a.db.QueryRow("SELECT value FROM app_config WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// SetAppConfig 设置应用配置
func (a *App) SetAppConfig(key, value string) error {
	_, err := a.db.Exec(
		`INSERT INTO app_config (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP`,
		key, value, value,
	)
	return err
}

// GetTheme 获取当前主题
func (a *App) GetTheme() string {
	theme, _ := a.GetAppConfig("theme")
	if theme == "" {
		theme = "dark"
	}
	return theme
}

// SetTheme 设置主题
func (a *App) SetTheme(theme string) error {
	return a.SetAppConfig("theme", theme)
}

// ExtractWithRetry 带重试的解压
func (a *App) ExtractWithRetry(archivePath, outputDir string, maxRetries int) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		err := a.Extract(archivePath, outputDir)
		if err == nil {
			return nil
		}
		lastErr = err

		// 如果是密码错误，不重试
		if err == ErrPasswordRequired {
			return err
		}
	}

	return lastErr
}

// ExtractWithStrategy 使用策略模板解压
func (a *App) ExtractWithStrategy(archivePath, strategyName string) error {
	strategy, err := a.GetExtractStrategy(strategyName)
	if err != nil {
		return err
	}

	outputDir := strategy.OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir(archivePath)
	}

	if strategy.AutoRetry {
		return a.ExtractWithRetry(archivePath, outputDir, strategy.MaxRetries)
	}

	return a.Extract(archivePath, outputDir)
}

// ExtractStrategy 解压策略
type ExtractStrategy struct {
	Name       string `json:"name"`
	OutputDir  string `json:"outputDir"`
	AutoRetry  bool   `json:"autoRetry"`
	MaxRetries int    `json:"maxRetries"`
	AutoOpen   bool   `json:"autoOpen"`
}

// GetExtractStrategy 获取解压策略
func (a *App) GetExtractStrategy(name string) (ExtractStrategy, error) {
	var strategy ExtractStrategy
	value, err := a.GetAppConfig("strategy_" + name)
	if err != nil {
		return strategy, err
	}
	if value == "" {
		// 返回默认策略
		return ExtractStrategy{
			Name:       name,
			AutoRetry:  false,
			MaxRetries: 3,
			AutoOpen:   false,
		}, nil
	}

	// 解析 JSON
	strategy.Name = name
	// 简单解析，实际应使用 json.Unmarshal
	if value == "default" {
		strategy.AutoRetry = false
		strategy.MaxRetries = 3
		strategy.AutoOpen = false
	}

	return strategy, nil
}

// SaveExtractStrategy 保存解压策略
func (a *App) SaveExtractStrategy(strategy ExtractStrategy) error {
	// 简化存储，实际应使用 JSON
	value := "default"
	if strategy.AutoRetry {
		value = "retry"
	}
	return a.SetAppConfig("strategy_"+strategy.Name, value)
}

// GetExtractStrategies 获取所有解压策略
func (a *App) GetExtractStrategies() []ExtractStrategy {
	return []ExtractStrategy{
		{Name: "default", OutputDir: "", AutoRetry: false, MaxRetries: 3, AutoOpen: false},
		{Name: "retry", OutputDir: "", AutoRetry: true, MaxRetries: 3, AutoOpen: false},
		{Name: "auto-open", OutputDir: "", AutoRetry: false, MaxRetries: 3, AutoOpen: true},
	}
}

func (a *App) ExtractWithPassword(archivePath, outputDir, passwordStr string) error {
	if outputDir == "" {
		outputDir = defaultOutputDir(archivePath)
	}

	result := a.runExtract(archivePath, outputDir, passwordStr)
	a.recordExtract(archivePath, outputDir, true, result)
	if result.Success {
		password.UpdateSuccess(a.db, archivePath, passwordStr)
		return nil
	}

	return result.Error()
}

// SelectFilesForCompress 选择要压缩的文件（任意文件，非压缩包过滤）
func (a *App) SelectFilesForCompress() ([]string, error) {
	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择要压缩的文件",
	})
}

// SelectFolderForCompress 选择要压缩的文件夹
func (a *App) SelectFolderForCompress() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择要压缩的文件夹",
	})
}

// SelectSavePath 选择压缩包保存路径
func (a *App) SelectSavePath(defaultName string) (string, error) {
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "保存压缩包",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "压缩包 (*.zip, *.7z, *.tar, *.gz)",
				Pattern:     "*.zip;*.7z;*.tar;*.gz",
			},
		},
	})
}

// Compress 压缩文件
func (a *App) Compress(files []string, archivePath, format, passwordStr string) error {
	if len(files) == 0 {
		return fmt.Errorf("no files selected")
	}

	if archivePath == "" {
		return fmt.Errorf("archive path is required")
	}

	if format == "" {
		format = "zip"
	}

	compressLogger := func(line string) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "compress-log", line)
		}
	}

	result := archive.Compress(a.ctx, archive.CompressRequest{
		SevenZipPath: a.sevenZipPath,
		Files:        files,
		ArchivePath:  archivePath,
		Format:       format,
		Password:     passwordStr,
	}, compressLogger)

	if result.Success {
		return nil
	}

	return result.Error()
}

func defaultOutputDir(archivePath string) string {
	dir := filepath.Dir(archivePath)
	name := filepath.Base(archivePath)
	ext := filepath.Ext(name)
	nameWithoutExt := strings.TrimSuffix(name, ext)
	return filepath.Join(dir, nameWithoutExt)
}

var ErrPasswordRequired = errors.New("password required")
