package archive

import (
	"bufio"
	"context"
	"io"
	"os/exec"
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

func Extract(ctx context.Context, req ExtractRequest, onLog LogHandler) ExtractResult {
	args := []string{"x", req.ArchivePath, "-o" + req.OutputDir, "-y"}
	if req.Password != "" {
		args = append(args, "-p"+req.Password)
	}

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
