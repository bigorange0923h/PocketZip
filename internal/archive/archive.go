package archive

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ExtractRequest struct {
	SevenZipPath string
	ArchivePath  string
	OutputDir    string
	Password     string
}

type LogHandler func(line string)

type ExtractResult struct {
	Success bool
	ExitErr error
	Output  string
}

func (r ExtractResult) Error() error {
	if r.Success {
		return nil
	}

	output := strings.TrimSpace(r.Output)
	if r.ExitErr == nil {
		if output == "" {
			return errors.New("extract failed")
		}
		return errors.New(output)
	}

	if output == "" {
		return r.ExitErr
	}
	return fmt.Errorf("%w: %s", r.ExitErr, output)
}

func buildArgs(req ExtractRequest) []string {
	args := []string{"x", req.ArchivePath, "-o" + req.OutputDir}
	if req.Password != "" {
		args = append(args, "-p"+req.Password)
	}
	args = append(args, "-y")
	return args
}

func Extract(ctx context.Context, req ExtractRequest, onLog LogHandler) ExtractResult {
	if ctx == nil {
		ctx = context.Background()
	}

	// 自动创建输出目录
	if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
		return ExtractResult{Success: false, ExitErr: err}
	}

	args := buildArgs(req)
	cmd := exec.CommandContext(ctx, req.SevenZipPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return ExtractResult{Success: false, ExitErr: err}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return ExtractResult{Success: false, ExitErr: err}
	}

	if err := cmd.Start(); err != nil {
		return ExtractResult{Success: false, ExitErr: err}
	}

	var output bytes.Buffer
	var outputMu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go scanPipe(stdout, onLog, &output, &outputMu, &wg)
	go scanPipe(stderr, onLog, &output, &outputMu, &wg)

	err = cmd.Wait()
	wg.Wait()

	return ExtractResult{Success: err == nil, ExitErr: err, Output: output.String()}
}

func scanPipe(r io.Reader, onLog LogHandler, output *bytes.Buffer, outputMu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if r == nil {
		return
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if onLog != nil {
			onLog(line)
		}
		outputMu.Lock()
		output.WriteString(line)
		output.WriteByte('\n')
		outputMu.Unlock()
	}
}

func IsPasswordError(output string) bool {
	normalizedOutput := strings.ToLower(output)
	passwordKeywords := []string{
		"wrong password",
		"password is incorrect",
		"cannot open encrypted archive",
		"can not open encrypted archive",
		"data error in encrypted file",
		"enter password",
		"encrypted archive",
		"密码错误",
	}
	for _, keyword := range passwordKeywords {
		if strings.Contains(normalizedOutput, keyword) {
			return true
		}
	}

	keywords := []string{
		"Wrong password",
		"Wrong password?",
		"密码错误",
		"Cannot open encrypted archive",
	}
	for _, keyword := range keywords {
		if strings.Contains(output, keyword) {
			return true
		}
	}
	return false
}

// Test 测试压缩包完整性
func Test(ctx context.Context, sevenZipPath, archivePath string) (bool, error) {
	args := []string{"t", archivePath, "-y"}
	cmd := exec.CommandContext(ctx, sevenZipPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	// 检查输出中是否包含错误信息
	outputStr := string(output)
	if strings.Contains(outputStr, "ERROR") || strings.Contains(outputStr, "错误") {
		return false, nil
	}

	return true, nil
}

// ArchiveEntry 压缩包内的文件/目录条目
type ArchiveEntry struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"isDir"`
	Modified string `json:"modified"`
}

// List 列出压缩包内容
func List(ctx context.Context, sevenZipPath, archivePath string) ([]ArchiveEntry, error) {
	args := []string{"l", archivePath, "-y"}
	cmd := exec.CommandContext(ctx, sevenZipPath, args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseListOutput(string(output)), nil
}

// parseListOutput 解析 7z l 命令的输出
func parseListOutput(output string) []ArchiveEntry {
	var entries []ArchiveEntry
	lines := strings.Split(output, "\n")

	// 找到文件列表开始的标记
	inList := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行
		if line == "" {
			continue
		}

		// 检测列表开始（以 "Date Time Attr" 或类似格式开头）
		if strings.HasPrefix(line, "Date") && strings.Contains(line, "Attr") {
			inList = true
			continue
		}

		// 检测列表结束（以 "---" 或 "----------------" 开头）
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "====") {
			inList = false
			continue
		}

		// 解析文件条目
		if inList && len(line) > 20 {
			entry := parseEntryLine(line)
			if entry != nil {
				entries = append(entries, *entry)
			}
		}
	}

	return entries
}

// parseEntryLine 解析单行文件条目
func parseEntryLine(line string) *ArchiveEntry {
	// 7z l 输出格式：
	// 2024-01-01 12:00:00 ....A 12345  file.txt
	parts := strings.Fields(line)
	if len(parts) < 5 {
		return nil
	}

	entry := &ArchiveEntry{
		Modified: parts[0] + " " + parts[1],
	}

	// 解析属性
	attr := parts[2]
	entry.IsDir = strings.Contains(attr, "D")

	// 解析大小
	sizeStr := parts[3]
	if size, err := parseInt64(sizeStr); err == nil {
		entry.Size = size
	}

	// 解析路径（可能是多个字段）
	entry.Path = strings.Join(parts[4:], " ")

	return entry
}

func parseInt64(s string) (int64, error) {
	var n int64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number")
		}
		n = n*10 + int64(c-'0')
	}
	return n, nil
}
