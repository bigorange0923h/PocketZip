package archive

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type CompressRequest struct {
	SevenZipPath string
	Files        []string
	ArchivePath  string
	Format       string
	Password     string
}

type CompressResult struct {
	Success bool
	ExitErr error
	Output  string
}

func (r CompressResult) Error() error {
	if r.Success {
		return nil
	}

	output := strings.TrimSpace(r.Output)
	if r.ExitErr == nil {
		if output == "" {
			return errors.New("compress failed")
		}
		return errors.New(output)
	}

	if output == "" {
		return r.ExitErr
	}
	return fmt.Errorf("%w: %s", r.ExitErr, output)
}

func buildCompressArgs(req CompressRequest) []string {
	args := []string{"a"}

	// Archive path
	args = append(args, req.ArchivePath)

	// Source files
	args = append(args, req.Files...)

	// Format
	if req.Format != "" {
		args = append(args, "-t"+req.Format)
	}

	// Password
	if req.Password != "" {
		args = append(args, "-p"+req.Password)
		// Enable header encryption for 7z format
		if req.Format == "7z" {
			args = append(args, "-mhe=on")
		}
	}

	// Auto-confirm
	args = append(args, "-y")

	return args
}

func Compress(ctx context.Context, req CompressRequest, onLog LogHandler) CompressResult {
	if ctx == nil {
		ctx = context.Background()
	}

	// Validate files exist
	for _, f := range req.Files {
		if _, err := os.Stat(f); err != nil {
			return CompressResult{
				Success: false,
				ExitErr: fmt.Errorf("file not found: %s", f),
			}
		}
	}

	args := buildCompressArgs(req)
	cmd := exec.CommandContext(ctx, req.SevenZipPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return CompressResult{Success: false, ExitErr: err}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return CompressResult{Success: false, ExitErr: err}
	}

	if err := cmd.Start(); err != nil {
		return CompressResult{Success: false, ExitErr: err}
	}

	var output bytes.Buffer
	var outputMu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go scanPipe(stdout, onLog, &output, &outputMu, &wg)
	go scanPipe(stderr, onLog, &output, &outputMu, &wg)

	err = cmd.Wait()
	wg.Wait()

	return CompressResult{Success: err == nil, ExitErr: err, Output: output.String()}
}
