package archive

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
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
	args := buildArgs(req)
	cmd := exec.CommandContext(ctx, req.SevenZipPath, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return ExtractResult{Success: false, ExitErr: err}
	}

	go scanPipe(stdout, onLog)
	go scanPipe(stderr, onLog)

	err := cmd.Wait()
	return ExtractResult{Success: err == nil, ExitErr: err}
}

func scanPipe(r io.Reader, onLog LogHandler) {
	if onLog == nil || r == nil {
		return
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		onLog(scanner.Text())
	}
}

func IsPasswordError(output string) bool {
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
